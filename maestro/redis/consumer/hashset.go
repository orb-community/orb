package consumer

import (
	"context"
	"github.com/ns1labs/orb/maestro"
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
	es.logger.Info("Received maestro DELETE event from sinks ID=" + event.sinkID + ", Owner ID=" + event.ownerID)
	deployment, err := es.GetDeploymentEntryFromSinkId(ctx, event.sinkID)
	if err != nil {
		return err
	}
	err = es.maestroService.DeleteOtelCollector(ctx, event.sinkID, deployment)
	if err != nil {
		return err
	}
	es.client.HDel(ctx, deploymentKey, event.sinkID)
	return nil
}

// handleSinksCreateCollector will create Deployment Entry in Redis
func (es eventStore) handleSinksCreateCollector(ctx context.Context, event sinksUpdateEvent) error {
	es.logger.Info("Received event to Create DeploymentEntry from sinks ID=" + event.sinkID + ", Owner ID=" + event.ownerID)
	sinkUrl := event.config["sink_url"].(string)
	sinkUsername := event.config["username"].(string)
	sinkPassword := event.config["password"].(string)
	deploy, err := maestro.GetDeploymentJson(event.sinkID, sinkUrl, sinkUsername, sinkPassword)
	if err != nil {
		es.logger.Error("error trying to get deployment json for sink ID", zap.String("sinkId", event.sinkID))
		return err
	}
	es.client.HSet(ctx, deploymentKey, event.sinkID, deploy)

	return nil
}

// handleSinksUpdateCollector will update Deployment Entry in Redis and force update otel collector
func (es eventStore) handleSinksUpdateCollector(ctx context.Context, event sinksUpdateEvent) error {
	es.logger.Info("Received event to Update DeploymentEntry from sinks ID=" + event.sinkID + ", Owner ID=" + event.ownerID)
	sinkUrl := event.config["sink_url"].(string)
	sinkUsername := event.config["username"].(string)
	sinkPassword := event.config["password"].(string)
	deploy, err := maestro.GetDeploymentJson(event.sinkID, sinkUrl, sinkUsername, sinkPassword)
	if err != nil {
		es.logger.Error("error trying to get deployment json for sink ID", zap.String("sinkId", event.sinkID))
		return err
	}
	es.client.HSet(ctx, deploymentKey, event.sinkID, deploy)
	err = es.maestroService.UpdateOtelCollector(ctx, event.sinkID, deploy)
	if err != nil {
		return err
	}

	return nil
}

func decodeSinksUpdate(event map[string]interface{}) sinksUpdateEvent {
	val := sinksUpdateEvent{
		ownerID:   read(event, "owner", ""),
		sinkID:    read(event, "sink_id", ""),
		config:    readMetadata(event, "config"),
		timestamp: time.Time{},
	}
	return val
}
