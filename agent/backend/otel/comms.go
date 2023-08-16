package otel

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
)

func (o *openTelemetryBackend) SetCommsClient(agentID string, client *mqtt.Client, baseTopic string) {
	o.mqttClient = client
	otelBaseTopic := strings.Replace(baseTopic, "?", "otlp", 1)
	o.otlpMetricsTopic = fmt.Sprintf("%s/m/%c", otelBaseTopic, agentID[0])
	o.otlpTracesTopic = fmt.Sprintf("%s/t/%c", otelBaseTopic, agentID[0])
	o.otlpLogsTopic = fmt.Sprintf("%s/l/%c", otelBaseTopic, agentID[0])
}
