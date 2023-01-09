package kubecontrol

import (
	"bytes"
	"context"
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

func (svc *monitorService) getPodLogs(ctx context.Context, pod k8scorev1.Pod) ([]string, error) {
	maxTailLines := int64(10)
	podLogOpts := k8scorev1.PodLogOptions{TailLines: &maxTailLines}
	config, err := rest.InClusterConfig()
	if err != nil {
		svc.logger.Error("error on get cluster config", zap.Error(err))
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		svc.logger.Error("error on get client", zap.Error(err))
		return nil, err
	}
	svc.logger.Info("DEBUG log pod", zap.Any("pod", pod), zap.Any("clientSet", clientSet))
	req := clientSet.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		svc.logger.Error("error on get logs", zap.Error(err))
		return nil, err
	}
	defer func(podLogs io.ReadCloser) {
		err := podLogs.Close()
		if err != nil {
			svc.logger.Error("error closing log stream", zap.Error(err))
		}
	}(podLogs)

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		svc.logger.Error("error on copying buffer", zap.Error(err))
		return nil, err
	}
	str := buf.String()
	splitLogs := strings.Split(str, "\n")
	svc.logger.Info("DEBUG podLogs", zap.Strings("logs", splitLogs))
	return splitLogs, nil
}

func (svc *monitorService) getRunningPods(ctx context.Context) ([]k8scorev1.Pod, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		svc.logger.Error("error on get cluster config", zap.Error(err))
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		svc.logger.Error("error on get client", zap.Error(err))
		return nil, err
	}
	pods, err := clientSet.CoreV1().Pods(namespace).List(ctx, k8smetav1.ListOptions{})
	return pods.Items, err
}

func (svc *monitorService) monitorSinks(ctx context.Context) {
	runningCollectors, err := svc.getRunningPods(ctx)
	if err != nil {
		svc.logger.Error("error getting running pods on namespace", zap.Error(err))
		return
	}
	if len(runningCollectors) == 0 {
		svc.logger.Info("skipping, no running collectors")
		return
	}
	sinksRes, err := svc.sinksClient.RetrieveSinks(ctx, &sinkspb.SinksFilterReq{OtelEnabled: "enabled"})
	if err != nil {
		svc.logger.Error("error collecting sinks", zap.Error(err))
		return
	}
	svc.logger.Info("reading logs from collectors", zap.Int("collectors_length", len(sinksRes.Sinks)))
	for _, sink := range sinksRes.Sinks {
		var sinkCollector *k8scorev1.Pod
		for _, collector := range runningCollectors {
			if strings.Contains(collector.Name, sink.Id) {
				sinkCollector = &collector
				break
			}
		}
		if sinkCollector == nil {
			svc.logger.Warn("collector not found for sink, skipping", zap.String("sinkID", sink.Id))
			continue
		}
		logs, err := svc.getPodLogs(ctx, *sinkCollector)
		if err != nil {
			svc.logger.Error("error on getting logs, skipping", zap.Error(err))
			continue
		}
		status, err := analyzeLogs(logs)
		if sink.State != status {
			if err != nil {
				svc.logger.Info("updating status", zap.Any("before", sink.GetState()), zap.String("new status", status), zap.String("error_message (opt)", err.Error()))
			} else {
				svc.logger.Info("updating status", zap.Any("before", sink.GetState()), zap.String("new status", status))
			}
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
