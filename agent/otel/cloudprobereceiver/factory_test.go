package cloudprobereceiver_test

import (
	"context"
	"github.com/ns1labs/orb/agent/otel/pktvisorreceiver"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"testing"
)

func TestFactory(t *testing.T) {
	f := pktvisorreceiver.NewFactory()
	cfg := f.CreateDefaultConfig()
	require.NotNil(t, cfg)

	receiver, err := f.CreateMetricsReceiver(
		context.Background(),
		componenttest.NewNopReceiverCreateSettings(),
		cfg,
		consumertest.NewNop(),
	)
	require.NoError(t, err)
	require.NotNil(t, receiver)

}
