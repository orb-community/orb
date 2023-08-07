package otel

import (
	"context"
	"fmt"
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

type openTelemetryBackend struct {
	logger    *zap.Logger
	startTime time.Time

	binaryPath string

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
func (o openTelemetryBackend) Configure(logger *zap.Logger, repo policies.PolicyRepo, configuration map[string]string, m2 map[string]interface{}) error {
	o.logger = logger
	o.policyRepo = repo

	var prs bool
	if o.binaryPath, prs = configuration["binary"]; !prs {
		o.binaryPath = DefaultBinary
	}

}

func (o openTelemetryBackend) SetCommsClient(agentID string, client *mqtt.Client, baseTopic string) {
	o.mqttClient = client
	otelBaseTopic := strings.Replace(baseTopic, "?", "otlp", 1)
	o.otlpMetricsTopic = fmt.Sprintf("%s/m/%c", otelBaseTopic, agentID[0])
	o.otlpTracesTopic = fmt.Sprintf("%s/t/%c", otelBaseTopic, agentID[0])
	o.otlpLogsTopic = fmt.Sprintf("%s/l/%c", otelBaseTopic, agentID[0])
}

func (o openTelemetryBackend) Version() (string, error) {
	//TODO implement me
	panic("implement me")
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
