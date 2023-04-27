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
	"strings"
	"sync"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/orb-community/orb/sinker/otel/bridgeservice"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

const OtelMetricsTopic = "otlp.*.m.>"

// OrbReceiver is the type that exposes Trace and Metrics reception.
type OrbReceiver struct {
	cfg             *Config
	ctx             context.Context
	cancelFunc      context.CancelFunc
	metricsReceiver *internalReceiver
	encoder         encoder
	metricsConsumer consumer.Metrics
	sinkerService   *bridgeservice.SinkerOtelBridgeService
	obsrepGRPC      *obsreport.Receiver

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

// registerMetricsConsumer creates a go routine that will monitor the channel
func (r *OrbReceiver) registerMetricsConsumer(mc consumer.Metrics) error {
	if mc == nil {
		return component.ErrNilNextConsumer
	}
	if r.ctx == nil {
		r.cfg.Logger.Warn("error context is nil, using background")
		r.ctx = context.Background()
	}
	var err error
	r.obsrepGRPC, err = obsreport.NewReceiver(obsreport.ReceiverSettings{
		ReceiverID:             component.NewIDWithName("otlp", "metrics"),
		Transport:              "grpc",
		ReceiverCreateSettings: r.settings,
	})
	if err != nil {
		return err
	}
	r.metricsReceiver = InternalReceiverNew(mc, r.obsrepGRPC)
	otelTopic := fmt.Sprintf("channels.*.%s", OtelMetricsTopic)
	if err = r.cfg.PubSub.Subscribe(otelTopic, r.MessageInbound); err != nil {
		return err
	}
	r.cfg.Logger.Info("started otel metrics consumer", zap.String("otel-topic", otelTopic))

	return nil
}

func (r *OrbReceiver) decompressBrotli(data []byte) []byte {
	rdata := bytes.NewReader(data)
	rec := brotli.NewReader(rdata)
	s, _ := io.ReadAll(rec)
	return []byte(s)
}

// extractAttribute extract attribute from metricsScope metrics
func (r *OrbReceiver) extractScopeAttribute(metricsScope pmetric.ScopeMetrics, attribute string) string {
	metrics := metricsScope.Metrics()
	if metrics.Len() > 0 {
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
	return ""
}

// inject attribute on all ScopeMetrics metrics
func (r *OrbReceiver) injectScopeAttribute(metricsScope pmetric.ScopeMetrics, attribute string, value string) pmetric.ScopeMetrics {
	metrics := metricsScope.Metrics()
	for i := 0; i < metrics.Len(); i++ {
		metricItem := metrics.At(i)

		switch metricItem.Type() {
		case pmetric.MetricTypeExponentialHistogram:
			for i := 0; i < metricItem.ExponentialHistogram().DataPoints().Len(); i++ {
				metricItem.ExponentialHistogram().DataPoints().At(i).Attributes().PutStr(attribute, value)
			}
		case pmetric.MetricTypeGauge:
			for i := 0; i < metricItem.Gauge().DataPoints().Len(); i++ {
				metricItem.Gauge().DataPoints().At(i).Attributes().PutStr(attribute, value)
			}
		case pmetric.MetricTypeHistogram:
			for i := 0; i < metricItem.Histogram().DataPoints().Len(); i++ {
				metricItem.Histogram().DataPoints().At(i).Attributes().PutStr(attribute, value)
			}
		case pmetric.MetricTypeSum:
			for i := 0; i < metricItem.Sum().DataPoints().Len(); i++ {
				metricItem.Sum().DataPoints().At(i).Attributes().PutStr(attribute, value)
			}
		case pmetric.MetricTypeSummary:
			for i := 0; i < metricItem.Summary().DataPoints().Len(); i++ {
				metricItem.Summary().DataPoints().At(i).Attributes().PutStr(attribute, value)
			}
		default:
			continue
		}
	}
	return metricsScope
}

// replace ScopeMetrics metrics timestamp
func (r *OrbReceiver) replaceScopeTimestamp(metricsScope pmetric.ScopeMetrics, ts pcommon.Timestamp) pmetric.ScopeMetrics {
	metricsList := metricsScope.Metrics()
	for i3 := 0; i3 < metricsList.Len(); i3++ {
		metricItem := metricsList.At(i3)
		switch metricItem.Type() {
		case pmetric.MetricTypeExponentialHistogram:
			for i := 0; i < metricItem.ExponentialHistogram().DataPoints().Len(); i++ {
				metricItem.ExponentialHistogram().DataPoints().At(i).SetTimestamp(ts)
			}
		case pmetric.MetricTypeGauge:
			for i := 0; i < metricItem.Gauge().DataPoints().Len(); i++ {
				metricItem.Gauge().DataPoints().At(i).SetTimestamp(ts)
			}
		case pmetric.MetricTypeHistogram:
			for i := 0; i < metricItem.Histogram().DataPoints().Len(); i++ {
				metricItem.Histogram().DataPoints().At(i).SetTimestamp(ts)
			}
		case pmetric.MetricTypeSum:
			for i := 0; i < metricItem.Sum().DataPoints().Len(); i++ {
				metricItem.Sum().DataPoints().At(i).SetTimestamp(ts)
			}
		case pmetric.MetricTypeSummary:
			for i := 0; i < metricItem.Summary().DataPoints().Len(); i++ {
				metricItem.Summary().DataPoints().At(i).SetTimestamp(ts)
			}
		}
	}
	return metricsScope
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

		r.sinkerService.IncreamentMessageCounter(msg.Publisher, msg.Subtopic, msg.Channel, msg.Protocol)

		if mr.Metrics().ResourceMetrics().Len() == 0 || mr.Metrics().ResourceMetrics().At(0).ScopeMetrics().Len() == 0 {
			r.cfg.Logger.Info("No data information from metrics request")
			return
		}

		scopes := mr.Metrics().ResourceMetrics().At(0).ScopeMetrics()
		for i := 0; i < scopes.Len(); i++ {
			r.ProccessPolicyContext(scopes.At(i), msg.Channel)
		}
	}()
	return nil
}

func (r *OrbReceiver) ProccessPolicyContext(scope pmetric.ScopeMetrics, channel string) {
	// Extract Datasets
	attrDataset, ok := scope.Scope().Attributes().Get("dataset_ids")
	if !ok {
		r.cfg.Logger.Info("No datasetIDs information on metrics scope attributes")
		return
	}
	datasets := attrDataset.AsString()
	if datasets == "" {
		r.cfg.Logger.Info("datasetIDs information is empty")
		return
	}
	datasetIDs := strings.Split(datasets, ",")
	//Extract policyID
	attrPolID, ok := scope.Scope().Attributes().Get("policy_id")
	if !ok {
		r.cfg.Logger.Info("No policyID information on metrics scope attributes")
		return
	}
	polID := attrPolID.AsString()
	if polID == "" {
		r.cfg.Logger.Info("policyID information is empty")
		return
	}
	// Delete datasets_ids and policy_ids from scope attributes
	scope.Scope().Attributes().Clear()

	// Add tags in Context
	execCtx, execCancelF := context.WithCancel(r.ctx)
	defer execCancelF()
	agentPb, err := r.sinkerService.ExtractAgent(execCtx, channel)
	if err != nil {
		execCancelF()
		r.cfg.Logger.Info("No data extracting agent information from fleet")
		return
	}
	for k, v := range agentPb.OrbTags {
		scope = r.injectScopeAttribute(scope, k, v)
	}

	scope = r.replaceScopeTimestamp(scope, pcommon.NewTimestampFromTime(time.Now()))
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
	for sinkId := range sinkIds {
		err := r.cfg.SinkerService.NotifyActiveSink(r.ctx, agentPb.OwnerID, sinkId, "active", "")
		if err != nil {
			r.cfg.Logger.Error("error notifying sink active, changing state, skipping sink", zap.String("sink-id", sinkId), zap.Error(err))
			continue
		}
		attributeCtx = context.WithValue(attributeCtx, "sink_id", sinkId)
		mr := pmetric.NewMetrics()
		scope.CopyTo(mr.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty())
		mr.ResourceMetrics().At(0).Resource().Attributes().PutStr("service.name", agentPb.AgentName)
		mr.ResourceMetrics().At(0).Resource().Attributes().PutStr("service.instance.id", polID)
		request := pmetricotlp.NewExportRequestFromMetrics(mr)
		_, err = r.metricsReceiver.Export(attributeCtx, request)
		if err != nil {
			r.cfg.Logger.Error("error during export, skipping sink", zap.Error(err))
			_ = r.cfg.SinkerService.NotifyActiveSink(r.ctx, agentPb.OwnerID, sinkId, "error", err.Error())
			continue
		}
	}
}
