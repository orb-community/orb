package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	"strconv"
	"time"

	"github.com/ns1labs/orb/maestro/config"
	"github.com/ns1labs/orb/maestro/redis"
	"github.com/ns1labs/orb/pkg/types"
	sinkspb "github.com/ns1labs/orb/sinks/pb"
	"go.uber.org/zap"
)

const (
	deploymentKey  = "orb.sinks.deployment"
	activityPrefix = "sinker_activity"
)

func (es eventStore) GetDeploymentEntryFromSinkId(ctx context.Context, sinkId string) (string, error) {
	cmd := es.streamRedisClient.HGet(ctx, deploymentKey, sinkId)
	if err := cmd.Err(); err != nil {
		es.logger.Error("error during redis reading of SinkId", zap.String("sink-id", sinkId), zap.Error(err))
		return "", err
	}
	return cmd.String(), nil
}

// handleSinksDeleteCollector will delete Deployment Entry and force delete otel collector
func (es eventStore) handleSinksDeleteCollector(ctx context.Context, event redis.SinksUpdateEvent) error {
	es.logger.Info("Received maestro DELETE event from sinks ID", zap.String("sinkID", event.SinkID), zap.String("owner", event.Owner))
	deployment, err := es.GetDeploymentEntryFromSinkId(ctx, event.SinkID)
	if err != nil {
		es.logger.Error("did not find collector entry for sink", zap.String("sink-id", event.SinkID))
		return err
	}
	err = es.kubecontrol.DeleteOtelCollector(ctx, event.Owner, event.SinkID, deployment)
	if err != nil {
		return err
	}
	es.streamRedisClient.HDel(ctx, deploymentKey, event.SinkID)
	return nil
}

// handleSinksCreateCollector will create Deployment Entry in Redis
func (es eventStore) handleSinksCreateCollector(ctx context.Context, event redis.SinksUpdateEvent) error {
	es.logger.Info("Received event to Create DeploymentEntry from sinks ID", zap.String("sinkID", event.SinkID), zap.String("owner", event.Owner))
	sinkData, err := es.sinksClient.RetrieveSink(ctx, &sinkspb.SinkByIDReq{
		SinkID:  event.SinkID,
		OwnerID: event.Owner,
	})
	if err != nil {
		es.logger.Error("could not fetch info for sink", zap.String("sink-id", event.SinkID), zap.Error(err))
	}
	var data config.SinkData
	if err := json.Unmarshal(sinkData.Config, &data); err != nil {
		return err
	}
	sinkUrl := data.Url
	sinkUsername := data.User
	sinkPassword := data.Password
	err2 := es.CreateDeploymentEntry(ctx, event.SinkID, sinkUrl, sinkUsername, sinkPassword)
	if err2 != nil {
		return err2
	}

	return nil
}

func (es eventStore) CreateDeploymentEntry(ctx context.Context, sinkId, sinkUrl, sinkUsername, sinkPassword string) error {
	deploy, err := config.GetDeploymentJson(es.kafkaUrl, sinkId, sinkUrl, sinkUsername, sinkPassword)
	if err != nil {
		es.logger.Error("error trying to get deployment json for sink ID", zap.String("sinkId", sinkId))
		return err
	}

	es.streamRedisClient.HSet(ctx, deploymentKey, sinkId, deploy)
	return nil
}

// handleSinksUpdateCollector will update Deployment Entry in Redis and force update otel collector
func (es eventStore) handleSinksUpdateCollector(ctx context.Context, event redis.SinksUpdateEvent) error {
	es.logger.Info("Received event to Update DeploymentEntry from sinks ID", zap.String("sinkID", event.SinkID), zap.String("owner", event.Owner))
	sinkData, err := es.sinksClient.RetrieveSink(ctx, &sinkspb.SinkByIDReq{
		SinkID:  event.SinkID,
		OwnerID: event.Owner,
	})
	if err != nil {
		es.logger.Error("could not fetch info for sink", zap.String("sink-id", event.SinkID), zap.Error(err))
	}
	var data config.SinkData
	if err := json.Unmarshal(sinkData.Config, &data); err != nil {
		return err
	}
	sinkUrl := data.Url
	sinkUsername := data.User
	sinkPassword := data.Password
	deploy, err := config.GetDeploymentJson(es.kafkaUrl, event.SinkID, sinkUrl, sinkUsername, sinkPassword)
	if err != nil {
		es.logger.Error("error trying to get deployment json for sink ID", zap.String("sinkId", event.SinkID))
		return err
	}
	es.streamRedisClient.HSet(ctx, deploymentKey, event.SinkID, deploy)
	err = es.kubecontrol.UpdateOtelCollector(ctx, event.Owner, event.SinkID, deploy)
	if err != nil {
		return err
	}

	return nil
}

func (es eventStore) UpdateSinkCache(ctx context.Context, data config.SinkData) (err error) {
	data.State = config.Unknown
	keyPrefix := "sinker_key"
	skey := fmt.Sprintf("%s-%s:%s", keyPrefix, data.OwnerID, data.SinkID)
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err = es.sinkerKeyRedisClient.Set(ctx, skey, bytes, 0).Err(); err != nil {
		return err
	}
	return
}

// GetActivity collector activity
func (es eventStore) GetActivity(sinkID string) (int64, error) {
	if sinkID == "" {
		return 0, errors.New("invalid parameters")
	}
	skey := fmt.Sprintf("%s:%s", activityPrefix, sinkID)
	secs, err := es.sinkerKeyRedisClient.Get(context.Background(), skey).Result()
	if err != nil {
		return 0, err
	}
	lastActivity, _ := strconv.ParseInt(secs, 10, 64)
	return lastActivity, nil
}

func (es eventStore) RemoveSinkActivity(ctx context.Context, sinkId string) error {
	skey := fmt.Sprintf("%s:%s", activityPrefix, sinkId)
	cmd := es.sinkerKeyRedisClient.Del(ctx, skey, sinkId)
	if err := cmd.Err(); err != nil {
		es.logger.Error("error during redis reading of SinkId", zap.String("sink-id", sinkId), zap.Error(err))
		return err
	}
	return nil
}

func (es eventStore) PublishSinkStateChange(sink *sinkspb.SinkRes, status string, logsErr error, err error) {
	streamID := "orb.sinker"
	logMessage := ""
	if logsErr != nil {
		logMessage = logsErr.Error()
	}
	event := redis.SinkerUpdateEvent{
		SinkID: sink.Id,
		Owner:  sink.OwnerID,
		State:  status,

		Msg:       logMessage,
		Timestamp: time.Now(),
	}

	record := &redis2.XAddArgs{
		Stream: streamID,
		Values: event.Encode(),
	}
	err = es.sinkerKeyRedisClient.XAdd(context.Background(), record).Err()
	if err != nil {
		es.logger.Error("error sending event to event store", zap.Error(err))
	}
}

func decodeSinksEvent(event map[string]interface{}, operation string) (redis.SinksUpdateEvent, error) {
	val := redis.SinksUpdateEvent{
		SinkID:    read(event, "sink_id", ""),
		Owner:     read(event, "owner", ""),
		Timestamp: time.Now(),
	}
	if operation != sinksDelete {
		var metadata types.Metadata
		if err := json.Unmarshal([]byte(read(event, "config", "")), &metadata); err != nil {
			return redis.SinksUpdateEvent{}, err
		}
		val.Config = metadata
		return val, nil
	}

	return val, nil
}
