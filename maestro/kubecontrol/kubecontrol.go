package kubecontrol

import (
	"bufio"
	"context"
	"encoding/json"
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
	CreateOtelCollector(ctx context.Context, ownerID, sinkID, deploymentEntry string) error

	// DeleteOtelCollector - delete an existing collector by id
	DeleteOtelCollector(ctx context.Context, ownerID, sinkID, deploymentEntry string) error

	// UpdateOtelCollector - update an existing collector by id
	UpdateOtelCollector(ctx context.Context, ownerID, sinkID, deploymentEntry string) error
}

func (svc *deployService) collectorDeploy(ctx context.Context, operation, ownerID, sinkId, manifest string) error {

	fileContent := []byte(manifest)
	tmp := strings.Split(string(fileContent), "\n")
	newContent := strings.Join(tmp[1:], "\n")
	status, err := svc.getDeploymentState(ctx, ownerID, sinkId)
	if err != nil {
		return err
	}
	if operation == "apply" {
		if status != "deleted" {
			svc.logger.Info("Already applied Sink ID", zap.String("ownerID", ownerID), zap.String("sinkID", sinkId), zap.String("status", status))
			return nil
		}
	} else if operation == "delete" {
		if status == "deleted" {
			svc.logger.Info("Already deleted Sink ID", zap.String("ownerID", ownerID), zap.String("sinkID", sinkId), zap.String("status", status))
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
			svc.logger.Info("Deploy Info: " + out.Text())
		}
		for err.Scan() {
			svc.logger.Info("Deploy Error: " + err.Text())
		}
	}

	// execute action
	cmd := exec.Command("kubectl", operation, "-f", "/tmp/otel-collector-"+sinkId+".json", "-n", namespace)
	_, _, err = execCmd(ctx, cmd, svc.logger, stdOutListenFunction)

	if err == nil {
		if operation == "apply" {
			err := svc.setNewDeploymentState(ctx, ownerID, sinkId, "new")
			if err != nil {
				return err
			}
		} else if operation == "delete" {
			err := svc.setNewDeploymentState(ctx, ownerID, sinkId, "stopped")
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
		logger.Error("Collector Deploy Error", zap.Error(err))
	}
	err = cmd.Wait()
	if err != nil {
		logger.Error("Collector Deploy Error", zap.Error(err))
	}
	return stdoutScanner, stderrScanner, err
}

func (svc *deployService) getDeploymentState(ctx context.Context, ownerID, sinkId string) (string, error) {
	key := CollectorStatusKey + "." + ownerID + "." + sinkId
	cmd := svc.redisClient.Get(ctx, key)
	collectorStatusAsString, err := cmd.Result()
	if err != nil {
		return "", cmd.Err()
	}
	var collectorStatus CollectorStatusSortedSetEntry
	err = json.Unmarshal([]byte(collectorStatusAsString), &collectorStatus)
	if err != nil {
		return "", err
	}
	return collectorStatus.Status, nil
}

func (svc *deployService) setNewDeploymentState(ctx context.Context, ownerID, sinkId, state string) error {
	key := CollectorStatusKey + "." + ownerID + "." + sinkId
	entry := redis.Z{
		Score: float64(time.Now().Unix()),
		Member: CollectorStatusSortedSetEntry{
			SinkID:       sinkId,
			Status:       state,
			ErrorMessage: nil,
		},
	}
	intCmd := svc.redisClient.Set(ctx, key, &entry, 0)
	if intCmd.Err() != nil {
		return intCmd.Err()
	}
	return nil
}

func (svc *deployService) CreateOtelCollector(ctx context.Context, ownerID, sinkID, deploymentEntry string) error {
	err := svc.collectorDeploy(ctx, "apply", ownerID, sinkID, deploymentEntry)
	if err != nil {
		return err
	}

	return nil
}

func (svc *deployService) UpdateOtelCollector(ctx context.Context, ownerID, sinkID, deploymentEntry string) error {
	err := svc.DeleteOtelCollector(ctx, ownerID, sinkID, deploymentEntry)
	if err != nil {
		return err
	}
	err = svc.CreateOtelCollector(ctx, ownerID, sinkID, deploymentEntry)
	if err != nil {
		return err
	}
	return nil
}

func (svc *deployService) DeleteOtelCollector(ctx context.Context, ownerID, sinkID, deploymentEntry string) error {
	err := svc.collectorDeploy(ctx, "delete", ownerID, sinkID, deploymentEntry)
	if err != nil {
		return err
	}
	return nil
}
