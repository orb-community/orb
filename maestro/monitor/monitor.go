package monitor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/orb-community/orb/maestro/deployment"
	"github.com/orb-community/orb/maestro/redis/producer"
	"io"
	"strings"
	"time"

	maestroconfig "github.com/orb-community/orb/maestro/config"
	"github.com/orb-community/orb/maestro/kubecontrol"
	sinkspb "github.com/orb-community/orb/sinks/pb"
	"go.uber.org/zap"
	k8scorev1 "k8s.io/api/core/v1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	idleTimeSeconds = 600
	TickerForScan   = 1 * time.Minute
	namespace       = "otelcollectors"
)

func NewMonitorService(logger *zap.Logger, sinksClient *sinkspb.SinkServiceClient, mp producer.Producer, kubecontrol *kubecontrol.Service) Service {
	return &monitorService{
		logger:          logger,
		sinksClient:     *sinksClient,
		maestroProducer: mp,
		kubecontrol:     *kubecontrol,
	}
}

type Service interface {
	Start(ctx context.Context, cancelFunc context.CancelFunc) error
	GetRunningPods(ctx context.Context) ([]string, error)
}

type monitorService struct {
	logger          *zap.Logger
	sinksClient     sinkspb.SinkServiceClient
	maestroProducer producer.Producer
	deploymentSvc   deployment.Service
	kubecontrol     kubecontrol.Service
}

func (svc *monitorService) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	go func(ctx context.Context, cancelFunc context.CancelFunc) {
		ticker := time.NewTicker(TickerForScan)
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
	sinceSeconds := int64(300)
	podLogOpts := k8scorev1.PodLogOptions{TailLines: &maxTailLines, SinceSeconds: &sinceSeconds}
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
	for _, collector := range runningCollectors {
		var sink *sinkspb.SinkRes
		for _, sinkRes := range sinksRes.Sinks {
			if strings.Contains(collector.Name, sinkRes.Id) {
				sink = sinkRes
				break
			}
		}
		if sink == nil {
			svc.logger.Warn("collector not found for sink, depleting collector", zap.String("collector name", collector.Name))
			sinkId := collector.Name[5:41]
			deploymentName := "otel-" + sinkId
			err = svc.kubecontrol.KillOtelCollector(ctx, deploymentName, sinkId)
			if err != nil {
				svc.logger.Error("error removing otel collector", zap.Error(err))
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
		var logsErr error
		var status string
		logs, err := svc.getPodLogs(ctx, collector)
		if err != nil {
			svc.logger.Error("error on getting logs, skipping", zap.Error(err))
			continue
		}
		var logErrMsg string
		status, logsErr = svc.analyzeLogs(logs)
		if status == "fail" {
			svc.logger.Error("error during analyze logs", zap.Error(logsErr))
			continue
		}
		if logsErr != nil {
			logErrMsg = logsErr.Error()
		}

		//set the new sink status if changed during checks
		if sink.GetState() != status && status != "" {
			svc.logger.Info("changing sink status", zap.Any("before", sink.GetState()), zap.String("new status", status), zap.String("SinkID", sink.Id), zap.String("ownerID", sink.OwnerID))
			if err != nil {
				svc.logger.Error("error updating status", zap.Any("before", sink.GetState()), zap.String("new status", status), zap.String("error_message (opt)", err.Error()), zap.String("SinkID", sink.Id), zap.String("ownerID", sink.OwnerID))
			} else {
				svc.logger.Info("updating status", zap.Any("before", sink.GetState()), zap.String("new status", status), zap.String("SinkID", sink.Id), zap.String("ownerID", sink.OwnerID))
				err = svc.deploymentSvc.UpdateStatus(ctx, sink.OwnerID, sink.Id, status, logErrMsg)
			}
		}
	}
}

// analyzeLogs, will check for errors in exporter, and will return as follows
// for errors 429 will send a "warning" state, plus message of too many requests
// for any other errors, will add error and message
// if no error message on exporter, will log as active
// logs from otel-collector are coming in the standard from https://pkg.go.dev/log,
func (svc *monitorService) analyzeLogs(logEntry []string) (status string, err error) {
	for _, logLine := range logEntry {
		if len(logLine) > 24 {
			// known errors
			if strings.Contains(logLine, "401 Unauthorized") {
				errorMessage := "error: remote write returned HTTP status 401 Unauthorized"
				return "error", errors.New(errorMessage)
			}
			if strings.Contains(logLine, "404 Not Found") {
				errorMessage := "error: remote write returned HTTP status 404 Not Found"
				return "error", errors.New(errorMessage)
			}
			if strings.Contains(logLine, "502 Bad Gateway") {
				errorMessage := "error: remote write returned HTTP status 502 Bad Gateway"
				return "error", errors.New(errorMessage)
			}
			if strings.Contains(logLine, "504 Gateway Timeout") {
				errorMessage := "error: remote write returned HTTP status 504 Gateway Timeout"
				return "error", errors.New(errorMessage)
			}
			// known warnings
			if strings.Contains(logLine, "429 Too Many Requests") {
				errorMessage := "error: remote write returned HTTP status 429 Too Many Requests"
				return "warning", errors.New(errorMessage)
			}
			if strings.Contains(logLine, "400 Bad Request") {
				errorMessage := "error: remote write returned HTTP status 400 Bad Request"
				return "warning", errors.New(errorMessage)
			}
			// other generic errors
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
				} else {
					return "error", errors.New("sink configuration error: please review your sink parameters")
				}
			}
		}
	}
	// if nothing happens on logs is active
	return "active", nil
}
