package kubecontrol

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

const namespace = "otelcollectors"

var _ Service = (*deployService)(nil)

type deployService struct {
	logger          *zap.Logger
	deploymentState map[string]bool
}

func NewService(logger *zap.Logger) Service {
	deploymentState := make(map[string]bool)
	return &deployService{logger: logger, deploymentState: deploymentState}
}

type Service interface {
	// CreateOtelCollector - create an existing collector by id
	CreateOtelCollector(ctx context.Context, sinkID, deploymentEntry string) error

	// DeleteOtelCollector - delete an existing collector by id
	DeleteOtelCollector(ctx context.Context, sinkID, deploymentEntry string) error

	// UpdateOtelCollector - update an existing collector by id
	UpdateOtelCollector(ctx context.Context, sinkID, deploymentEntry string) error
}

func (svc *deployService) collectorDeploy(_ context.Context, operation, sinkId, manifest string) error {

	fileContent := []byte(manifest)
	tmp := strings.Split(string(fileContent), "\n")
	newContent := strings.Join(tmp[1:], "\n")

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

	// execute action
	cmd := exec.Command("kubectl", operation, "-f", "/tmp/otel-collector-"+sinkId+".json", "-n", namespace)
	stdoutReader, _ := cmd.StdoutPipe()
	stdoutScanner := bufio.NewScanner(stdoutReader)
	go func() {
		for stdoutScanner.Scan() {
			fmt.Println(stdoutScanner.Text())
			svc.logger.Info("Deploy Info: " + stdoutScanner.Text())
		}
	}()
	stderrReader, _ := cmd.StderrPipe()
	stderrScanner := bufio.NewScanner(stderrReader)
	go func() {
		for stderrScanner.Scan() {
			fmt.Println(stderrScanner.Text())
			svc.logger.Info("Deploy Error: " + stderrScanner.Text())
		}
	}()
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error : %v \n", err)
		svc.logger.Error("Collector Deploy Error", zap.Error(err))
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		svc.logger.Error("Collector Deploy Error", zap.Error(err))
	}

	if err == nil {
		if operation == "apply" {
			svc.deploymentState[sinkId] = true
		} else if operation == "delete" {
			svc.deploymentState[sinkId] = false
		}
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
