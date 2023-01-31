package kubecontrol

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	maestroconfig "github.com/ns1labs/orb/maestro/config"
	maestroredis "github.com/ns1labs/orb/maestro/redis"
	sinkspb "github.com/ns1labs/orb/sinks/pb"
	"go.uber.org/zap"
	k8scorev1 "k8s.io/api/core/v1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	activityPrefix       = "sinker_activity"
	MonitorFixedDuration = 1 * time.Minute
	TimeDiffActiveIdle   = 5 * time.Minute
)

func NewMonitorService(logger *zap.Logger, sinksClient *sinkspb.SinkServiceClient, redisClient *redis.Client, redisCache *redis.Client, kubecontrol *Service) MonitorService {
	deploymentChecks := make(map[string]int)
	return &monitorService{
		logger:           logger,
		sinksClient:      *sinksClient,
		redisClient:      redisClient,
		redisCache:       redisCache,
		kubecontrol:      *kubecontrol,
		deploymentChecks: deploymentChecks,
	}
}

type MonitorService interface {
	Start(ctx context.Context, cancelFunc context.CancelFunc) error
	GetRunningPods(ctx context.Context) ([]string, error)
}

type monitorService struct {
	logger           *zap.Logger
	sinksClient      sinkspb.SinkServiceClient
	redisClient      *redis.Client
	redisCache       *redis.Client
	kubecontrol      Service
	deploymentChecks map[string]int //to check deployment error
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
	svc.logger.Info("logs length", zap.Int("amount line logs", len(splitLogs)))
	return splitLogs, nil
}

func (svc *monitorService) GetRunningPods(ctx context.Context) ([]string, error) {
	pods, err := svc.getRunningPods(ctx)
	if err != nil {
		svc.logger.Error("error getting running collectors")
		return nil, err
	}
	runningSinks := make([]string, len(pods))
	if len(pods) > 0 {
		for i, pod := range pods {
			runningSinks[i] = strings.TrimPrefix(pod.Name, "otel-")
		}
		return runningSinks, nil
	}
	return nil, nil
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
			svc.logger.Warn("collector not found for sink, checking to set state as error", zap.String("sinkID", sink.Id))
			// if collector dont spin up in 30 minutes should report error on collector deployment
			svc.deploymentChecks[sink.Id]++
			if svc.deploymentChecks[sink.Id] >= 30 {
				err := errors.New("permanent error: opentelemetry collector deployment error")
				svc.publishSinkStateChange(sink, "error", err, err)
				svc.deploymentChecks[sink.Id] = 0
			}
			continue
		}
		var data maestroconfig.SinkData
		if err := json.Unmarshal(sink.Config, &data); err != nil {
			svc.logger.Warn("failed to unmarshal sink config, skipping", zap.String("sink-id", sink.Id))
			continue
		}
		data.SinkID = sink.Id
		data.OwnerID = sink.OwnerID
		data.LastRemoteWrite = time.Now()
		logs, err := svc.getPodLogs(ctx, *sinkCollector)
		if err != nil {
			svc.logger.Error("error on getting logs, skipping", zap.Error(err))
			continue
		}
		status, logsErr := svc.analyzeLogs(logs)
		if status == "fail" {
			svc.logger.Error("error during analyze logs", zap.Error(logsErr))
			continue
		}
		// here we should check if LastRemoteWrite is up-to-date, otherwise we need to set sink as idle
		lastActivity, activityErr := svc.GetActivity(sink.Id)
		if activityErr != nil {
			svc.logger.Error("error on getting last collector activity, skipping", zap.Error(activityErr))
			continue
		}
		idleLimit := lastActivity + 1800 // 30 minutes
		// we should change sink state just when it is 'active'.
		// this state can be 'error', if any error was found on otel collector during analyzeLogs() or
		// we can set it as idle when we see that lastRemoteWrite is older than 30 minutes, however this
		// monitoring state function already exists on sinker
		if sink.GetState() == "active" {
			// check if idle
			if time.Now().Unix() >= idleLimit {
				svc.publishSinkStateChange(sink, "idle", logsErr, err)
				// TODOs:

				// - set state on sinkerCache to idle too or remove from sinker redis
				// - remove collector from k8s
				// - remove sinkid from redis sinker_activity
				//svc.kubecontrol.DeleteOtelCollector(sink.OwnerID, sink.Id, ?)

			} else if sink.GetState() != status { //updating status
				if err != nil {
					svc.logger.Info("updating status", zap.Any("before", sink.GetState()), zap.String("new status", status), zap.String("error_message (opt)", err.Error()), zap.String("SinkID", sink.Id), zap.String("ownerID", sink.OwnerID))
				} else {
					svc.logger.Info("updating status", zap.Any("before", sink.GetState()), zap.String("new status", status), zap.String("SinkID", sink.Id), zap.String("ownerID", sink.OwnerID))
				}
				svc.publishSinkStateChange(sink, status, logsErr, err)

			}
		}
	}
}

// check last collector activity
func (svc *monitorService) GetActivity(sinkID string) (int64, error) {
	if sinkID == "" {
		return 0, errors.New("invalid parameters")
	}
	skey := fmt.Sprintf("%s:%s", activityPrefix, sinkID)
	secs, err := svc.redisCache.Get(context.Background(), skey).Result()
	if err != nil {
		return 0, err
	}
	lastActivity, _ := strconv.ParseInt(secs, 10, 64)
	return lastActivity, nil
}

func (svc *monitorService) publishSinkStateChange(sink *sinkspb.SinkRes, status string, logsErr error, err error) {
	streamID := "orb.sinker"
	logMessage := ""
	if logsErr != nil {
		logMessage = logsErr.Error()
	}
	event := maestroredis.SinkerUpdateEvent{
		SinkID:    sink.Id,
		Owner:     sink.OwnerID,
		State:     status,
		Msg:       logMessage,
		Timestamp: time.Now(),
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

// analyzeLogs, will check for errors in exporter, and will return as follows
//
//		for errors 429 will send a "warning" state, plus message of too many requests
//		for any other errors, will add error and message
//		if no error message on exporter, will log as active
//	 logs from otel-collector are coming in the standard from https://pkg.go.dev/log,
func (svc *monitorService) analyzeLogs(logEntry []string) (status string, err error) {
	for _, logLine := range logEntry {
		if len(logLine) > 24 {
			// known errors
			if strings.Contains(logLine, "Permanent error: remote write returned HTTP status 401 Unauthorized") {
				errorMessage := "Permanent error: remote write returned HTTP status 401 Unauthorized"
				return "error", errors.New(errorMessage)
			}
			if strings.Contains(logLine, "Permanent error: remote write returned HTTP status 429 Too Many Requests") {
				errorMessage := "Permanent error: remote write returned HTTP status 429 Too Many Requests"
				return "warning", errors.New(errorMessage)
			}
			// other errors
			if strings.Contains(logLine, "error") {
				errStringLog := strings.TrimRight(logLine, "error")
				if len(errStringLog) > 4 {
					aux := strings.Split(errStringLog, "\t")
					numItems := len(aux)
					if numItems > 3 {
						jsonError := aux[4]
						errorJson := make(map[string]interface{})
						err := json.Unmarshal([]byte(jsonError), &errorJson)
						if err != nil {
							return "fail", err
						}
						if errorJson != nil && errorJson["error"] != nil {
							errorMessage := errorJson["error"].(string)
							return "error", errors.New(errorMessage)
						}
					} else {
						return "error", errors.New("sink configuration error: please review your sink parameters")
					}
				}
			}
		}
	}
	// if nothing happens is active, cuz idle state sinker mciro service already is handling
	return "active", nil
}
