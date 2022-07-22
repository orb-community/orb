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

// Program otelcollector is an extension to the OpenTelemetry Collector
// that includes additional components, some vendor-specific, contributed
// from the wider community.

package otelcollector

import (
	"context"
	"github.com/ns1labs/orb/otelcollector/components"
	"github.com/ns1labs/orb/pkg/config"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/service"
	"go.uber.org/zap"
)

func StartCollector(ctx context.Context, logger zap.Logger, svcCfg config.BaseSvcConfig, sinkerGrpcCfg,
	policiesGrpcCfg, sinksGrpcCfg config.GRPCConfig) {
	// define factories
	factories, err := components.Components(logger)
	if err != nil {
		ctx.Done()
		logger.Fatal("failed to build components", zap.Error(err))
	}

	orbConfigProvider, err := getConfigProvider(svcCfg, sinkerGrpcCfg)
	if err != nil {
		ctx.Done()
		logger.Fatal("failed during build of config provider", zap.Error(err))
	}
	logger.Info("config provider load successfully")

	// dataCollector firstly starting the collector which will receive the data from sinker
	dataCollector, err := service.New(service.CollectorSettings{
		Factories: factories,
		BuildInfo: component.BuildInfo{
			Description: "DataEntry",
			Version:     "alpha",
		},
		DisableGracefulShutdown: false,
		ConfigProvider:          orbConfigProvider,
		LoggingOptions:          nil,
		SkipSettingGRPCLogger:   false,
	})
	if err != nil {
		ctx.Done()
		logger.Fatal("fatal error during data collector initialization", zap.Error(err))
	}
	logger.Info("started dataCollector successfully", zap.String("state", dataCollector.GetState().String()))
	dataColCtx := context.WithValue(ctx, "collector", "main")
	err = dataCollector.Run(dataColCtx)
	if err != nil {
		ctx.Done()
		dataColCtx.Done()
		logger.Fatal("fatal error during data collector execution", zap.Error(err))
	}
}

func getConfigProvider(svcCfg config.BaseSvcConfig, sinkerGrpcCfg config.GRPCConfig) (service.ConfigProvider, error) {
	// Figure out to create a provider that fetches information from a Database/Cache
	// Make it watch changes based on the GRPC call to change from Datasets
	return service.NewConfigProvider(service.ConfigProviderSettings{
		Locations:     []string{""},
		MapProviders:  map[string]confmap.Provider{},
		MapConverters: nil,
	})
}
