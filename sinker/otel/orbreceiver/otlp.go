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
	"github.com/ns1labs/orb/sinker/otel/orbreceiver/internal/metrics"
	"go.uber.org/zap"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
)

// orbReceiver is the type that exposes Trace and Metrics reception.
type orbReceiver struct {
	cfg             *Config
	ctx             context.Context
	cancelFunc      context.CancelFunc
	metricsReceiver *metrics.Receiver
	encoder         encoder

	shutdownWG sync.WaitGroup

	settings component.ReceiverCreateSettings
}

// NewOrbReceiver just creates the OpenTelemetry receiver services. It is the caller's
// responsibility to invoke the respective Start*Reception methods as well
// as the various Stop*Reception methods to end it.
func NewOrbReceiver(ctx context.Context, cfg *Config, settings component.ReceiverCreateSettings) *orbReceiver {
	r := &orbReceiver{
		ctx:      ctx,
		cfg:      cfg,
		settings: settings,
	}

	return r
}

// Start appends the message channel that Orb-Sinker will deliver the message
func (r *orbReceiver) Start(ctx context.Context, _ component.Host) error {
	r.ctx, r.cancelFunc = context.WithCancel(ctx)
	r.encoder = jsEncoder
	return nil
}

// Shutdown is a method to turn off receiving.
func (r *orbReceiver) Shutdown(ctx context.Context) error {
	r.cfg.Logger.Warn("shutting down orb-receiver")
	defer func() {
		r.cancelFunc()
		ctx.Done()
	}()
	return nil
}

// registerMetricsConsumer creates a go routine that will monitor the channel
func (r *orbReceiver) registerMetricsConsumer(mc consumer.Metrics) error {
	if mc == nil {
		return component.ErrNilNextConsumer
	}
	if r.ctx == nil {
		r.cfg.Logger.Warn("error context is nil, using background")
		r.ctx = context.Background()
	}
	metricsReceiverCtx, cancelMetricsReceiver := context.WithCancel(r.ctx)
	r.cfg.Logger.Info("registering go routine to read new messages from orb-sinker")
	go func(ctx context.Context, cancelFunc context.CancelFunc, logger *zap.Logger) {
		defer cancelFunc()
		logger.Info("started routine to listen to channel.")
	LOOP:
		for {
			select {
			case <-ctx.Done():
				logger.Warn("closing receiver routine.")
				close(r.cfg.MetricsChannel)
				break LOOP
			case message := <-r.cfg.MetricsChannel:
				r.cfg.Logger.Info("received metric message, pushing to exporter")
				mr, err := r.encoder.unmarshalMetricsRequest(message)
				if err != nil {
					r.cfg.Logger.Error("error during unmarshalling, skipping message", zap.Error(err))
					continue
				}
				_, err = r.metricsReceiver.Export(ctx, mr)
				if err != nil {
					r.cfg.Logger.Error("error during export, skipping message", zap.Error(err))
					continue
				}
			}
		}
	}(metricsReceiverCtx, cancelMetricsReceiver, r.cfg.Logger)

	return nil
}
