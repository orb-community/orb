package kubecontrol

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	rediscons1 "github.com/ns1labs/orb/maestro/redis"
	sinkspb "github.com/ns1labs/orb/sinks/pb"
	"go.uber.org/zap"
	"io"
	k8scorev1 "k8s.io/api/core/v1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"
	"time"
)

const streamID = "orb.sinker"

const MonitorFixedDuration = 1 * time.Minute
const TimeDiffForFetchingLogs = 5 * time.Minute

func NewMonitorService(logger *zap.Logger, sinksClient *sinkspb.SinkServiceClient, redisClient *redis.Client, kubecontrol *Service) MonitorService {
	return &monitorService{
		logger:      logger,
		sinksClient: *sinksClient,
		redisClient: redisClient,
		kubecontrol: *kubecontrol,
	}
}

type MonitorService interface {
	Start(ctx context.Context, cancelFunc context.CancelFunc) error
}

type monitorService struct {
	logger      *zap.Logger
	sinksClient sinkspb.SinkServiceClient
	redisClient *redis.Client
	kubecontrol Service
}

func (svc *monitorService) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	go func(ctx context.Context, cancelFunc context.CancelFunc) {
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
	}(ctx, cancelFunc)
	return nil
}

func (svc *monitorService) getPodLogs(ctx context.Context, collectorName string) ([]string, error) {
	pod := k8scorev1.Pod{ObjectMeta: k8smetav1.ObjectMeta{Name: collectorName}}
	podLogOpts := k8scorev1.PodLogOptions{}
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return nil, err
	}
	defer func(podLogs io.ReadCloser) {
		err := podLogs.Close()
		if err != nil {

		}
	}(podLogs)

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return nil, err
	}
	str := buf.String()

	return strings.Split(str, "\n"), nil
}

func (svc *monitorService) monitorSinks(ctx context.Context) {

	sinksRes, err := svc.sinksClient.RetrieveSinks(ctx, &sinkspb.SinksFilterReq{OtelEnabled: "enabled"})
	if err != nil {
		svc.logger.Error("error collecting collectors keys", zap.Error(err))
		return
	}
	svc.logger.Info("reading logs from collectors", zap.Int("collectors_length", len(sinksRes.Sinks)))
	for _, sink := range sinksRes.Sinks {
		logs, err := svc.getPodLogs(ctx, fmt.Sprintf("otel-%s", sink.Id))
		if err != nil {
			return
		}
		status, err := analyzeLogs(logs)
		if status != sink.GetState() {
			svc.logger.Info("updating status", zap.Any("before", sink.GetState()), zap.String("new status", status), zap.String("error_message (opt)", err.Error()))
			event := rediscons1.SinkerUpdateEvent{
				SinkID:    sink.Id,
				Owner:     sink.OwnerID,
				State:     sink.State,
				Timestamp: time.Now(),
			}
			if status == "error" {
				event.Msg = err.Error()
			}
			record := &redis.XAddArgs{
				Stream: streamID,
				Values: event.Encode(),
			}
			err = svc.redisClient.XAdd(context.Background(), record).Err()
			if err != nil {
				svc.logger.Error("error sending event to event store", zap.Error(err))
			}
		}
		return
	}

}

// WIP
func analyzeLogs(logEntry []string) (string, error) {

	return "active", nil
}
