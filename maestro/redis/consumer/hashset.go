package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/orb-community/orb/maestro/deployment"
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

type DeploymentHashsetRepository interface {
	GetDeploymentEntryFromSinkId(ctx context.Context, ownerId string, sinkId string) (string, error)
	CreateDeploymentEntry(ctx context.Context, deployment *deployment.Deployment) error
	UpdateDeploymentEntry(ctx context.Context, data config.SinkData) (err error)
	DeleteDeploymentEntry(ctx context.Context, sinkId string) error
}

type hashsetRepository struct {
	logger             *zap.Logger
	hashsetRedisClient *redis2.Client
}

func (es eventStore) GetDeploymentEntryFromSinkId(ctx context.Context, ownerId string, sinkId string) (string, error) {
	cmd := es.sinkerKeyRedisClient.HGet(ctx, deploymentKey, sinkId)
	if err := cmd.Err(); err != nil {
		es.logger.Error("error during redis reading of SinkId", zap.String("sink-id", sinkId), zap.Error(err))
		return "", err
	}
	return cmd.String(), nil
}

func (es eventStore) CreateDeploymentEntry(ctx context.Context, d *deployment.Deployment) error {
	deploy, err := config.BuildDeploymentJson(es.kafkaUrl, d)
	if err != nil {
		es.logger.Error("error trying to get deployment json for sink ID", zap.String("sinkId", d.SinkID), zap.Error(err))
		return err
	}

	// Instead create the deployment entry in postgres
	es.sinkerKeyRedisClient.HSet(ctx, deploymentKey, d.SinkID, deploy)

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
		Backend:   read(event, "backend", ""),
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
