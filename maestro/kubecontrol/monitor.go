package kubecontrol

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

const MonitorFixedDuration = 1 * time.Minute
const TimeDiffForFetchingLogs = 5 * time.Minute

func NewMonitorService(logger *zap.Logger, redis *redis.Client, kubecontrol *Service) MonitorService {
	return &monitorService{
		logger:      logger,
		redisClient: redis,
		kubecontrol: *kubecontrol,
	}
}

type MonitorService interface {
	Start(ctx context.Context, cancelFunc context.CancelFunc) error
}

type monitorService struct {
	logger      *zap.Logger
	redisClient *redis.Client
	kubecontrol Service
}

func (svc *monitorService) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	go func() {
		ticker := time.NewTicker(MonitorFixedDuration)
		svc.logger.Info("start monitor routine", zap.Any("routine", ctx))
		defer func() {
			cancelFunc()
			svc.logger.Info("stopping monitor routine")
		}()
		for {
			select {
			case <-ctx.Done():
				cancelFunc()
				return
			case _ = <-ticker.C:
				svc.logger.Info("monitoring sinks")
				svc.monitorSinks(ctx)
			}
		}
	}()
	return nil
}

func (svc *monitorService) monitorSinks(ctx context.Context) {
	queryCmd := svc.redisClient.ZRange(ctx, CollectorStatusKey, time.Now().Add(-TimeDiffForFetchingLogs).Unix(), 0)
	if queryCmd.Err() != nil {
		svc.logger.Error("error collecting collectors keys", zap.Error(queryCmd.Err()))
		return
	}
	collectorSlice, err := queryCmd.Result()
	if err != nil {
		svc.logger.Error("error collecting collectors keys", zap.Error(queryCmd.Err()))
		return
	}
	svc.logger.Info("reading logs from collectors", zap.Int("collectors_length", len(collectorSlice)))
	for _, collectorJson := range collectorSlice {
		var collector CollectorStatusSortedSetEntry
		err = json.Unmarshal([]byte(collectorJson), &collector)
		logs, err := svc.kubecontrol.CollectLogs(ctx, collector.SinkId)
		if err != nil {
			return
		}
		status, err := analyzeLogs(logs)
		if status != collector.Status {
			svc.logger.Info("updating status", zap.Any("before", collector), zap.String("new status", status), zap.String("error_message (opt)", err.Error()))
			collector.Status = status
			if status == "error" {
				errorStr := err.Error()
				collector.ErrorMessage = &errorStr
			}
		}
		return
	}

}

// WIP
func analyzeLogs(logEntry []string) (string, error) {

	return "active", nil
}
