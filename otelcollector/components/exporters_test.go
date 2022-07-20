// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Skip tests on Windows temporarily, see https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/11451
//go:build !windows
// +build !windows

package components

import (
	"context"
	"errors"
	testutil "github.com/ns1labs/orb/otelcollector/testutil"
	"go.uber.org/zap"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter"
)

func TestDefaultExporters(t *testing.T) {
	logger := zap.NewNop()
	factories, err := Components(*logger)
	assert.NoError(t, err)

	expFactories := factories.Exporters
	endpoint := testutil.GetAvailableLocalAddress(t)

	tests := []struct {
		exporter      config.Type
		getConfigFn   getExporterConfigFn
		skipLifecycle bool
	}{
		{
			exporter:      "logging",
			skipLifecycle: runtime.GOOS == "darwin", // TODO: investigate why this fails on darwin.
		},
		{
			exporter: "otlp",
			getConfigFn: func() config.Exporter {
				cfg := expFactories["otlp"].CreateDefaultConfig().(*otlpexporter.Config)
				cfg.GRPCClientSettings = configgrpc.GRPCClientSettings{
					Endpoint: endpoint,
				}
				return cfg
			},
		},
		{
			exporter: "otlphttp",
			getConfigFn: func() config.Exporter {
				cfg := expFactories["otlphttp"].CreateDefaultConfig().(*otlphttpexporter.Config)
				cfg.Endpoint = "http://" + endpoint
				return cfg
			},
		},
		{
			exporter: "prometheus",
			getConfigFn: func() config.Exporter {
				cfg := expFactories["prometheus"].CreateDefaultConfig().(*prometheusexporter.Config)
				cfg.Endpoint = endpoint
				return cfg
			},
		},
		{
			exporter: "prometheusremotewrite",
		},
		{
			exporter:      "googlecloud",
			skipLifecycle: true, // Requires credentials to be able to successfully load the exporter
		},
		{
			exporter:      "googlemanagedprometheus",
			skipLifecycle: true, // Requires credentials to be able to successfully load the exporter
		},
		{
			exporter: "googlecloudpubsub",
		},
	}

	assert.Len(t, tests, len(expFactories), "All user configurable components must be added to the lifecycle test")
	for _, tt := range tests {
		t.Run(string(tt.exporter), func(t *testing.T) {
			t.Parallel()

			factory, ok := expFactories[tt.exporter]
			require.True(t, ok)
			assert.Equal(t, tt.exporter, factory.Type())
			assert.Equal(t, config.NewComponentID(tt.exporter), factory.CreateDefaultConfig().ID())

			if tt.skipLifecycle {
				t.Skip("Skipping lifecycle test", tt.exporter)
				return
			}

			verifyExporterLifecycle(t, factory, tt.getConfigFn)
		})
	}
}

// GetExporterConfigFn is used customize the configuration passed to the verification.
// This is used to change ports or provide values required but not provided by the
// default configuration.
type getExporterConfigFn func() config.Exporter

// verifyExporterLifecycle is used to test if an exporter type can handle the typical
// lifecycle of a component. The getConfigFn parameter only need to be specified if
// the test can't be done with the default configuration for the component.
func verifyExporterLifecycle(t *testing.T, factory component.ExporterFactory, getConfigFn getExporterConfigFn) {
	ctx := context.Background()
	host := newAssertNoErrorHost(t)
	expCreateSettings := componenttest.NewNopExporterCreateSettings()

	cfg := factory.CreateDefaultConfig()
	if getConfigFn != nil {
		cfg = getConfigFn()
	}

	createFns := []createExporterFn{
		wrapCreateLogsExp(factory),
		wrapCreateTracesExp(factory),
		wrapCreateMetricsExp(factory),
	}

	for i := 0; i < 2; i++ {
		var exps []component.Exporter
		for _, createFn := range createFns {
			exp, err := createFn(ctx, expCreateSettings, cfg)
			if errors.Is(err, component.ErrDataTypeIsNotSupported) {
				continue
			}
			require.NoError(t, err)
			require.NoError(t, exp.Start(ctx, host))
			exps = append(exps, exp)
		}
		for _, exp := range exps {
			assert.NoError(t, exp.Shutdown(ctx))
		}
	}
}

type createExporterFn func(
	ctx context.Context,
	set component.ExporterCreateSettings,
	cfg config.Exporter,
) (component.Exporter, error)

func wrapCreateLogsExp(factory component.ExporterFactory) createExporterFn {
	return func(ctx context.Context, set component.ExporterCreateSettings, cfg config.Exporter) (component.Exporter, error) {
		return factory.CreateLogsExporter(ctx, set, cfg)
	}
}

func wrapCreateTracesExp(factory component.ExporterFactory) createExporterFn {
	return func(ctx context.Context, set component.ExporterCreateSettings, cfg config.Exporter) (component.Exporter, error) {
		return factory.CreateTracesExporter(ctx, set, cfg)
	}
}

func wrapCreateMetricsExp(factory component.ExporterFactory) createExporterFn {
	return func(ctx context.Context, set component.ExporterCreateSettings, cfg config.Exporter) (component.Exporter, error) {
		return factory.CreateMetricsExporter(ctx, set, cfg)
	}
}
