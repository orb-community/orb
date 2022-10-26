package kubecontrol

import (
	"bufio"
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/exec"
)

const namespace = "otelcollectors"

var _ Service = (*deployService)(nil)

type deployService struct {
	logger *zap.Logger
}

func NewService(logger *zap.Logger) Service {
	return &deployService{logger: logger}
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
	err := os.WriteFile("/tmp/otel-collector-"+sinkId+".json", fileContent, 0644)
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
	err := svc.CreateOtelCollector(ctx, sinkID, deploymentEntry)
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
