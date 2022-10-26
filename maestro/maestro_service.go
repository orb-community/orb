// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package maestro

import (
	"bufio"
	"context"
	"fmt"
	"github.com/ns1labs/orb/pkg/errors"
	"go.uber.org/zap"
	"os"
	"os/exec"
)

var (
	ErrCreateMaestro   = errors.New("failed to create Otel Collector")
	ErrConflictMaestro = errors.New("Otel collector already exists")
)

const namespace = "otelcollectors"

func (svc maestroService) collectorDeploy(ctx context.Context, operation, sinkId, manifest string) error {

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

func (svc maestroService) getConfigFromSinkId(config SinkConfig) (sinkID, sinkUrl, sinkUsername, sinkPassword string) {

	return config.Id, config.Url, config.Username, config.Password
}

func (svc maestroService) CreateOtelCollector(ctx context.Context, sinkID, deploymentEntry string) error {
	err := svc.collectorDeploy(ctx, "apply", sinkID, deploymentEntry)

	if err != nil {
		return err
	}
	return nil
}

func (svc maestroService) UpdateOtelCollector(ctx context.Context, sinkID, deploymentEntry string) error {
	err := svc.DeleteOtelCollector(ctx, sinkID)
	if err != nil {
		return err
	}
	err = svc.CreateOtelCollector(ctx, sinkID, deploymentEntry)
	if err != nil {
		return err
	}
	return nil
}

func (svc maestroService) DeleteOtelCollector(ctx context.Context, sinkID, deploymentEntry string) error {
	err := svc.collectorDeploy(ctx, "delete", sinkID, deploymentEntry)
	if err != nil {
		return err
	}
	return nil
}
