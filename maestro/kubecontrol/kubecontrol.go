package kubecontrol

import (
	"bufio"
	"context"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/plgd-dev/kit/v2/codec/json"
	"go.uber.org/zap"
	k8sappsv1 "k8s.io/api/apps/v1"
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
	DeleteOtelCollector(ctx context.Context, ownerID, sinkID string) error

	// UpdateOtelCollector - update an existing collector by id
	UpdateOtelCollector(ctx context.Context, ownerID, sinkID, deploymentEntry string) error
}

func (svc *deployService) removeDeployment(ctx context.Context, ownerID string, sinkId string) (error, bool) {
	deploymentName, status, err := svc.getDeploymentState(ctx, ownerID, sinkId)
	if err != nil {
		if status == "broken" {
			svc.logger.Info("get the deployed pod broke", zap.String("sinkID", sinkId), zap.Error(err))
			return nil, true
		}
	}
	if status == "deleted" {
		svc.logger.Info("Already deleted collector for Sink ID", zap.String("sinkID", sinkId))
		return nil, true
	}
	err = svc.clientSet.AppsV1().Deployments(namespace).Delete(ctx, deploymentName, k8smetav1.DeleteOptions{})
	if err != nil {
		svc.logger.Info("failed to remove deployment", zap.Error(err))
		return err, true
	}
	return nil, false
}

func (svc *deployService) applyDeployment(ctx context.Context, ownerID string, sinkId string, manifest string) (error, bool) {
	fileContent := []byte(manifest)
	tmp := strings.Split(string(fileContent), "\n")
	newContent := strings.Join(tmp[1:], "\n")
	_, status, err := svc.getDeploymentState(ctx, ownerID, sinkId)
	if err != nil {
		return err, false
	}
	fileName := "/tmp/otel-collector-" + sinkId + ".json"
	err = os.WriteFile(fileName, []byte(newContent), 0644)
	if err != nil {
		svc.logger.Error("failed to write file content", zap.Error(err))
		return err, true
	}
	if status == "active" {
		err := os.Remove(fileName)
		if err != nil {
			svc.logger.Info("failed to remove file", zap.Error(err))
			return err, true
		}
		svc.logger.Info("Already applied collector for Sink ID", zap.String("sinkID", sinkId))
		return nil, true
	}
	var applyConfig k8sv1.DeploymentApplyConfiguration
	err = json.Decode(fileContent, applyConfig)
	if err != nil {
		return err, true
	}
	_, err = svc.clientSet.AppsV1().Deployments(namespace).Apply(ctx, &applyConfig, k8smetav1.ApplyOptions{
		Force: true,
	})
	if err != nil {
		svc.logger.Info("failed to apply deployment", zap.Error(err))
		return err, true
	}
	// wait until deployment is active before returning
	for i := 0; i < 4; i++ {
		time.Sleep(1 * time.Second)
		_, status, err := svc.getDeploymentState(ctx, ownerID, sinkId)
		if status == "broken" || err != nil {
			svc.logger.Error("Failed during deployment, deployment is not active", zap.String("sinkID", sinkId), zap.Error(err))
			return err, false
		}
		if status == "active" {
			break
		}
	}
	return nil, false
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
	err2, done := svc.applyDeployment(ctx, ownerID, sinkID, deploymentEntry)
	if done {
		return err2
	}

	return nil
}

func (svc *deployService) UpdateOtelCollector(ctx context.Context, ownerID, sinkID, deploymentEntry string) error {
	err := svc.DeleteOtelCollector(ctx, ownerID, sinkID)
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

func (svc *deployService) DeleteOtelCollector(ctx context.Context, ownerID, sinkID string) error {
	err2, done := svc.removeDeployment(ctx, ownerID, sinkID)
	if done {
		return err2
	}
	return nil
}
