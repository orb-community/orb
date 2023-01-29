package kubecontrol

import (
	"bufio"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"os/exec"
	"strings"
)

const namespace = "otelcollectors"

var _ Service = (*deployService)(nil)

type deployService struct {
	logger      *zap.Logger
	redisClient *redis.Client
	deploymentState map[string]bool
}

func NewService(logger *zap.Logger, redisClient *redis.Client) Service {
	deploymentState := make(map[string]bool)
	return &deployService{logger: logger, redisClient: redisClient, deploymentState: deploymentState}
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
	
	// due the delay in k8s api to answer, we should control it internally using a map or redis for more maestro replicas, and not with k8s
	// otherwise can generate wrong deployment status during sink update operation
	if operation == "apply" {
		if value, ok := svc.deploymentState[sinkId]; ok && value {
			svc.logger.Info("Already applied Sink ID=" + sinkId)
			return nil
		}
	} else if operation == "delete" {
		if value, ok := svc.deploymentState[sinkId]; ok && !value {
			svc.logger.Info("Already deleted Sink ID=" + sinkId)
			return nil
		}
	}

	err := os.WriteFile("/tmp/otel-collector-"+sinkId+".json", []byte(newContent), 0644)
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
		svc.logger.Info(fmt.Sprintf("successfully %s the otel-collector for sink-id: %s", operation, sinkId))
		// update deployment state map
		if operation == "apply" {
			svc.deploymentState[sinkId] = true
		} else if operation == "delete" {
			svc.deploymentState[sinkId] = false
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

func (svc *deployService) getDeploymentState(ctx context.Context, _, sinkId string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		svc.logger.Error("error on get cluster config", zap.Error(err))
		return "", err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		svc.logger.Error("error on get client", zap.Error(err))
		return "", err
	}
	pods, err := clientSet.CoreV1().Pods(namespace).List(ctx, k8smetav1.ListOptions{})
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, sinkId) {
			return "active", nil
		}
	}
	return "deleted", nil
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
