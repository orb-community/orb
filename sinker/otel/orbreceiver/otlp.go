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
	"io/ioutil"
	"strings"

	"sync"

	"github.com/andybalholm/brotli"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/ns1labs/orb/sinker/otel/bridgeservice"
	"github.com/ns1labs/orb/sinker/otel/orbreceiver/internal/metrics"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"go.uber.org/zap"
)

const OtelMetricsTopic = "otlp.*.m.>"

// OrbReceiver is the type that exposes Trace and Metrics reception.
type OrbReceiver struct {
	cfg             *Config
	ctx             context.Context
	cancelFunc      context.CancelFunc
	metricsReceiver *metrics.Receiver
	encoder         encoder
	sinkerService   *bridgeservice.SinkerOtelBridgeService

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

func (r *OrbReceiver) decompressBrotli(data []byte) []byte {
	rdata := bytes.NewReader(data)
	rec := brotli.NewReader(rdata)
	s, _ := ioutil.ReadAll(rec)
	return []byte(s)
}

// extractAttribute extract attribute from metricsRequest metrics
func (r *OrbReceiver) extractAttribute(metricsRequest pmetricotlp.Request, attribute string) string {
	if metricsRequest.Metrics().ResourceMetrics().Len() > 0 {
		if metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().Len() > 0 {
			metrics := metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()
			for i := 0; i < metrics.Len(); i++ {
				metricItem := metrics.At(i)
				switch metricItem.Type() {
				case pmetric.MetricTypeGauge:
					if metricItem.Gauge().DataPoints().Len() > 0 {
						p, ok := metricItem.Gauge().DataPoints().At(0).Attributes().Get(attribute)
						if ok {
							return p.AsString()
						}
					}
				case pmetric.MetricTypeHistogram:
					if metricItem.Histogram().DataPoints().Len() > 0 {
						p, ok := metricItem.Histogram().DataPoints().At(0).Attributes().Get(attribute)
						if ok {
							return p.AsString()
						}
					}
				case pmetric.MetricTypeSum:
					if metricItem.Sum().DataPoints().Len() > 0 {
						p, ok := metricItem.Sum().DataPoints().At(0).Attributes().Get(attribute)
						if ok {
							return p.AsString()
						}
					}
				case pmetric.MetricTypeSummary:
					if metricItem.Summary().DataPoints().Len() > 0 {
						p, ok := metricItem.Summary().DataPoints().At(0).Attributes().Get(attribute)
						if ok {
							return p.AsString()
						}
					}
				case pmetric.MetricTypeExponentialHistogram:
					if metricItem.ExponentialHistogram().DataPoints().Len() > 0 {
						p, ok := metricItem.ExponentialHistogram().DataPoints().At(0).Attributes().Get(attribute)
						if ok {
							return p.AsString()
						}
					}
				}
			}
		}
	}
	return ""
}

// inject attribute on all metricsRequest metrics
func (r *OrbReceiver) injectAttribute(metricsRequest pmetricotlp.Request, attribute string, value string) pmetricotlp.Request {
	metrics := metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()
	for i := 0; i < metrics.Len(); i++ {
		metricItem := metrics.At(i)
		switch metricItem.Type() {
		case pmetric.MetricTypeExponentialHistogram:
			metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(i).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr(attribute, value)
		case pmetric.MetricTypeGauge:
			metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(i).Gauge().DataPoints().At(0).Attributes().PutStr(attribute, value)
		case pmetric.MetricTypeHistogram:
			metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(i).Histogram().DataPoints().At(0).Attributes().PutStr(attribute, value)
		case pmetric.MetricTypeSum:
			metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(i).Sum().DataPoints().At(0).Attributes().PutStr(attribute, value)
		case pmetric.MetricTypeSummary:
			metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(i).Summary().DataPoints().At(0).Attributes().PutStr(attribute, value)
		default:
			r.cfg.Logger.Error("Unknown metric type: " + metricItem.Type().String())
		}
	}
	return metricsRequest
}

// delete attribute on all metricsRequest metrics
func (r *OrbReceiver) deleteAttribute(metricsRequest pmetricotlp.Request, attribute string) pmetricotlp.Request {
	if metricsRequest.Metrics().ResourceMetrics().Len() > 0 {
		if metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().Len() > 0 {
			metrics := metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()
			for i := 0; i < metrics.Len(); i++ {
				metricItem := metrics.At(i)
				switch metricItem.Type() {
				case pmetric.MetricTypeExponentialHistogram:
					metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(i).ExponentialHistogram().DataPoints().At(0).Attributes().Remove(attribute)
				case pmetric.MetricTypeGauge:
					metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(i).Gauge().DataPoints().At(0).Attributes().Remove(attribute)
				case pmetric.MetricTypeHistogram:
					metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(i).Histogram().DataPoints().At(0).Attributes().Remove(attribute)
				case pmetric.MetricTypeSum:
					metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(i).Sum().DataPoints().At(0).Attributes().Remove(attribute)
				case pmetric.MetricTypeSummary:
					metricsRequest.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(i).Summary().DataPoints().At(0).Attributes().Remove(attribute)
				default:
					r.cfg.Logger.Error("Unknown metric type: " + metricItem.Type().String())
				}
			}
		} else {
			r.cfg.Logger.Error("Unable to delete attribute, ScopeMetrics length 0")
		}
	} else {
		r.cfg.Logger.Error("Unable to delete attribute, ResourceMetrics length 0")
	}
	return metricsRequest
}

func (r *OrbReceiver) MessageInbound(msg messaging.Message) error {
	go func() {
		r.cfg.Logger.Debug("received agent message",
			zap.String("subtopic", msg.Subtopic),
			zap.String("channel", msg.Channel),
			zap.String("protocol", msg.Protocol),
			zap.Int64("created", msg.Created),
			zap.String("publisher", msg.Publisher))
		r.cfg.Logger.Info("received metric message, pushing to kafka exporter")
		decompressedPayload := r.decompressBrotli(msg.Payload)
		mr, err := r.encoder.unmarshalMetricsRequest(decompressedPayload)
		if err != nil {
			r.cfg.Logger.Error("error during unmarshalling, skipping message", zap.Error(err))
			return
		}

		// Extract Datasets
		datasets := r.extractAttribute(mr, "dataset_ids")
		if datasets == "" {
			r.cfg.Logger.Info("No data extracting datasetIDs information from metrics request")
			return
		}
		datasetIDs := strings.Split(datasets, ",")
		
		// Delete datasets_ids and policy_ids from metricsRequest
		mr = r.deleteAttribute(mr, "dataset_ids")
		mr = r.deleteAttribute(mr, "policy_ids")

		// Add tags in Context
		execCtx, execCancelF := context.WithCancel(r.ctx)
		defer execCancelF()
		agentPb, err := r.sinkerService.ExtractAgent(execCtx, msg.Channel)
		if err != nil {
			execCancelF()
			r.cfg.Logger.Info("No data extracting agent information from fleet")
			return
		}
		sinkIds, err := r.sinkerService.GetSinkIdsFromDatasetIDs(execCtx, agentPb.OwnerID, datasetIDs)
		if err != nil {
			execCancelF()
			r.cfg.Logger.Info("No data extracting sinks information from datasetIds = " + datasets)
			return
		}
		attributeCtx := context.WithValue(r.ctx, "agent_name", agentPb.AgentName)
		attributeCtx = context.WithValue(attributeCtx, "agent_tags", agentPb.AgentTags)
		attributeCtx = context.WithValue(attributeCtx, "orb_tags", agentPb.OrbTags)
		attributeCtx = context.WithValue(attributeCtx, "agent_groups", agentPb.AgentGroupIDs)
		attributeCtx = context.WithValue(attributeCtx, "agent_ownerID", agentPb.OwnerID)
		for sinkId, _ := range sinkIds {
			err := r.cfg.SinkerService.NotifyActiveSink(r.ctx, agentPb.OwnerID, sinkId, "active", "")
			if err != nil {
				r.cfg.Logger.Error("error notifying sink active, changing state, skipping sink", zap.String("sink-id", sinkId), zap.Error(err))
				continue
			}
			attributeCtx = context.WithValue(attributeCtx, "sink_id", sinkId)
			_, err = r.metricsReceiver.Export(attributeCtx, mr)
			if err != nil {
				r.cfg.Logger.Error("error during export, skipping sink", zap.Error(err))
				_ = r.cfg.SinkerService.NotifyActiveSink(r.ctx, agentPb.OwnerID, sinkId, "error", err.Error())
				continue
			}
		}
	}()
	return nil
}
