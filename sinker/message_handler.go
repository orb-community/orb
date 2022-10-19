package sinker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/pb"
	"github.com/ns1labs/orb/pkg/types"
	pb2 "github.com/ns1labs/orb/policies/pb"
	"github.com/ns1labs/orb/sinker/backend"
	"github.com/ns1labs/orb/sinker/config"
	"github.com/ns1labs/orb/sinker/prometheus"
	pb3 "github.com/ns1labs/orb/sinks/pb"
	"go.uber.org/zap"
	"strings"
	"time"
)

func (svc SinkerService) remoteWriteToPrometheus(tsList prometheus.TSList, ownerID string, sinkID string) error {
	cfgRepo, err := svc.sinkerCache.Get(ownerID, sinkID)
	if err != nil {
		svc.logger.Error("unable to retrieve the sink config", zap.Error(err))
		return err
	}

	cfg := prometheus.NewConfig(
		prometheus.WriteURLOption(cfgRepo.Url),
	)

	promClient, err := prometheus.NewClient(cfg)
	if err != nil {
		svc.logger.Error("unable to construct client", zap.Error(err))
		return err
	}

	var headers = make(map[string]string)
	headers["Authorization"] = svc.encodeBase64(cfgRepo.User, cfgRepo.Password)
	result, writeErr := promClient.WriteTimeSeries(context.Background(), tsList,
		prometheus.WriteOptions{Headers: headers})
	if err := error(writeErr); err != nil {
		if cfgRepo.State != config.Error || cfgRepo.Msg != fmt.Sprint(err) {
			cfgRepo.State = config.Error
			cfgRepo.Msg = fmt.Sprint(err)
			cfgRepo.LastRemoteWrite = time.Now()
			err := svc.sinkerCache.Edit(cfgRepo)
			if err != nil {
				svc.logger.Error("error during update sink cache", zap.Error(err))
				return err
			}
		}

		svc.logger.Error("remote write error", zap.String("sink_id", sinkID), zap.Error(err))
		return err
	}

	svc.logger.Debug("successful sink", zap.Int("payload_size_b", result.PayloadSize), zap.String("sink_id", sinkID), zap.String("url", cfgRepo.Url), zap.String("user", cfgRepo.User))

	if cfgRepo.State != config.Active {
		cfgRepo.State = config.Active
		cfgRepo.Msg = ""
		cfgRepo.LastRemoteWrite = time.Now()
		err := svc.sinkerCache.Edit(cfgRepo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (svc SinkerService) encodeBase64(user string, password string) string {
	defer func(t time.Time) {
		svc.logger.Debug("encodeBase64 took", zap.String("execution", time.Since(t).String()))
	}(time.Now())
	sEnc := base64.URLEncoding.EncodeToString([]byte(user + ":" + password))
	return fmt.Sprintf("Basic %s", sEnc)
}

func (svc SinkerService) handleMetrics(agentID string, channelID string, subtopic string, payload []byte) error {

	// find backend to send it to
	beName := strings.Split(subtopic, ".")
	if len(beName) < 3 || beName[0] != "be" || beName[2] != "m" {
		return errors.New(fmt.Sprintf("invalid subtopic, ignoring: %s", subtopic))
	}
	if !backend.HaveBackend(beName[1]) {
		return errors.New(fmt.Sprintf("unknown agent backend, ignoring: %s", beName[1]))
	}
	be := backend.GetBackend(beName[1])

	// unpack metrics RPC
	var versionCheck fleet.SchemaVersionCheck
	if err := json.Unmarshal(payload, &versionCheck); err != nil {
		return fleet.ErrSchemaMalformed
	}
	if versionCheck.SchemaVersion != fleet.CurrentRPCSchemaVersion {
		return fleet.ErrSchemaVersion
	}
	var rpc fleet.RPC
	if err := json.Unmarshal(payload, &rpc); err != nil {
		return fleet.ErrSchemaMalformed
	}
	if rpc.Func != fleet.AgentMetricsRPCFunc {
		return errors.New(fmt.Sprintf("unexpected RPC function: %s", rpc.Func))
	}
	var metricsRPC fleet.AgentMetricsRPC
	if err := json.Unmarshal(payload, &metricsRPC); err != nil {
		return fleet.ErrSchemaMalformed
	}

	agentPb, err2 := svc.ExtractAgent(channelID)
	if err2 != nil {
		return err2
	}

	agentName, _ := types.NewIdentifier(agentPb.AgentName)

	agent := fleet.Agent{
		Name:        agentName,
		MFOwnerID:   agentPb.OwnerID,
		MFThingID:   agentID,
		MFChannelID: channelID,
		OrbTags:     agentPb.OrbTags,
		AgentTags:   agentPb.AgentTags,
	}

	for _, metricsPayload := range metricsRPC.Payload {
		// this payload loop is per policy. each policy has a list of datasets it is associated with, and each dataset may contain multiple sinks
		// however, per policy, we want a unique set of sink IDs as we don't want to send the same metrics twice to the same sink for the same policy
		datasetSinkIDs := make(map[string]bool)
		// first go through the datasets and gather the unique set of sinks we need for this particular policy
		err2 := svc.GetSinks(agent, metricsPayload, datasetSinkIDs)
		if err2 != nil {
			return err2
		}

		// ensure there are sinks
		if len(datasetSinkIDs) == 0 {
			svc.logger.Error("unable to attach any sinks to policy", zap.String("policy_id", metricsPayload.PolicyID), zap.String("agent_id", agentID), zap.String("owner_id", agent.MFOwnerID))
			continue
		}

		// now that we have the sinks, process the metrics for this policy
		tsList, err := be.ProcessMetrics(agentPb, agentID, metricsPayload)
		if err != nil {
			svc.logger.Error("ProcessMetrics failed", zap.String("policy_id", metricsPayload.PolicyID), zap.String("agent_id", agentID), zap.String("owner_id", agent.MFOwnerID), zap.Error(err))
			continue
		}

		// finally, sink this policy
		svc.SinkPolicy(agent, metricsPayload, datasetSinkIDs, tsList)
	}

	return nil
}

func (svc SinkerService) ExtractAgent(ctx context.Context, channelID string) (*pb.AgentInfoRes, error) {
	agentPb, err := svc.fleetClient.RetrieveAgentInfoByChannelID(ctx, &pb.AgentInfoByChannelIDReq{Channel: channelID})
	if err != nil {
		return nil, err
	}
	return agentPb, nil
}

func (svc SinkerService) SinkPolicy(agent fleet.Agent, metricsPayload fleet.AgentMetricsRPCPayload, datasetSinkIDs map[string]bool, tsList []prometheus.TimeSeries) {
	sinkIDList := make([]string, len(datasetSinkIDs))
	i := 0
	for k := range datasetSinkIDs {
		sinkIDList[i] = k
		i++
	}
	svc.logger.Info("sinking agent metric RPC",
		zap.String("owner_id", agent.MFOwnerID),
		zap.String("agent", agent.Name.String()),
		zap.String("policy", metricsPayload.PolicyName),
		zap.String("policy_id", metricsPayload.PolicyID),
		zap.Strings("sinks", sinkIDList))

	for _, id := range sinkIDList {
		err := svc.remoteWriteToPrometheus(tsList, agent.MFOwnerID, id)
		if err != nil {
			svc.logger.Warn(fmt.Sprintf("unable to remote write to sinkID: %s", id), zap.String("policy_id", metricsPayload.PolicyID), zap.String("agent_id", agent.MFThingID), zap.String("owner_id", agent.MFOwnerID), zap.Error(err))
		}

		// send operational metrics
		labels := []string{
			"method", "sinker_payload_size",
			"agent_id", agent.MFThingID,
			"agent", agent.Name.String(),
			"policy_id", metricsPayload.PolicyID,
			"policy", metricsPayload.PolicyName,
			"sink_id", id,
			"owner_id", agent.MFOwnerID,
		}
		svc.requestCounter.With(labels...).Add(1)
		svc.requestGauge.With(labels...).Add(float64(len(metricsPayload.Data)))
	}
}

func (svc SinkerService) GetSinks(agent fleet.Agent, agentMetricsRPCPayload fleet.AgentMetricsRPCPayload, datasetSinkIDs map[string]bool) error {
	for _, ds := range agentMetricsRPCPayload.Datasets {
		if ds == "" {
			svc.logger.Error("malformed agent RPC: empty dataset", zap.String("agent_id", agent.MFThingID), zap.String("owner_id", agent.MFOwnerID))
			continue
		}
		dataset, err := svc.policiesClient.RetrieveDataset(context.Background(), &pb2.DatasetByIDReq{
			DatasetID: ds,
			OwnerID:   agent.MFOwnerID,
		})
		if err != nil {
			svc.logger.Error("unable to retrieve dataset", zap.String("dataset_id", ds), zap.String("owner_id", agent.MFOwnerID), zap.Error(err))
			continue
		}
		for _, sid := range dataset.SinkIds {
			if !svc.sinkerCache.Exists(agent.MFOwnerID, sid) {
				// Use the retrieved sinkID to get the backend config
				sink, err := svc.sinksClient.RetrieveSink(context.Background(), &pb3.SinkByIDReq{
					SinkID:  sid,
					OwnerID: agent.MFOwnerID,
				})
				if err != nil {
					return err
				}

				var data config.SinkConfig
				if err := json.Unmarshal(sink.Config, &data); err != nil {
					return err
				}

				data.SinkID = sid
				data.OwnerID = agent.MFOwnerID
				err = svc.sinkerCache.Add(data)
				if err != nil {
					return err
				}
			}
			datasetSinkIDs[sid] = true
		}
	}
	return nil
}

func (svc SinkerService) handleMsgFromAgent(msg messaging.Message) error {
	inputContext := context.WithValue(context.Background(), "trace-id", uuid.NewString())
	go func(ctx context.Context) {
		defer func(t time.Time) {
			svc.logger.Info("message consumption time", zap.String("execution", time.Since(t).String()))
		}(time.Now())
		// NOTE: we need to consider ALL input from the agent as untrusted, the same as untrusted HTTP API would be
		var payload map[string]interface{}
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			svc.logger.Error("metrics processing failure", zap.Any("trace-id", ctx.Value("trace-id")), zap.Error(err))
			return
		}

		svc.logger.Debug("received agent message",
			zap.String("subtopic", msg.Subtopic),
			zap.String("channel", msg.Channel),
			zap.String("protocol", msg.Protocol),
			zap.Int64("created", msg.Created),
			zap.String("publisher", msg.Publisher))

		labels := []string{
			"method", "handleMsgFromAgent",
			"agent_id", msg.Publisher,
			"subtopic", msg.Subtopic,
			"channel", msg.Channel,
			"protocol", msg.Protocol,
		}
		svc.messageInputCounter.With(labels...).Add(1)

		if len(msg.Payload) > MaxMsgPayloadSize {
			svc.logger.Error("metrics processing failure", zap.Any("trace-id", ctx.Value("trace-id")), zap.Error(ErrPayloadTooBig))
			return
		}

		if err := svc.handleMetrics(msg.Publisher, msg.Channel, msg.Subtopic, msg.Payload); err != nil {
			svc.logger.Error("metrics processing failure", zap.Any("trace-id", ctx.Value("trace-id")), zap.Error(err))
			return
		}
	}(inputContext)

	return nil
}
