package otel

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-cmd/cmd"
	"github.com/orb-community/orb/agent/policies"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

const tempFileNamePattern = "otel-%s-config.yml"

type runningPolicy struct {
	ctx        context.Context
	cancel     context.CancelFunc
	policyId   string
	policyData policies.PolicyData
	statusChan *cmd.Status
}

func (o *openTelemetryBackend) ApplyPolicy(newPolicyData policies.PolicyData, updatePolicy bool) error {
	policyYaml, err := yaml.Marshal(newPolicyData.Data)
	if err != nil {
		o.logger.Warn("yaml policy marshal failure", zap.String("policy_id", newPolicyData.ID), zap.Any("policy", newPolicyData.Data))
		return err
	}
	err = o.ValidatePolicy(string(policyYaml))
	if err != nil {
		return err
	}
	if !updatePolicy || !o.policyRepo.Exists(newPolicyData.ID) {
		temporaryFile, err := os.CreateTemp(o.policyConfigDirectory, fmt.Sprintf(tempFileNamePattern, newPolicyData.ID))
		if err != nil {
			o.logger.Error("failed to create temporary file", zap.Error(err), zap.String("policy_id", newPolicyData.ID))
			return err
		}
		o.logger.Debug("writing policy to temporary file", zap.String("policy_id", newPolicyData.ID), zap.String("policyData", string(policyYaml)))
		_, err = temporaryFile.Write(policyYaml)
		if err != nil {
			o.logger.Error("failed to write temporary file", zap.Error(err), zap.String("policy_id", newPolicyData.ID))
			return err
		}
		err = o.addRunner(newPolicyData, temporaryFile.Name())
		if err != nil {
			return err
		}
	} else {
		currentPolicyData, err := o.policyRepo.Get(newPolicyData.ID)
		if err != nil {
			return err
		}
		if currentPolicyData.Version <= newPolicyData.Version {
			currentPolicyPath := o.policyConfigDirectory + fmt.Sprintf(tempFileNamePattern, currentPolicyData.ID)
			o.logger.Info("new policy version received, updating",
				zap.String("policy_id", newPolicyData.ID),
				zap.Int32("version", newPolicyData.Version))
			err := os.WriteFile(currentPolicyPath, []byte(policyYaml), os.ModeTemporary)
			if err != nil {
				return err
			}
			err = o.policyRepo.Update(newPolicyData)
			if err != nil {
				return err
			}
			o.otelReceiverTaps = append(o.otelReceiverTaps, newPolicyData.ID)
		} else {
			o.logger.Info("current policy version is newer than the one being applied, skipping",
				zap.String("policy_id", newPolicyData.ID),
				zap.Int32("current_version", currentPolicyData.Version),
				zap.Int32("incoming_version", newPolicyData.Version))
		}
	}

	return nil
}

func (o *openTelemetryBackend) addRunner(policyData policies.PolicyData, policyFilePath string) error {
	policyContext, policyCancel := context.WithCancel(context.WithValue(o.mainContext, "policy_id", policyData.ID))
	command := cmd.NewCmdOptions(cmd.Options{Buffered: false, Streaming: true}, o.otelExecutablePath, "--config", policyFilePath)
	go func(ctx context.Context, logger *zap.Logger) {
		status := command.Start()
		o.logger.Info("starting otel policy", zap.String("policy_id", policyData.ID),
			zap.Any("status", command.Status()), zap.Int("process id", command.Status().PID))
		for command.Status().Complete == false {
			select {
			case <-ctx.Done():
				err := command.Stop()
				if err != nil {
					logger.Error("failed to stop otel", zap.String("policy_id", policyData.ID), zap.Error(err))
					return
				}
			case line := <-command.Stdout:
				logger.Info("otel stdout", zap.String("policy_id", policyData.ID), zap.String("line", line))
			case line := <-command.Stderr:
				logger.Warn("otel stderr", zap.String("policy_id", policyData.ID), zap.String("line", line))
			case finalStatus := <-status:
				logger.Info("otel finished", zap.String("policy_id", policyData.ID), zap.Any("status", finalStatus))
			}
		}
	}(policyContext, o.logger)
	status := command.Status()
	policyEntry := runningPolicy{
		ctx:        policyContext,
		cancel:     policyCancel,
		policyId:   policyData.ID,
		policyData: policyData,
		statusChan: &status,
	}
	o.addPolicyControl(policyEntry, policyData.ID)

	return nil
}

func (o *openTelemetryBackend) addPolicyControl(policyEntry runningPolicy, policyID string) {
	o.runningCollectors[policyID] = policyEntry
}

func (o *openTelemetryBackend) removePolicyControl(policyID string) {
	policy, ok := o.runningCollectors[policyID]
	if !ok {
		o.logger.Error("did not find a running collector for policy id", zap.String("policy_id", policyID))
		return
	}
	policy.cancel()
	select {
	case <-policy.ctx.Done():
		o.logger.Info("policy context done", zap.String("policy_id", policyID))
	}
	delete(o.runningCollectors, policyID)
}

func (o *openTelemetryBackend) RemovePolicy(data policies.PolicyData) error {
	if o.policyRepo.Exists(data.ID) {
		o.removePolicyControl(data.ID)
		err := o.policyRepo.Remove(data.ID)
		if err != nil {
			return err
		}
		policyPath := o.policyConfigDirectory + fmt.Sprintf(tempFileNamePattern, data.ID)
		err = os.Remove(policyPath)
		if err != nil {
			o.logger.Error("error removing temporary file with policy")
			return err
		}
	}
	return nil
}

func (o *openTelemetryBackend) ValidatePolicy(policyData string) error {
	if policyData == "" {
		return errors.New("policy data is empty")
	}
	return nil
}
