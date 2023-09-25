package kubecontrol

import (
	"bufio"
	"context"
	"fmt"
	_ "github.com/orb-community/orb/maestro/config"
	"github.com/orb-community/orb/pkg/errors"
	"go.uber.org/zap"
	k8sappsv1 "k8s.io/api/apps/v1"
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
	logger    *zap.Logger
	clientSet *kubernetes.Clientset
}

const OperationDeploy CollectorOperation = iota
const OperationDelete = 1

type CollectorOperation int

func (o CollectorOperation) Name() string {
	switch o {
	case OperationDeploy:
		return "deploy"
	case OperationDelete:
		return "delete"
	default:
		return "unknown"
	}
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
	CreateOtelCollector(ctx context.Context, ownerID, sinkID, deploymentEntry string) (string, error)

	// KillOtelCollector - kill an existing collector by id, terminating by the ownerID, sinkID without the file
	KillOtelCollector(ctx context.Context, deploymentName, sinkID string) error
}

func (svc *deployService) collectorDeploy(ctx context.Context, operation, ownerID, sinkId, manifest string) (string, error) {
	_, status, err := svc.getDeploymentState(ctx, ownerID, sinkId)
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
		return "", err
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
	// TODO this will be retrieved once we move to K8s SDK
	collectorName := fmt.Sprintf("otelcol-%s-%s", ownerID, sinkId)
	return collectorName, nil
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

func (svc *deployService) CreateOtelCollector(ctx context.Context, ownerID, sinkID, deploymentEntry string) (string, error) {
	col, err := svc.collectorDeploy(ctx, "apply", ownerID, sinkID, deploymentEntry)
	if err != nil {
		return "", err
	}

	return col, nil
}

func (svc *deployService) KillOtelCollector(ctx context.Context, deploymentName string, sinkId string) error {
	stdOutListenFunction := func(out *bufio.Scanner, err *bufio.Scanner) {
		for out.Scan() {
			svc.logger.Info("Deploy Info: " + out.Text())
		}
		for err.Scan() {
			svc.logger.Info("Deploy Error: " + err.Text())
		}
	}

	// execute action
	cmd := exec.Command("kubectl", "delete", "deploy", deploymentName, "-n", namespace)
	_, _, err := execCmd(ctx, cmd, svc.logger, stdOutListenFunction)
	if err == nil {
		svc.logger.Info(fmt.Sprintf("successfully killed the otel-collector for sink-id: %s", sinkId))
	}

	return nil
}
