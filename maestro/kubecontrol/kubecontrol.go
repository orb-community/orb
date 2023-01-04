package kubecontrol

import (
	"bufio"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"os/exec"
	"strings"
	"time"

	"go.uber.org/zap"
)

const namespace = "otelcollectors"

var _ Service = (*deployService)(nil)

type deployService struct {
	logger      *zap.Logger
	redisClient *redis.Client
}

func NewService(logger *zap.Logger, redisClient *redis.Client) Service {
	return &deployService{logger: logger, redisClient: redisClient}
}

type Service interface {
	// CreateOtelCollector - create an existing collector by id
	CreateOtelCollector(ctx context.Context, sinkID, deploymentEntry string) error

	// DeleteOtelCollector - delete an existing collector by id
	DeleteOtelCollector(ctx context.Context, sinkID, deploymentEntry string) error

	// UpdateOtelCollector - update an existing collector by id
	UpdateOtelCollector(ctx context.Context, sinkID, deploymentEntry string) error

	// CollectLogs - collect logs from the collector by sink-id
	CollectLogs(ctx context.Context, sinkID string) ([]string, error)
}

func (svc *deployService) CollectLogs(ctx context.Context, sinkId string) ([]string, error) {
	cmd := exec.Command("kubectl", "get logs", fmt.Sprintf("otel-collector-%s", sinkId), "-n", namespace)
	exporterLogs := make([]string, 10)
	watchLogsFunction := func(out *bufio.Scanner, err *bufio.Scanner) {
		if err.Scan() || out.Err() != nil {
			svc.logger.Error("failed to get logs for collector on sink")
			return
		}
		for out.Scan() && len(exporterLogs) < 10 {
			logEntry := out.Text()
			svc.logger.Info("debugging logEntry", zap.String("sinkId", sinkId), zap.String("logEntry", logEntry))
			exporterLogs = append(exporterLogs, logEntry)
		}
	}
	_, _, err := execCmd(ctx, cmd, svc.logger, watchLogsFunction)
	if err != nil {
		svc.logger.Error("Error reading the logs")
		exporterLogs = nil
	}
	return exporterLogs, err
}

func (svc *deployService) collectorDeploy(ctx context.Context, operation, sinkId, manifest string) error {

	fileContent := []byte(manifest)
	tmp := strings.Split(string(fileContent), "\n")
	newContent := strings.Join(tmp[1:], "\n")
	status, err := svc.getDeploymentState(ctx, sinkId)
	if err != nil {
		return err
	}
	if operation == "apply" {
		if status != "deleted" {
			svc.logger.Info("Already applied Sink ID", zap.String("sinkID", sinkId), zap.String("status", status))
			return nil
		}
	} else if operation == "delete" {
		if status == "deleted" {
			svc.logger.Info("Already deleted Sink ID", zap.String("sinkID", sinkId), zap.String("status", status))
			return nil
		}
	}

	err = os.WriteFile("/tmp/otel-collector-"+sinkId+".json", []byte(newContent), 0644)
	if err != nil {
		svc.logger.Error("failed to write file content", zap.Error(err))
		return err
	}

	stdOutListenFunction := func(out *bufio.Scanner, err *bufio.Scanner) {
		for out.Scan() {
			fmt.Println(out.Text())
			svc.logger.Info("Deploy Info: " + out.Text())
		}
		for err.Scan() {
			fmt.Println(err.Text())
			svc.logger.Info("Deploy Error: " + err.Text())
		}
	}

	// execute action
	cmd := exec.Command("kubectl", operation, "-f", "/tmp/otel-collector-"+sinkId+".json", "-n", namespace)
	_, _, err = execCmd(ctx, cmd, svc.logger, stdOutListenFunction)

	if err == nil {
		if operation == "apply" {
			err := svc.setNewDeploymentState(ctx, sinkId, "idle")
			if err != nil {
				return err
			}
		} else if operation == "delete" {
			err := svc.setNewDeploymentState(ctx, sinkId, "idle")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func execCmd(_ context.Context, cmd *exec.Cmd, logger *zap.Logger, stdOutFunc func(stdOut *bufio.Scanner, stdErr *bufio.Scanner)) (*bufio.Scanner, *bufio.Scanner, error) {
	stdoutReader, _ := cmd.StdoutPipe()
	stdoutScanner := bufio.NewScanner(stdoutReader)
	stderrReader, _ := cmd.StderrPipe()
	stderrScanner := bufio.NewScanner(stderrReader)
	go stdOutFunc(stdoutScanner, stderrScanner)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error : %v \n", err)
		logger.Error("Collector Deploy Error", zap.Error(err))
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		logger.Error("Collector Deploy Error", zap.Error(err))
	}
	return stdoutScanner, stderrScanner, err
}

func (svc *deployService) getDeploymentState(ctx context.Context, sinkId string) (string, error) {
	key := CollectorStatusKey + "." + sinkId
	args := redis.ZRangeArgs{
		Key:     key,
		Start:   nil,
		Stop:    nil,
		ByScore: false,
		ByLex:   false,
		Rev:     false,
		Offset:  0,
		Count:   0,
	}
	cmd := svc.redisClient.ZRangeArgsWithScores(ctx, args)
	slice, err := cmd.Result()
	if err != nil {
		return "", cmd.Err()
	}
	svc.logger.Info("debug returned slice", zap.Any("slice", slice))
	value := slice[0]
	entry := value.Member.(CollectorStatusSortedSetEntry)
	return entry.Status, nil
}

func (svc *deployService) setNewDeploymentState(ctx context.Context, sinkId, state string) error {
	key := CollectorStatusKey + "." + sinkId
	entry := redis.Z{
		Score: float64(time.Now().Unix()),
		Member: CollectorStatusSortedSetEntry{
			SinkId:       sinkId,
			Status:       "idle",
			ErrorMessage: nil,
		},
	}
	intCmd := svc.redisClient.ZAdd(ctx, key, &entry)
	if intCmd.Err() != nil {
		return intCmd.Err()
	}
	return nil
}

func (svc *deployService) CreateOtelCollector(ctx context.Context, sinkID, deploymentEntry string) error {
	err := svc.collectorDeploy(ctx, "apply", sinkID, deploymentEntry)
	if err != nil {
		return err
	}

	return nil
}

func (svc *deployService) UpdateOtelCollector(ctx context.Context, sinkID, deploymentEntry string) error {
	err := svc.DeleteOtelCollector(ctx, sinkID, deploymentEntry)
	if err != nil {
		return err
	}
	err = svc.CreateOtelCollector(ctx, sinkID, deploymentEntry)
	if err != nil {
		return err
	}
	return nil
}

func (svc *deployService) DeleteOtelCollector(ctx context.Context, sinkID, deploymentEntry string) error {
	err := svc.collectorDeploy(ctx, "delete", sinkID, deploymentEntry)
	if err != nil {
		return err
	}
	return nil
}
