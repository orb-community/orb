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
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/andybalholm/brotli"
	"github.com/orb-community/orb/sinker/otel/bridgeservice"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

const (
	OtelMetricsTopic   = "otlp.*.m.>"
	OtelLogsTopic      = "otlp.*.l.>"
	OtelTraceTopic     = "otlp.*.t.>"
	dataFormatProtobuf = "protobuf"
)

// OrbReceiver is the type that exposes Trace and Metrics reception.
type OrbReceiver struct {
	cfg             *Config
	ctx             context.Context
	cancelFunc      context.CancelFunc
	metricsReceiver internalMetricsReceiver
	logsReceiver    internalLogsReceiver
	tracesReceiver  internalTracesReceiver
	encoder         encoder
	sinkerService   *bridgeservice.SinkerOtelBridgeService

	shutdownWG sync.WaitGroup

	settings receiver.CreateSettings
}

// NewOrbReceiver just creates the OpenTelemetry receiver services. It is the caller's
// responsibility to invoke the respective Start*Reception methods as well
// as the various Stop*Reception methods to end it.
func NewOrbReceiver(ctx context.Context, cfg *Config, settings receiver.CreateSettings) *OrbReceiver {
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

func (r *OrbReceiver) DecompressBrotli(data []byte) []byte {
	rdata := bytes.NewReader(data)
	rec := brotli.NewReader(rdata)
	s, _ := io.ReadAll(rec)
	return []byte(s)
}

func (r *OrbReceiver) registerMetricsConsumer(mc consumer.Metrics) error {
	if mc == nil {
		return component.ErrNilNextConsumer
	}
	if r.ctx == nil {
		r.cfg.Logger.Warn("error context is nil, using background")
		r.ctx = context.Background()
	}
	var err error
	obsrecv, err := obsreport.NewReceiver(obsreport.ReceiverSettings{
		ReceiverID:             component.NewIDWithName("otlp", "metrics"),
		Transport:              "grpc",
		ReceiverCreateSettings: r.settings,
	})
	if err != nil {
		return err
	}
	r.metricsReceiver = internalMetricsReceiver{
		nextConsumer: mc,
		obsrecv:      obsrecv,
	}
	otelTopic := fmt.Sprintf("channels.*.%s", OtelMetricsTopic)
	if err = r.cfg.PubSub.Subscribe(otelTopic, r.MessageMetricsInbound); err != nil {
		return err
	}
	r.cfg.Logger.Info("started otel metrics consumer", zap.String("otel-topic", otelTopic))

	return nil
}

func (r *OrbReceiver) registerLogsConsumer(lc consumer.Logs) error {
	if lc == nil {
		return component.ErrNilNextConsumer
	}
	if r.ctx == nil {
		r.cfg.Logger.Warn("error context is nil, using background")
		r.ctx = context.Background()
	}
	var err error
	obsrecv, err := obsreport.NewReceiver(obsreport.ReceiverSettings{
		ReceiverID:             component.NewIDWithName("otlp", "logs"),
		Transport:              "grpc",
		ReceiverCreateSettings: r.settings,
	})
	if err != nil {
		return err
	}
	r.logsReceiver = internalLogsReceiver{
		nextConsumer: lc,
		obsrecv:      obsrecv,
	}
	otelTopic := fmt.Sprintf("channels.*.%s", OtelLogsTopic)
	if err = r.cfg.PubSub.Subscribe(otelTopic, r.MessageLogsInbound); err != nil {
		return err
	}
	r.cfg.Logger.Info("started otel logs consumer", zap.String("otel-topic", otelTopic))

	return nil
}

func (r *OrbReceiver) registerTracesConsumer(tc consumer.Traces) error {
	if tc == nil {
		return component.ErrNilNextConsumer
	}
	if r.ctx == nil {
		r.cfg.Logger.Warn("error context is nil, using background")
		r.ctx = context.Background()
	}
	var err error
	obsrecv, err := obsreport.NewReceiver(obsreport.ReceiverSettings{
		ReceiverID:             component.NewIDWithName("otlp", "traces"),
		Transport:              "grpc",
		ReceiverCreateSettings: r.settings,
	})
	if err != nil {
		return err
	}
	r.tracesReceiver = internalTracesReceiver{
		nextConsumer: tc,
		obsrecv:      obsrecv,
	}
	otelTopic := fmt.Sprintf("channels.*.%s", OtelTraceTopic)
	if err = r.cfg.PubSub.Subscribe(otelTopic, r.MessageTracesInbound); err != nil {
		return err
	}
	r.cfg.Logger.Info("started otel traces consumer", zap.String("otel-topic", otelTopic))

	return nil
}
