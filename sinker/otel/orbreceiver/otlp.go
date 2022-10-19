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

package orbreceiver

import (
	"context"
	"fmt"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/ns1labs/orb/sinker"
	"github.com/ns1labs/orb/sinker/otel/orbreceiver/internal/metrics"
	"go.opentelemetry.io/collector/config"
	"go.uber.org/zap"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
)

const OtelMetricsTopic = "otlp.*.m.>"

// OrbReceiver is the type that exposes Trace and Metrics reception.
type OrbReceiver struct {
	cfg             *Config
	ctx             context.Context
	cancelFunc      context.CancelFunc
	metricsReceiver *metrics.Receiver
	encoder         encoder
	sinkerService   *sinker.SinkerService

	shutdownWG sync.WaitGroup

	settings component.ReceiverCreateSettings
}

// NewOrbReceiver just creates the OpenTelemetry receiver services. It is the caller's
// responsibility to invoke the respective Start*Reception methods as well
// as the various Stop*Reception methods to end it.
func NewOrbReceiver(ctx context.Context, cfg *Config, settings component.ReceiverCreateSettings) *OrbReceiver {
	r := &OrbReceiver{
		ctx:           ctx,
		cfg:           cfg,
		settings:      settings,
		sinkerService: cfg.SinkerService,
	}

	return r
}

// Start appends the message channel that Orb-Sinker will deliver the message
func (r *OrbReceiver) Start(ctx context.Context, _ component.Host) error {
	r.ctx, r.cancelFunc = context.WithCancel(ctx)

	r.encoder = pbEncoder
	return nil
}

// Shutdown is a method to turn off receiving.
func (r *OrbReceiver) Shutdown(ctx context.Context) error {
	r.cfg.Logger.Warn("shutting down orb-receiver")
	defer func() {
		r.cancelFunc()
		ctx.Done()
	}()
	return nil
}

// registerMetricsConsumer creates a go routine that will monitor the channel
func (r *OrbReceiver) registerMetricsConsumer(mc consumer.Metrics) error {
	if mc == nil {
		return component.ErrNilNextConsumer
	}
	if r.ctx == nil {
		r.cfg.Logger.Warn("error context is nil, using background")
		r.ctx = context.Background()
	}
	r.metricsReceiver = metrics.New(config.NewComponentIDWithName("otlp", "metrics"), mc, r.settings)
	otelTopic := fmt.Sprintf("channels.*.%s", OtelMetricsTopic)
	if err := r.cfg.PubSub.Subscribe(otelTopic, r.MessageInbound); err != nil {
		return err
	}
	r.cfg.Logger.Info("started otel metrics consumer", zap.String("otel-topic", otelTopic))

	return nil
}

func (r *OrbReceiver) MessageInbound(msg messaging.Message) error {
	go func() {
		r.cfg.Logger.Debug("received agent message",
			zap.String("subtopic", msg.Subtopic),
			zap.String("channel", msg.Channel),
			zap.String("protocol", msg.Protocol),
			zap.Int64("created", msg.Created),
			zap.String("publisher", msg.Publisher))
		r.cfg.Logger.Info("received metric message, pushing to exporter")
		mr, err := r.encoder.unmarshalMetricsRequest(msg.Payload)
		if err != nil {
			r.cfg.Logger.Error("error during unmarshalling, skipping message", zap.Error(err))
			return
		}
		// Add tags in Context
		execCtx, execCancelF := context.WithCancel(r.ctx)
		defer execCancelF()
		agentPb, err := r.sinkerService.ExtractAgent(execCtx, msg.Channel)
		if err != nil {
			execCancelF()
			r.cfg.Logger.Error("error during extracting agent information from fleet", zap.Error(err))
			return
		}
		attributeCtx := context.WithValue(r.ctx, "agent-name", agentPb.AgentName)
		attributeCtx = context.WithValue(attributeCtx, "", "")
		_, err = r.metricsReceiver.Export(attributeCtx, mr)
		if err != nil {
			r.cfg.Logger.Error("error during export, skipping message", zap.Error(err))
			return
		}
	}()
	return nil
}
