package otel

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"github.com/amenzhinsky/go-memexec"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/orb-community/orb/agent/backend"
	"github.com/orb-community/orb/agent/config"
	"github.com/orb-community/orb/agent/policies"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

var _ backend.Backend = (*openTelemetryBackend)(nil)

const DefaultBinary = "/usr/local/sbin/otelcol"

//go:embed otelcol-contrib
var openTelemtryContribBinary []byte

type openTelemetryBackend struct {
	logger    *zap.Logger
	startTime time.Time

	//policies
	policyRepo            policies.PolicyRepo
	policyConfigDirectory string
	agentTags             map[string]string

	// Context for controlling the context cancellation
	mainContext        context.Context
	runningCollectors  map[string]context.Context
	mainCancelFunction context.CancelFunc

	// MQTT Config for OTEL MQTT Exporter
	mqttConfig config.MQTTConfig
	mqttClient *mqtt.Client

	otlpMetricsTopic string
	otlpTracesTopic  string
	otlpLogsTopic    string

	otelReceiverHost string
	otelReceiverPort int

	metricsReceiver receiver.Metrics
	metricsExporter exporter.Metrics
	//tracesReceiver  receiver.Traces
	//tracesExporter  exporter.Traces
	//logsReceiver    receiver.Logs
	//logsExporter    exporter.Logs
}

// Configure initializes the backend with the given configuration
func (o openTelemetryBackend) Configure(logger *zap.Logger, repo policies.PolicyRepo, configuration map[string]string, otelConfig map[string]interface{}) error {
	o.logger = logger
	o.policyRepo = repo
	var err error
	o.policyConfigDirectory, err = os.MkdirTemp("", "otel-policies")
	if err != nil {
		return err
	}
	if agentTags, ok := otelConfig["agent_tags"]; ok {
		o.agentTags = agentTags.(map[string]string)
	}
	for k, v := range otelConfig {
		switch k {
		case "Host":
			o.otelReceiverHost = v.(string)
		case "Port":
			o.otelReceiverPort = v.(int)
		}
	}
	o.mqttConfig = config.MQTTConfig{
		Address:   "",
		Id:        "",
		Key:       "",
		ChannelID: "",
	}

	return nil
}

func (o openTelemetryBackend) SetCommsClient(agentID string, client *mqtt.Client, baseTopic string) {
	o.mqttClient = client
	otelBaseTopic := strings.Replace(baseTopic, "?", "otlp", 1)
	o.otlpMetricsTopic = fmt.Sprintf("%s/m/%c", otelBaseTopic, agentID[0])
	o.otlpTracesTopic = fmt.Sprintf("%s/t/%c", otelBaseTopic, agentID[0])
	o.otlpLogsTopic = fmt.Sprintf("%s/l/%c", otelBaseTopic, agentID[0])
}

func (o openTelemetryBackend) Version() (string, error) {
	executable, err := memexec.New(openTelemtryContribBinary)
	if err != nil {
		return "", err
	}
	defer func(executable *memexec.Exec) {
		err := executable.Close()
		if err != nil {
			o.logger.Error("error closing executable", zap.Error(err))
		}
	}(executable)
	ctx, cancel := context.WithTimeout(o.mainContext, 5*time.Second)
	defer cancel()
	cmd := executable.CommandContext(ctx, "--version")
	if cmd.Err != nil {
		return "", cmd.Err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	var versionOutput string
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			o.logger.Info("DEBUG", zap.String("line", scanner.Text()))
			versionOutput = scanner.Text()
		}
	}()
	if err := cmd.Start(); err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		return "", err
	}
	o.logger.Info("running opentelemetry-contrib version", zap.String("version", versionOutput))
	return versionOutput, nil
}

