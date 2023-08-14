package otel

import (
	"bufio"
	"context"
	"fmt"
	"github.com/amenzhinsky/go-memexec"
	"github.com/orb-community/orb/agent/policies"
	"go.uber.org/zap"
	"os"
)

func (o openTelemetryBackend) ApplyPolicy(newPolicyData policies.PolicyData, updatePolicy bool) error {
	if !updatePolicy || !o.policyRepo.Exists(newPolicyData.ID) {
		temporaryFile, err := os.CreateTemp(o.policyConfigDirectory, fmt.Sprintf("otel-%s-config.yml", newPolicyData.ID))
		if err != nil {
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
			dataAsByte := []byte(newPolicyData.Data.(string))
			currentPolicyPath := o.policyConfigDirectory + fmt.Sprintf("otel-%s-config.yml", currentPolicyData.ID)
			o.logger.Info("new policy version received, updating", zap.String("policy_id", newPolicyData.ID), zap.Int32("version", newPolicyData.Version))
			err := os.WriteFile(currentPolicyPath, dataAsByte, os.ModeTemporary)
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

func (o openTelemetryBackend) addRunner(policyData policies.PolicyData, policyFilePath string) error {
	policyContext := context.WithValue(o.mainContext, "policy_id", policyData.ID)
	executable, err := memexec.New(openTelemetryContribBinary)
	if err != nil {
		return err
	}
	defer func(executable *memexec.Exec) {
		err := executable.Close()
		if err != nil {
			o.logger.Error("error closing executable", zap.Error(err))
		}
	}(executable)
	command := executable.CommandContext(policyContext, "--config", policyFilePath)
	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}
	go func(ctx context.Context) {
		err := command.Start()
		if err != nil {
			o.logger.Error("error starting command", zap.Error(err))
			ctx.Done()
			return
		}
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			o.logger.Info("stderr output",
				zap.String("policy_id", policyData.ID),
				zap.String("line", scanner.Text()))
			if command.Err != nil {
				o.logger.Error("error running command", zap.Error(command.Err))
				ctx.Done()
				return
			}
		}
	}(policyContext)
	o.addPolicyControl(policyContext, policyData.ID)

	return nil
}

func (o openTelemetryBackend) addPolicyControl(policyCtx context.Context, policyID string) {
	o.runningCollectors[policyID] = policyCtx
}

func (o openTelemetryBackend) removePolicyControl(policyID string) {
	o.runningCollectors[policyID].Done()
	delete(o.runningCollectors, policyID)
}

func (o openTelemetryBackend) RemovePolicy(data policies.PolicyData) error {
	if o.policyRepo.Exists(data.ID) {
		o.removePolicyControl(data.ID)
		err := o.policyRepo.Remove(data.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
