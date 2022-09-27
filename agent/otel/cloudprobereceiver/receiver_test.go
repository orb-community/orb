package cloudprobereceiver_test

import (
	"context"
	"github.com/ns1labs/orb/agent/otel/pktvisorreceiver"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"testing"
)

func TestReceiver(t *testing.T) {
	f := pktvisorreceiver.NewFactory()
	tests := map[string]struct {
		useServiceAccount bool
		wantError         bool
	}{
		"success": {},
		"fails to get prometheus config": {
			useServiceAccount: true,
			wantError:         true,
		},
	}
	for desc, tt := range tests {
		t.Run(desc, func(t *testing.T) {
			cfg := (f.CreateDefaultConfig()).(*pktvisorreceiver.Config)
			cfg.UseServiceAccount = tt.useServiceAccount

			r, err := f.CreateMetricsReceiver(
				context.Background(),
				componenttest.NewNopReceiverCreateSettings(),
				cfg,
				consumertest.NewNop(),
			)

			if !tt.wantError {
				require.NoError(t, err)
				require.NotNil(t, r)

				require.NoError(t, r.Start(context.Background(), componenttest.NewNopHost()))
				require.NoError(t, r.Shutdown(context.Background()))
				return
			}

			require.Error(t, r.Start(context.Background(), componenttest.NewNopHost()))
		})
	}
}
