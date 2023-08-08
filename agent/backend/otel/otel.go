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
	policyRepo policies.PolicyRepo
	agentTags  map[string]string

	// Context for controlling the context cancellation
	startCtx           context.Context
	mainCancelFunction context.CancelFunc

	// MQTT Config for OTEL MQTT Exporter
	mqttConfig config.MQTTConfig

	mqttClient *mqtt.Client

	otlpMetricsTopic string
	otlpTracesTopic  string
	otlpLogsTopic    string

	otelReceiverHost string
	otelReceiverPort int
	receiver         receiver.Metrics
	exporter         exporter.Metrics
}

// Configure initializes the backend with the given configuration
func (o openTelemetryBackend) Configure(logger *zap.Logger, repo policies.PolicyRepo, configuration map[string]string, openTelemetryConfiguration map[string]interface{}) error {
	o.logger = logger
	o.policyRepo = repo

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	//TODO implement me
	panic("implement me")
}

func (o openTelemetryBackend) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
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

func (o openTelemetryBackend) ApplyPolicy(data policies.PolicyData, updatePolicy bool) error {
	//TODO implement me
	panic("implement me")
}

func (o openTelemetryBackend) RemovePolicy(data policies.PolicyData) error {
	//TODO implement me
	panic("implement me")
}
