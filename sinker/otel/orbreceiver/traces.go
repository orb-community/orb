/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package orbreceiver

import (
	"context"
	"strings"

	"github.com/mainflux/mainflux/pkg/messaging"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
	"go.uber.org/zap"
)

type internalTracesReceiver struct {
	ptraceotlp.UnimplementedGRPCServer
	nextConsumer consumer.Traces
	obsrecv      *obsreport.Receiver
}

func (r *OrbReceiver) MessageTracesInbound(msg messaging.Message) error {
	go func() {
		r.cfg.Logger.Debug("received agent message",
			zap.String("subtopic", msg.Subtopic),
			zap.String("channel", msg.Channel),
			zap.String("protocol", msg.Protocol),
			zap.Int64("created", msg.Created),
			zap.String("publisher", msg.Publisher))
		r.cfg.Logger.Info("received trace message, pushing to kafka exporter")
		decompressedPayload := r.DecompressBrotli(msg.Payload)
		tr, err := r.encoder.unmarshalTracesRequest(decompressedPayload)
		if err != nil {
			r.cfg.Logger.Error("error during unmarshalling, skipping message", zap.Error(err))
			return
		}

		r.sinkerService.IncrementMessageCounter(msg.Publisher, msg.Subtopic, msg.Channel, msg.Protocol)

		if tr.Traces().ResourceSpans().Len() == 0 || tr.Traces().ResourceSpans().At(0).ScopeSpans().Len() == 0 {
			r.cfg.Logger.Info("No data information from traces request")
			return
		}

		scopes := tr.Traces().ResourceSpans().At(0).ScopeSpans()
		for i := 0; i < scopes.Len(); i++ {
			r.ProccessTracesContext(scopes.At(i), msg.Channel)
		}
	}()
	return nil
}

func (r *OrbReceiver) ProccessTracesContext(scope ptrace.ScopeSpans, channel string) {
	// Extract Datasets
	attrDataset, ok := scope.Scope().Attributes().Get("dataset_ids")
	if !ok {
		r.cfg.Logger.Info("No datasetIDs information on spans scope attributes")
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
		r.cfg.Logger.Info("No policyID information on spans scope attributes")
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
		scope = r.injectScopeSpansAttribute(scope, k, v)
	}
	r.injectScopeSpansAttribute(scope, "agent", agentPb.AgentName)
	r.injectScopeSpansAttribute(scope, "policy_id", polID)

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
		lr := ptrace.NewTraces()
		scope.CopyTo(lr.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty())
		lr.ResourceSpans().At(0).Resource().Attributes().PutStr("service.name", agentPb.AgentName)
		lr.ResourceSpans().At(0).Resource().Attributes().PutStr("service.instance.id", polID)
		request := ptraceotlp.NewExportRequestFromTraces(lr)
		_, err = r.exportTraces(attributeCtx, request)
		if err != nil {
			r.cfg.Logger.Error("error during export, skipping sink", zap.Error(err))
			_ = r.cfg.SinkerService.NotifyActiveSink(r.ctx, agentPb.OwnerID, sinkId, "error", err.Error())
			continue
		}
	}
}

// inject attribute on all ScopeSpans spans
func (e *OrbReceiver) injectScopeSpansAttribute(spanScope ptrace.ScopeSpans, attribute string, value string) ptrace.ScopeSpans {
	spans := spanScope.Spans()
	for i := 0; i < spans.Len(); i++ {
		spanItem := spans.At(i)
		spanItem.Attributes().PutStr(attribute, value)
	}
	return spanScope
}

func (r *OrbReceiver) exportTraces(ctx context.Context, req ptraceotlp.ExportRequest) (ptraceotlp.ExportResponse, error) {
	ts := req.Traces()
	spanCount := ts.SpanCount()
	if spanCount == 0 {
		return ptraceotlp.NewExportResponse(), nil
	}

	ctx = r.tracesReceiver.obsrecv.StartTracesOp(ctx)
	err := r.tracesReceiver.nextConsumer.ConsumeTraces(ctx, ts)
	r.logsReceiver.obsrecv.EndTracesOp(ctx, dataFormatProtobuf, spanCount, err)

	return ptraceotlp.NewExportResponse(), err
}
