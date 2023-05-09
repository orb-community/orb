// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package orbreceiver // import "go.opentelemetry.io/collector/receiver/otlpreceiver"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const (
	typeStr = "orb"
)

// NewFactory creates a new OTLP receiver factory.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		CreateDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, component.StabilityLevelAlpha),
		receiver.WithLogs(createLogsReceiver, component.StabilityLevelAlpha),
		receiver.WithTraces(createTracesReceiver, component.StabilityLevelAlpha))
}

// createDefaultConfig creates the default configuration for receiver.
func CreateDefaultConfig() component.Config {
	return &Config{}
}

// createMetricsReceiver creates a metrics receiver based on provided config.
func createMetricsReceiver(
	ctx context.Context,
	set receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	receiverCfg := cfg.(*Config)
	r := NewOrbReceiver(ctx, cfg.(*Config), set)
	err := r.registerMetricsConsumer(consumer)
	if err != nil {
		receiverCfg.Logger.Info("error on register metrics consumer")
		return nil, err
	}
	return r, nil
}

// createLogsReceiver creates a logs receiver based on provided config.
func createLogsReceiver(
	ctx context.Context,
	set receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Logs,
) (receiver.Logs, error) {
	receiverCfg := cfg.(*Config)
	r := NewOrbReceiver(ctx, cfg.(*Config), set)
	err := r.registerLogsConsumer(consumer)
	if err != nil {
		receiverCfg.Logger.Info("error on register logs consumer")
		return nil, err
	}
	return r, nil
}

// createLogsReceiver creates a logs receiver based on provided config.
func createTracesReceiver(
	ctx context.Context,
	set receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Traces,
) (receiver.Traces, error) {
	receiverCfg := cfg.(*Config)
	r := NewOrbReceiver(ctx, cfg.(*Config), set)
	err := r.registerTracesConsumer(consumer)
	if err != nil {
		receiverCfg.Logger.Info("error on register traces consumer")
		return nil, err
	}
	return r, nil
}
