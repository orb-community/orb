package otlpexporter_test

import (
	"github.com/ns1labs/orb/agent/otel/otlpexporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/service/servicetest"
	"path"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Skip("Configuration is not done yet")
	factories, err := componenttest.NopFactories()
	assert.NoError(t, err)

	factory := otlpexporter.NewFactory()
	factories.Exporters["otlp"] = factory

	cfg, err := servicetest.LoadConfigAndValidate(path.Join(".", "testdata", "config.yaml"), factories)
	require.NoError(t, err)
	require.NotNil(t, cfg)
}
