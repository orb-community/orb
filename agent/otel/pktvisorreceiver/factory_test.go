package pktvisorreceiver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestFactory(t *testing.T) {
	f := NewFactory()
	cfg := f.CreateDefaultConfig()
	require.NotNil(t, cfg)

	receiver, err := f.CreateMetricsReceiver(
		context.Background(),
		receivertest.NewNopCreateSettings(),
		cfg,
		consumertest.NewNop(),
	)
	require.NoError(t, err)
	require.NotNil(t, receiver)

}
