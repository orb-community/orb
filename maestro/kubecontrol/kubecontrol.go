package kubecontrol

import (
	"bufio"
	"context"
	"fmt"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/plgd-dev/kit/v2/codec/json"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"os/exec"
	"strings"
	"time"
)

const namespace = "otelcollectors"

var _ Service = (*deployService)(nil)

type deployService struct {
	logger    *zap.Logger
	clientSet *kubernetes.Clientset
}

func NewService(logger *zap.Logger) Service {
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Error("error on get cluster config", zap.Error(err))
		return nil
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error("error on get client", zap.Error(err))
		return nil
	}
	return &deployService{logger: logger, clientSet: clientSet}
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

	_, status, err := svc.getDeploymentState(ctx, ownerID, sinkId)
	if err != nil {
		if status == "broken" {
			operation = "delete"
		}
	}
	if operation == "apply" {
		if status == "active" {
			svc.logger.Info("Already applied Sink ID=" + sinkId)
			return nil
		}
	} else if operation == "delete" {
		if status == "deleted" {
			svc.logger.Info("Already deleted Sink ID=" + sinkId)
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
		svc.logger.Info(fmt.Sprintf("successfully %s the otel-collector for sink-id: %s", operation, sinkId))
	}
	svc.logger.Info(fmt.Sprintf("successfully %s the otel-collector for sink-id: %s", operation, sinkId))

	return nil
}

func (svc *deployService) newCollectorDeploy(ctx context.Context, operation, ownerID, sinkId, manifest string) error {
	fileContent := []byte(manifest)
	tmp := strings.Split(string(fileContent), "\n")
	newContent := strings.Join(tmp[1:], "\n")
	deploymentName, status, err := svc.getDeploymentState(ctx, ownerID, sinkId)
	if err != nil {
		if status == "broken" {
			operation = "delete"
		}
	}
	fileName := "/tmp/otel-collector-" + sinkId + ".json"
	err = os.WriteFile(fileName, []byte(newContent), 0644)
	if err != nil {
		svc.logger.Error("failed to write file content", zap.Error(err))
		return err
	}
	if operation == "apply" {
		if status == "active" {
			err := os.Remove(fileName)
			if err != nil {
				svc.logger.Info("failed to remove file", zap.Error(err))
				return err
			}
			svc.logger.Info("Already applied collector for Sink ID", zap.String("sinkID", sinkId))
			return nil
		}
		var applyConfig k8sv1.DeploymentApplyConfiguration
		err := json.Decode(fileContent, applyConfig)
		if err != nil {
			return err
		}
		svc.logger.Info("DEBUG applyConfig", zap.Any("applyConfig", applyConfig))
		result, err := svc.clientSet.AppsV1().Deployments(namespace).Apply(ctx, &applyConfig, k8smetav1.ApplyOptions{
			Force: true,
		})
		if err != nil {
			svc.logger.Info("failed to apply deployment", zap.Error(err))
			return err
		}
		svc.logger.Info("DEBUG result", zap.Any("result", result))
	} else if operation == "delete" {
		if status == "deleted" {
			svc.logger.Info("Already deleted collector for Sink ID", zap.String("sinkID", sinkId))
			return nil
		}
		err := svc.clientSet.AppsV1().Deployments(namespace).Delete(ctx, deploymentName, k8smetav1.DeleteOptions{})
		if err != nil {
			svc.logger.Info("failed to remove deployment", zap.Error(err))
			return err
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

func (svc *deployService) getDeploymentState(ctx context.Context, _, sinkId string) (deploymentName string, status string, err error) {
	// Since this can take a while to be retrieved, we need to have a wait mechanism
	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		pods, err2 := svc.clientSet.CoreV1().Pods(namespace).List(ctx, k8smetav1.ListOptions{})
		if err2 != nil {
			svc.logger.Error("error on reading pods", zap.Error(err2))
			return "", "", err2
		}
		for _, pod := range pods.Items {
			if strings.Contains(pod.Name, sinkId) {
				deploymentName = pod.Name
				if pod.Status.Phase == v1.PodFailed {
					svc.logger.Error("error on retrieving collector, pod is broken")
					return "", "broken", errors.New(pod.Status.Message)
				}
				if pod.Status.Phase != v1.PodRunning {
					break
				}
				status = "active"
				return
			}
		}
	}
	status = "deleted"
	return "", "deleted", nil
}

func (svc *deployService) CreateOtelCollector(ctx context.Context, ownerID, sinkID, deploymentEntry string) error {
	err := svc.newCollectorDeploy(ctx, "apply", ownerID, sinkID, deploymentEntry)
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
	// Time to wait until K8s completely removes before re-creating
	time.Sleep(3 * time.Second)
	err = svc.CreateOtelCollector(ctx, ownerID, sinkID, deploymentEntry)
	if err != nil {
		return err
	}
	return nil
}

func (svc *deployService) DeleteOtelCollector(ctx context.Context, ownerID, sinkID, deploymentEntry string) error {
	err := svc.newCollectorDeploy(ctx, "delete", ownerID, sinkID, deploymentEntry)
	if err != nil {
		return err
	}
	return nil
}
