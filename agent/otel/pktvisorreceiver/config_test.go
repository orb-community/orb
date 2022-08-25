package pktvisorreceiver_test

import (
	"github.com/ns1labs/orb/agent/otel/otlpmqttexporter"
	"github.com/ns1labs/orb/agent/otel/pktvisorreceiver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/service/servicetest"
	"path"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("load config", func(t *testing.T) {
		factories, err := componenttest.NopFactories()
		assert.NoError(t, err)

		factories.Receivers[typeStr] = pktvisorreceiver.NewFactory()
		factories.Exporters["otlpmqtt"] = otlpmqttexporter.NewFactory()
		cfgPath := path.Join(".", "testdata", "config.yaml")
		cfg, err := servicetest.LoadConfigAndValidate(cfgPath, factories)

		require.NoError(t, err)
		require.NotNil(t, cfg)
	})
}
