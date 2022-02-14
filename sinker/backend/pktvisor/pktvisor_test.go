package pktvisor_test

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/pb"
	"github.com/ns1labs/orb/sinker/backend"
	"github.com/ns1labs/orb/sinker/backend/pktvisor"
	"github.com/ns1labs/orb/sinker/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestDHCPConversion(t *testing.T) {
	var logger *zap.Logger
	pktvisor.Register(logger)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	agentID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var agent = &pb.OwnerRes{
		OwnerID:   ownerID.String(),
		AgentName: "agent-test",
	}

	data := fleet.AgentMetricsRPCPayload{
		PolicyID:   policyID.String(),
		PolicyName: "policy-test",
		Datasets:   nil,
		Format:     "json",
		BEVersion:  "1.0",
	}

	be := backend.GetBackend("pktvisor")

	cases := map[string]struct {
		data     []byte
		expected prometheus.TimeSeries
	}{
		"DHCPPayloadWirePacketsFiltered": {
			data: []byte(`
{
	"policy_dhcp": {
        "dhcp": {
            "wire_packets": {
                "filtered": 10
            }
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "dhcp_wire_packets_filtered",
					},
					{
						Name:  "instance",
						Value: "agent-test",
					},
					{
						Name:  "agent_id",
						Value: agentID.String(),
					},
					{
						Name:  "agent",
						Value: "agent-test",
					},
					{
						Name:  "policy_id",
						Value: policyID.String(),
					},
					{
						Name:  "policy",
						Value: "policy-test",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 10,
				},
			},
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			data.Data = c.data
			res, err := be.ProcessMetrics(agent, agentID.String(), data)
			require.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s", desc, err))
			var receivedLabel []prometheus.Label
			var receivedDatapoint prometheus.Datapoint
			for _, value := range res {
				if c.expected.Labels[0] == value.Labels[0] {
					receivedLabel = value.Labels
					receivedDatapoint = value.Datapoint
				}
			}
			assert.True(t, reflect.DeepEqual(c.expected.Labels, receivedLabel), fmt.Sprintf("%s: expected %v got %v", desc, c.expected.Labels, receivedLabel))
			assert.Equal(t, c.expected.Datapoint.Value, receivedDatapoint.Value, fmt.Sprintf("%s: expected value %f got %f", desc, c.expected.Datapoint.Value, receivedDatapoint.Value))
		})
	}

}

func TestASNConversion(t *testing.T) {
	var logger *zap.Logger
	pktvisor.Register(logger)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	agentID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var agent = &pb.OwnerRes{
		OwnerID:   ownerID.String(),
		AgentName: "agent-test",
	}

	data := fleet.AgentMetricsRPCPayload{
		PolicyID:   policyID.String(),
		PolicyName: "policy-test",
		Datasets:   nil,
		Format:     "json",
		BEVersion:  "1.0",
	}

	be := backend.GetBackend("pktvisor")

	cases := map[string]struct {
		data     []byte
		expected prometheus.TimeSeries
	}{
		"PacketPayloadTopASN": {
			data: []byte(`
{
    "policy_packets": {
        "packets": {
            "top_ASN": [
                {
                    "estimate": 996,
                    "name": "36236/NETACTUATE"
                }
            ]
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "packets_top_ASN",
					},
					{
						Name:  "instance",
						Value: "agent-test",
					},
					{
						Name:  "agent_id",
						Value: agentID.String(),
					},
					{
						Name:  "agent",
						Value: "agent-test",
					},
					{
						Name:  "policy_id",
						Value: policyID.String(),
					},
					{
						Name:  "policy",
						Value: "policy-test",
					},
					{
						Name:  "asn",
						Value: "36236/NETACTUATE",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 996,
				},
			},
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			data.Data = c.data
			res, err := be.ProcessMetrics(agent, agentID.String(), data)
			require.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s", desc, err))
			var receivedLabel []prometheus.Label
			var receivedDatapoint prometheus.Datapoint
			for _, value := range res {
				if c.expected.Labels[0] == value.Labels[0] {
					receivedLabel = value.Labels
					receivedDatapoint = value.Datapoint
				}
			}
			assert.True(t, reflect.DeepEqual(c.expected.Labels, receivedLabel), fmt.Sprintf("%s: expected %v got %v", desc, c.expected.Labels, receivedLabel))
			assert.Equal(t, c.expected.Datapoint.Value, receivedDatapoint.Value, fmt.Sprintf("%s: expected value %f got %f", desc, c.expected.Datapoint.Value, receivedDatapoint.Value))
		})
	}

}

func TestGeoLocConversion(t *testing.T) {
	var logger *zap.Logger
	pktvisor.Register(logger)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	agentID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var agent = &pb.OwnerRes{
		OwnerID:   ownerID.String(),
		AgentName: "agent-test",
	}

	data := fleet.AgentMetricsRPCPayload{
		PolicyID:   policyID.String(),
		PolicyName: "policy-test",
		Datasets:   nil,
		Format:     "json",
		BEVersion:  "1.0",
	}

	be := backend.GetBackend("pktvisor")

	cases := map[string]struct {
		data     []byte
		expected prometheus.TimeSeries
	}{
		"PacketPayloadTopASN": {
			data: []byte(`
{
    "policy_packets": {
        "packets": {
            "top_geoLoc": [
                {
                    "estimate": 996,
                    "name": "AS/Hong Kong/HCW/Central"
                }
            ]
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "packets_top_geoLoc",
					},
					{
						Name:  "instance",
						Value: "agent-test",
					},
					{
						Name:  "agent_id",
						Value: agentID.String(),
					},
					{
						Name:  "agent",
						Value: "agent-test",
					},
					{
						Name:  "policy_id",
						Value: policyID.String(),
					},
					{
						Name:  "policy",
						Value: "policy-test",
					},
					{
						Name:  "geo_loc",
						Value: "AS/Hong Kong/HCW/Central",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 996,
				},
			},
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			data.Data = c.data
			res, err := be.ProcessMetrics(agent, agentID.String(), data)
			require.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s", desc, err))
			var receivedLabel []prometheus.Label
			var receivedDatapoint prometheus.Datapoint
			for _, value := range res {
				if c.expected.Labels[0] == value.Labels[0] {
					receivedLabel = value.Labels
					receivedDatapoint = value.Datapoint
				}
			}
			assert.True(t, reflect.DeepEqual(c.expected.Labels, receivedLabel), fmt.Sprintf("%s: expected %v got %v", desc, c.expected.Labels, receivedLabel))
			assert.Equal(t, c.expected.Datapoint.Value, receivedDatapoint.Value, fmt.Sprintf("%s: expected value %f got %f", desc, c.expected.Datapoint.Value, receivedDatapoint.Value))
		})
	}

}
