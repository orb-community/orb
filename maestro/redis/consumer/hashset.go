package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	redis2 "github.com/go-redis/redis/v8"

	"github.com/orb-community/orb/maestro/config"
	"github.com/orb-community/orb/maestro/redis"
	"github.com/orb-community/orb/pkg/types"
	sinkspb "github.com/orb-community/orb/sinks/pb"
	"go.uber.org/zap"
)

const (
	deploymentKey  = "orb.sinks.deployment"
	activityPrefix = "sinker_activity"
	streamLen      = 1000
)

func (es eventStore) GetDeploymentEntryFromSinkId(ctx context.Context, sinkId string) (string, error) {
	cmd := es.sinkerKeyRedisClient.HGet(ctx, deploymentKey, sinkId)
	if err := cmd.Err(); err != nil {
		es.logger.Error("error during redis reading of SinkId", zap.String("sink-id", sinkId), zap.Error(err))
		return "", err
	}
	return cmd.String(), nil
}

// handleSinksDeleteCollector will delete Deployment Entry and force delete otel collector
func (es eventStore) handleSinksDeleteCollector(ctx context.Context, event redis.SinksUpdateEvent) error {
	es.logger.Info("Received maestro DELETE event from sinks ID", zap.String("sinkID", event.SinkID), zap.String("owner", event.Owner))
	deploymentEntry, err := es.GetDeploymentEntryFromSinkId(ctx, event.SinkID)
	if err != nil {
		es.logger.Error("did not find collector entry for sink", zap.String("sink-id", event.SinkID))
		return err
	}
	err = es.kubecontrol.DeleteOtelCollector(ctx, event.Owner, event.SinkID, deploymentEntry)
	if err != nil {
		return err
	}
	err = es.sinkerKeyRedisClient.HDel(ctx, deploymentKey, event.SinkID).Err()
	if err != nil {
		return err
	}
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

	var metadata types.Metadata
	if err := json.Unmarshal(sinkData.Config, &metadata); err != nil {
		return err
	}
	data := config.SinkData{
		SinkID:  sinkData.Id,
		OwnerID: sinkData.OwnerID,
		Backend: sinkData.Backend,
		Config:  metadata,
	}
	err2 := es.CreateDeploymentEntry(ctx, data)
	if err2 != nil {
		return err2
	}

	return nil
}

func (es eventStore) CreateDeploymentEntry(ctx context.Context, sink config.SinkData) error {
	deploy, err := config.GetDeploymentJson(es.kafkaUrl, sink)
	if err != nil {
		es.logger.Error("error trying to get deployment json for sink ID", zap.String("sinkId", sink.SinkID), zap.Error(err))
		return err
	}

	es.sinkerKeyRedisClient.HSet(ctx, deploymentKey, sink.SinkID, deploy)
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
	var metadata types.Metadata
	if err := json.Unmarshal(sinkData.Config, &metadata); err != nil {
		return err
	}
	es.logger.Info("metadata", zap.Any("metadata", metadata))
	data := config.SinkData{
		SinkID:  sinkData.Id,
		OwnerID: sinkData.OwnerID,
		Backend: sinkData.Backend,
		Config:  metadata,
	}
	_ = data.State.SetFromString(sinkData.State)

	deploy, err := config.GetDeploymentJson(es.kafkaUrl, data)
	es.logger.Info("deploy", zap.Any("deploy", deploy))
	if err != nil {
		es.logger.Error("error trying to get deployment json for sink ID", zap.String("sinkId", event.SinkID), zap.Error(err))
		return err
	}
	err = es.sinkerKeyRedisClient.HSet(ctx, deploymentKey, event.SinkID, deploy).Err()
	if err != nil {
		es.logger.Error("error trying to update deployment json for sink ID", zap.String("sinkId", event.SinkID), zap.Error(err))
		return err
	}
	err = es.kubecontrol.UpdateOtelCollector(ctx, event.Owner, event.SinkID, deploy)
	if err != nil {
		return err
	}
	return nil
}

func (es eventStore) UpdateSinkCache(ctx context.Context, data config.SinkData) (err error) {
	keyPrefix := "sinker_key"
	skey := fmt.Sprintf("%s-%s:%s", keyPrefix, data.OwnerID, data.SinkID)
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err = es.sinkerKeyRedisClient.Set(ctx, skey, bytes, 0).Err(); err != nil {
		es.logger.Error("failed to update sink cache", zap.Error(err))
		return err
	}
	return
}

func (es eventStore) UpdateSinkStateCache(ctx context.Context, data config.SinkData) (err error) {
	keyPrefix := "sinker_key"
	skey := fmt.Sprintf("%s-%s:%s", keyPrefix, data.OwnerID, data.SinkID)
	bytes, err := json.Marshal(data)
	if err != nil {
		es.logger.Error("error update sink cache state", zap.Error(err))
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
		SinkID:    sink.Id,
		Owner:     sink.OwnerID,
		State:     status,
		Msg:       logMessage,
		Timestamp: time.Now(),
	}

	record := &redis2.XAddArgs{
		Stream: streamID,
		Values: event.Encode(),
		MaxLen: streamLen,
		Approx: true,
	}
	err = es.streamRedisClient.XAdd(context.Background(), record).Err()
	if err != nil {
		es.logger.Error("error sending event to event store", zap.Error(err))
	}
	es.logger.Info("Maestro notified change of status for sink", zap.String("newState", status), zap.String("sink-id", sink.Id))
}

func decodeSinksEvent(event map[string]interface{}, operation string) (redis.SinksUpdateEvent, error) {
	val := redis.SinksUpdateEvent{
		SinkID:    read(event, "sink_id", ""),
		Owner:     read(event, "owner", ""),
		Config:    readMetadata(event, "config"),
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
