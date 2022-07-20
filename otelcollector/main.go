// Copyright 2019 OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Program otelcontribcol is an extension to the OpenTelemetry Collector
// that includes additional components, some vendor-specific, contributed
// from the wider community.

package otelcollector

import (
	"fmt"
	"github.com/ns1labs/orb/pkg/config"
	"go.uber.org/zap"
	"log"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"
)

type ComponentsFunc func(logger zap.Logger, svcCfg config.BaseSvcConfig, grpcCfgs []config.GRPCConfig) (component.Factories, error)

func RunWithComponents(logger zap.Logger, svcCfg config.BaseSvcConfig, grpcCfgs []config.GRPCConfig, componentsFunc ComponentsFunc) {
	factories, err := componentsFunc(logger, svcCfg, grpcCfgs)
	if err != nil {
		log.Fatalf("failed to build components: %v", err)
	}

	info := component.BuildInfo{
		Command:     "otelcollector",
		Description: "Otel Collector and Sinker",
		Version:     "latest",
	}

	if err = runInteractive(service.CollectorSettings{BuildInfo: info, Factories: factories}); err != nil {
		log.Fatal(err)
	}
}

func runInteractive(params service.CollectorSettings) error {
	cmd := service.NewCommand(params)
	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("collector server run finished with error: %w", err)
	}

	return nil
}