func (o openTelemetryBackend) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	// initialize otlpreceiver and mqttexporter for scraping
	if o.policyRepo == nil {
		return fmt.Errorf("backend not properly configured, call Configure() first")
	}
	o.runningCollectors = make(map[string]context.Context)
	o.mainCancelFunction = cancelFunc

	currentVersion, err := o.Version()
	if err != nil {
		o.logger.Error("error during ")
	}
	o.logger.Info("starting open-telemetry backend using version", zap.String("version", currentVersion))

	policiesData, err := o.policyRepo.GetAll()
	if err != nil {
		defer cancelFunc()
		o.logger.Error("failed to start otel backend, policies are absent")
		return err
	}
	for _, policyData := range policiesData {
		if err := o.ApplyPolicy(policyData, true); err != nil {
			o.logger.Error("failed to start otel backend, failed to apply policy", zap.Error(err))
			cancelFunc()
			return err
		}
		o.logger.Info("policy applied successfully", zap.String("policy_id", policyData.ID))
	}

	return nil
}

func (o openTelemetryBackend) Stop(ctx context.Context) error {
	o.logger.Info("stopping all running policies")
	o.mainCancelFunction()
	for policyID, policyCtx := range o.runningCollectors {
		o.logger.Debug("stopping policy context", zap.String("policy_id", policyID))
		policyCtx.Done()
	}
	return nil
}

func (o openTelemetryBackend) FullReset(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (o openTelemetryBackend) GetStartTime() time.Time {
	//TODO implement me
	panic("implement me")
}

func (o openTelemetryBackend) GetCapabilities() (map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (o openTelemetryBackend) GetRunningStatus() (backend.RunningStatus, string, error) {
	//TODO implement me
	panic("implement me")
}

func (o openTelemetryBackend) ApplyPolicy(newPolicyData policies.PolicyData, updatePolicy bool) error {
	if !o.policyRepo.Exists(newPolicyData.ID) {
		temporaryFile, err := os.CreateTemp(o.policyConfigDirectory, fmt.Sprintf("otel-%s-config.yml", newPolicyData.ID))
		if err != nil {
			return err
		}
		err = o.addRunner(newPolicyData, temporaryFile.Name())
		if err != nil {
			return err
		}
	} else {
		currentPolicyData, err := o.policyRepo.Get(newPolicyData.ID)
		if err != nil {
			return err
		}
		if currentPolicyData.Version <= newPolicyData.Version {
			dataAsByte := []byte(newPolicyData.Data.(string))
			currentPolicyPath := o.policyConfigDirectory + fmt.Sprintf("otel-%s-config.yml", currentPolicyData.ID)
			o.logger.Info("new policy version received, updating", zap.String("policy_id", newPolicyData.ID), zap.Int32("version", newPolicyData.Version))
			err := os.WriteFile(currentPolicyPath, dataAsByte, os.ModeTemporary)
			if err != nil {
				return err
			}
			err = o.policyRepo.Update(newPolicyData)
			if err != nil {
				return err
			}
		} else {
			o.logger.Info("current policy version is newer than the one being applied, skipping",
				zap.String("policy_id", newPolicyData.ID),
				zap.Int32("current_version", currentPolicyData.Version),
				zap.Int32("incoming_version", newPolicyData.Version))
		}
	}

	return nil
}

func (o openTelemetryBackend) addRunner(policyData policies.PolicyData, policyFilePath string) error {
	policyContext := context.WithValue(o.mainContext, "policy_id", policyData.ID)
	executable, err := memexec.New(openTelemtryContribBinary)
	if err != nil {
		return err
	}
	defer func(executable *memexec.Exec) {
		err := executable.Close()
		if err != nil {
			o.logger.Error("error closing executable", zap.Error(err))
		}
	}(executable)
	command := executable.CommandContext(policyContext, "--config", policyFilePath)
	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}
	go func(ctx context.Context) {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			o.logger.Info("stderr output",
				zap.String("policy_id", policyData.ID),
				zap.String("line", scanner.Text()))
			if command.Err != nil {
				o.logger.Error("error running command", zap.Error(command.Err))
				ctx.Done()
				return
			}
		}
	}(policyContext)

	return nil
}

func (o openTelemetryBackend) addPolicyControl(policyCtx context.Context, policyID string) {
	o.runningCollectors[policyID] = policyCtx
}

func (o openTelemetryBackend) removePolicyControl(policyID string) {
	o.runningCollectors[policyID].Done()
	delete(o.runningCollectors, policyID)
}

func (o openTelemetryBackend) RemovePolicy(data policies.PolicyData) error {
	//TODO implement me
	panic("implement me")
}
