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
	"github.com/ns1labs/orb/sinker/otel/orbreceiver/internal/metrics"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
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

	shutdownWG sync.WaitGroup

	settings component.ReceiverCreateSettings
}

// NewOrbReceiver just creates the OpenTelemetry receiver services. It is the caller's
// responsibility to invoke the respective Start*Reception methods as well
// as the various Stop*Reception methods to end it.
func NewOrbReceiver(cfg *Config, settings component.ReceiverCreateSettings) *orbReceiver {
	r := &orbReceiver{
		cfg:      cfg,
		settings: settings,
	}

	return r
}

// Start appends the message channel that Orb-Sinker will deliver the message
func (r *orbReceiver) Start(ctx context.Context, host component.Host) error {
	r.ctx, r.cancelFunc = context.WithCancel(ctx)

	return nil
}

// Shutdown is a method to turn off receiving.
func (r *orbReceiver) Shutdown(ctx context.Context) error {
	r.shutdownWG.Wait()
	defer r.cancelFunc()
	return nil
}

// registerMetricsConsumer creates a go routine that will monitor the channel
func (r *orbReceiver) registerMetricsConsumer(mc consumer.Metrics) error {
	if mc == nil {
		return component.ErrNilNextConsumer
	}
	metricsReceiverCtx, cancelMetricsReceiver := context.WithCancel(r.ctx)
	go func(ctx context.Context, cancelFunc context.CancelFunc) {
		defer cancelFunc()
	LOOP:
		for {
			select {
			case <-ctx.Done():
				close(*r.cfg.MetricsChannel)
				break LOOP
			case message := <-*r.cfg.MetricsChannel:
				r.cfg.Logger.Info("received metric, pushing to exporter")
				mr := pmetricotlp.NewRequest()
				err := mr.UnmarshalJSON(message)
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
	}(metricsReceiverCtx, cancelMetricsReceiver)

	return nil
}
