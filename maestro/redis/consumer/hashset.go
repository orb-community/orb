package consumer

import (
	"context"
	"encoding/json"
	"github.com/ns1labs/orb/maestro/config"
	"github.com/ns1labs/orb/pkg/types"
	"go.uber.org/zap"
	"time"
)

const deploymentKey = "orb.sinks.deployment"

func (es eventStore) GetDeploymentEntryFromSinkId(ctx context.Context, sinkId string) (string, error) {
	cmd := es.client.HGet(ctx, deploymentKey, sinkId)
	if err := cmd.Err(); err != nil {
		es.logger.Error("error during redis reading of SinkId", zap.String("sink-id", sinkId), zap.Error(err))
		return "", err
	}
	return cmd.String(), nil
}

// handleSinksDeleteCollector will delete Deployment Entry and force delete otel collector
func (es eventStore) handleSinksDeleteCollector(ctx context.Context, event sinksUpdateEvent) error {
	es.logger.Info("Received maestro DELETE event from sinks ID=" + event.sinkID + ", Owner ID=" + event.owner)
	deployment, err := es.GetDeploymentEntryFromSinkId(ctx, event.sinkID)
	if err != nil {
		es.logger.Error("did not find collector entry for sink", zap.String("sink-id", event.sinkID))
		return err
	}
	err = es.kubecontrol.DeleteOtelCollector(ctx, event.sinkID, deployment)
	if err != nil {
		return err
	}
	es.client.HDel(ctx, deploymentKey, event.sinkID)
	return nil
}

// handleSinksCreateCollector will create Deployment Entry in Redis
func (es eventStore) handleSinksCreateCollector(ctx context.Context, event sinksUpdateEvent) error {
	es.logger.Info("Received event to Create DeploymentEntry from sinks ID=" + event.sinkID + ", Owner ID=" + event.owner)
	sinkUrl := event.config["sink_url"].(string)
	sinkUsername := event.config["username"].(string)
	sinkPassword := event.config["password"].(string)
	err2 := es.CreateDeploymentEntry(ctx, event.sinkID, sinkUrl, sinkUsername, sinkPassword)
	if err2 != nil {
		return err2
	}

	return nil
}

func (es eventStore) CreateDeploymentEntry(ctx context.Context, sinkId, sinkUrl, sinkUsername, sinkPassword string) error {
	deploy, err := config.GetDeploymentJson(sinkId, sinkUrl, sinkUsername, sinkPassword)
	if err != nil {
		es.logger.Error("error trying to get deployment json for sink ID", zap.String("sinkId", sinkId))
		return err
	}
	es.client.HSet(ctx, deploymentKey, sinkId, deploy)
	return nil
}

// handleSinksUpdateCollector will update Deployment Entry in Redis and force update otel collector
func (es eventStore) handleSinksUpdateCollector(ctx context.Context, event sinksUpdateEvent) error {
	es.logger.Info("Received event to Update DeploymentEntry from sinks ID=" + event.sinkID + ", Owner ID=" + event.owner)
	sinkUrl := event.config["sink_url"].(string)
	sinkUsername := event.config["username"].(string)
	sinkPassword := event.config["password"].(string)
	deploy, err := config.GetDeploymentJson(event.sinkID, sinkUrl, sinkUsername, sinkPassword)
	if err != nil {
		es.logger.Error("error trying to get deployment json for sink ID", zap.String("sinkId", event.sinkID))
		return err
	}
	es.client.HSet(ctx, deploymentKey, event.sinkID, deploy)
	err = es.kubecontrol.UpdateOtelCollector(ctx, event.sinkID, deploy)
	if err != nil {
		return err
	}

	return nil
}

func decodeSinksEvent(event map[string]interface{}, operation string) (sinksUpdateEvent, error) {
	val := sinksUpdateEvent{
		sinkID:    read(event, "sink_id", ""),
		owner:     read(event, "owner", ""),
		timestamp: time.Now(),
	}
	if operation != sinksDelete {
		var metadata types.Metadata
		if err := json.Unmarshal([]byte(read(event, "config", "")), &metadata); err != nil {
			return sinksUpdateEvent{}, err
		}
		val.config = metadata
		return val, nil
	}
	return val, nil
}
