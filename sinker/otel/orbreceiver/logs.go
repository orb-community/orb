/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package orbreceiver

import (
	"context"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
	"strconv"
	"strings"

	"github.com/mainflux/mainflux/pkg/messaging"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"
	"go.uber.org/zap"
)

type internalLogsReceiver struct {
	plogotlp.UnimplementedGRPCServer
	nextConsumer consumer.Logs
	obsrecv      *receiverhelper.ObsReport
}

func (r *OrbReceiver) MessageLogsInbound(msg messaging.Message) error {
	go func() {
		r.cfg.Logger.Debug("received agent message",
			zap.String("subtopic", msg.Subtopic),
			zap.String("channel", msg.Channel),
			zap.String("protocol", msg.Protocol),
			zap.Int64("created", msg.Created),
			zap.String("publisher", msg.Publisher))
		r.cfg.Logger.Info("received log message, pushing to kafka exporter")
		size := len(msg.Payload)
		decompressedPayload := r.DecompressBrotli(msg.Payload)
		lr, err := r.encoder.unmarshalLogsRequest(decompressedPayload)
		if err != nil {
			r.cfg.Logger.Error("error during unmarshalling, skipping message", zap.Error(err))
			return
		}

		r.sinkerService.IncrementMessageCounter(msg.Publisher, msg.Subtopic, msg.Channel, msg.Protocol)

		if lr.Logs().ResourceLogs().Len() == 0 || lr.Logs().ResourceLogs().At(0).ScopeLogs().Len() == 0 {
			r.cfg.Logger.Info("No data information from logs request")
			return
		}

		scopes := lr.Logs().ResourceLogs().At(0).ScopeLogs()
		for i := 0; i < scopes.Len(); i++ {
			r.ProccessLogsContext(scopes.At(i), msg.Channel, size)
		}
	}()
	return nil
}

func (r *OrbReceiver) ProccessLogsContext(scope plog.ScopeLogs, channel string, size int) {
	// Extract Datasets
	attrDataset, ok := scope.Scope().Attributes().Get("dataset_ids")
	if !ok {
		r.cfg.Logger.Info("No datasetIDs information on logs scope attributes")
		return
	}
	datasets := attrDataset.AsString()
	if datasets == "" {
		r.cfg.Logger.Info("datasetIDs information on logs is empty")
		return
	}
	datasetIDs := strings.Split(datasets, ",")
	//Extract policyID
	attrPolID, ok := scope.Scope().Attributes().Get("policy_id")
	if !ok {
		r.cfg.Logger.Info("No policyID information on logs scope attributes")
		return
	}
	polID := attrPolID.AsString()
	if polID == "" {
		r.cfg.Logger.Info("policyID information on logs is empty")
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
		scope = r.injectScopeLogsAttribute(scope, k, v)
	}
	r.injectScopeLogsAttribute(scope, "agent", agentPb.AgentName)
	r.injectScopeLogsAttribute(scope, "policy_id", polID)

	sinkIds, err := r.sinkerService.GetSinkIdsFromDatasetIDs(execCtx, agentPb.OwnerID, datasetIDs)
	if err != nil {
		execCancelF()
		r.cfg.Logger.Info("No data extracting log sinks information from datasetIds = " + datasets)
		return
	}
	attributeCtx := context.WithValue(r.ctx, "agent_name", agentPb.AgentName)
	attributeCtx = context.WithValue(attributeCtx, "agent_tags", agentPb.AgentTags)
	attributeCtx = context.WithValue(attributeCtx, "orb_tags", agentPb.OrbTags)
	attributeCtx = context.WithValue(attributeCtx, "agent_groups", agentPb.AgentGroupIDs)
	attributeCtx = context.WithValue(attributeCtx, "agent_ownerID", agentPb.OwnerID)
	for sinkId := range sinkIds {
		if err != nil {
			r.cfg.Logger.Error("error notifying logs sink active, changing state, skipping sink", zap.String("sink-id", sinkId), zap.Error(err))
			continue
		}
		attributeCtx = context.WithValue(attributeCtx, "sink_id", sinkId)
		lr := plog.NewLogs()
		scope.CopyTo(lr.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty())
		lr.ResourceLogs().At(0).Resource().Attributes().PutStr("service.name", agentPb.AgentName)
		lr.ResourceLogs().At(0).Resource().Attributes().PutStr("service.instance.id", polID)
		request := plogotlp.NewExportRequestFromLogs(lr)
		_, err = r.exportLogs(attributeCtx, request)
		if err != nil {
			r.cfg.Logger.Error("error during logs export, skipping sink", zap.Error(err))
			_ = r.cfg.SinkerService.NotifyActiveSink(r.ctx, agentPb.OwnerID, sinkId, "0")
			continue
		} else {
			_ = r.cfg.SinkerService.NotifyActiveSink(r.ctx, agentPb.OwnerID, sinkId, strconv.Itoa(size))
		}
	}
}

// inject attribute on all ScopeLogs records
func (e *OrbReceiver) injectScopeLogsAttribute(logsScope plog.ScopeLogs, attribute string, value string) plog.ScopeLogs {
	logs := logsScope.LogRecords()
	for i := 0; i < logs.Len(); i++ {
		logItem := logs.At(i)
		logItem.Attributes().PutStr(attribute, value)
	}
	return logsScope
}

func (r *OrbReceiver) exportLogs(ctx context.Context, req plogotlp.ExportRequest) (plogotlp.ExportResponse, error) {
	lr := req.Logs()
	recordsCount := lr.LogRecordCount()
	if recordsCount == 0 {
		return plogotlp.NewExportResponse(), nil
	}

	ctx = r.logsReceiver.obsrecv.StartLogsOp(ctx)
	err := r.logsReceiver.nextConsumer.ConsumeLogs(ctx, lr)
	r.logsReceiver.obsrecv.EndLogsOp(ctx, dataFormatProtobuf, recordsCount, err)

	return plogotlp.NewExportResponse(), err
}
