/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package orbreceiver

import (
	"context"
	"strings"
	"time"

	"github.com/mainflux/mainflux/pkg/messaging"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"go.uber.org/zap"
)

type internalMetricsReceiver struct {
	pmetricotlp.UnimplementedGRPCServer
	nextConsumer consumer.Metrics
	obsrecv      *obsreport.Receiver
}

func (r *OrbReceiver) MessageMetricsInbound(msg messaging.Message) error {
	go func() {
		r.cfg.Logger.Debug("received agent message",
			zap.String("subtopic", msg.Subtopic),
			zap.String("channel", msg.Channel),
			zap.String("protocol", msg.Protocol),
			zap.Int64("created", msg.Created),
			zap.String("publisher", msg.Publisher))
		r.cfg.Logger.Info("received metric message, pushing to kafka exporter")
		decompressedPayload := r.DecompressBrotli(msg.Payload)
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
			r.ProccessMetricsContext(scopes.At(i), msg.Channel)
		}
	}()
	return nil
}

func (r *OrbReceiver) ProccessMetricsContext(scope pmetric.ScopeMetrics, channel string) {
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
		scope = r.injectScopeMetricsAttribute(scope, k, v)
	}

	scope = r.replaceScopeMetricsTimestamp(scope, pcommon.NewTimestampFromTime(time.Now()))
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
		_, err = r.exportMetrics(attributeCtx, request)
		if err != nil {
			r.cfg.Logger.Error("error during export, skipping sink", zap.Error(err))
			_ = r.cfg.SinkerService.NotifyActiveSink(r.ctx, agentPb.OwnerID, sinkId, "error", err.Error())
			continue
		}
	}
}

// inject attribute on all ScopeMetrics metrics
func (r *OrbReceiver) injectScopeMetricsAttribute(metricsScope pmetric.ScopeMetrics, attribute string, value string) pmetric.ScopeMetrics {
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
func (r *OrbReceiver) replaceScopeMetricsTimestamp(metricsScope pmetric.ScopeMetrics, ts pcommon.Timestamp) pmetric.ScopeMetrics {
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

func (r *OrbReceiver) exportMetrics(ctx context.Context, req pmetricotlp.ExportRequest) (pmetricotlp.ExportResponse, error) {
	md := req.Metrics()
	dataPointCount := md.DataPointCount()
	if dataPointCount == 0 {
		return pmetricotlp.NewExportResponse(), nil
	}

	ctx = r.metricsReceiver.obsrecv.StartMetricsOp(ctx)
	err := r.metricsReceiver.nextConsumer.ConsumeMetrics(ctx, md)
	r.metricsReceiver.obsrecv.EndMetricsOp(ctx, dataFormatProtobuf, dataPointCount, err)

	return pmetricotlp.NewExportResponse(), err
}
