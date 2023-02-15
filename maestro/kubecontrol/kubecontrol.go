package kubecontrol

import (
	"bufio"
	"context"
	"fmt"
	_ "github.com/ns1labs/orb/maestro/config"
	"github.com/ns1labs/orb/pkg/errors"
	"go.uber.org/zap"
	k8sappsv1 "k8s.io/api/apps/v1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		logger.Error("error on get cluster config", zap.Error(err))
		return nil
	}
	clientSet, err := kubernetes.NewForConfig(clusterConfig)
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
	_, status, err := svc.getDeploymentState(ctx, ownerID, sinkId)
	if operation == "apply" {
		fileContent := []byte(manifest)
		tmp := strings.Split(string(fileContent), "\n")
		newContent := strings.Join(tmp[1:], "\n")
		if err != nil {
			if status == "broken" {
				operation = "delete"
			}
		}
		err = os.WriteFile("/tmp/otel-collector-"+sinkId+".json", []byte(newContent), 0644)
		if err != nil {
			svc.logger.Error("failed to write file content", zap.Error(err))
			return err
		}
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

		} else if operation == "delete" {

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
		deploymentList, err2 := svc.clientSet.AppsV1().Deployments(namespace).List(ctx, k8smetav1.ListOptions{})
		if err2 != nil {
			svc.logger.Error("error on reading pods", zap.Error(err2))
			return "", "", err2
		}
		for _, deployment := range deploymentList.Items {
			if strings.Contains(deployment.Name, sinkId) {
				svc.logger.Info("found deployment for sink")
				deploymentName = deployment.Name
				if len(deployment.Status.Conditions) == 0 || deployment.Status.Conditions[0].Type == k8sappsv1.DeploymentReplicaFailure {
					svc.logger.Error("error on retrieving collector, deployment is broken")
					return "", "broken", errors.New("error on retrieving collector, deployment is broken")
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
	// Time to wait until K8s completely removes before re-creating
	time.Sleep(3 * time.Second)
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
