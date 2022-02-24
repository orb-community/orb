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
		"PacketPayloadTopGeoLoc": {
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

func TestPCAPConversion(t *testing.T) {
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
		"PCAPPayload_Tcp_Reassembly_Errors": {
			data: []byte(`
{
	"policy_pcap": {
        "pcap": {
            "tcp_reassembly_errors": 2
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "pcap_tcp_reassembly_errors",
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
					Value: 2,
				},
			},
		},
		"PCAPPayload_if_drops": {
			data: []byte(`
{
	"policy_pcap": {
		"pcap": {
			"if_drops": 2
    	}
	}
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "pcap_if_drops",
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
					Value: 2,
				},
			},
		},
		"PCAPPayload_os_drops": {
			data: []byte(`
{
	"policy_pcap": {
		"pcap": {
			"os_drops": 2
    	}
	}
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "pcap_os_drops",
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
					Value: 2,
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

func TestDNSConversion(t *testing.T) {
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
		"DNSPayloadWirePacketsIpv4": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
            "wire_packets": {
				"ipv4": 1
			}
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "dns_wire_packets_ipv4",
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
					Value: 1,
				},
			},
		},
		"DNSPayloadXactInQuantiles": {
			data: []byte(`
{
	"policy_dns": {
		"dns": {
        	"xact": {
            	"in": {
					"quantiles_us": {
						"p90": 4
					}
				}
        	}
    	}
	}
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "dns_xact_in_quantiles_us",
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
						Name:  "quantile",
						Value: "0.9",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 4,
				},
			},
		},
		"DNSPayloadCardinalityTotal": {
			data: []byte(`
{
	"policy_dns": {
		"dns": {
        	"cardinality": {
				"qname": 4
        	}
    	}
	}
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "dns_cardinality_qname",
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
					Value: 4,
				},
			},
		},
		"DNSPayloadRatesTotal": {
			data: []byte(`
{
	"policy_dns": {
		"dns": {
        	"rates": {
				"total": {
					"p50": 0,
          			"p90": 0,
          			"p95": 2,
          			"p99": 6
				}
        	}
    	}
	}
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "dns_rates_total",
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
						Name:  "quantile",
						Value: "0.9",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 0,
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
					if len(c.expected.Labels) < 7 {
						receivedLabel = value.Labels
						receivedDatapoint = value.Datapoint
					} else {
						if c.expected.Labels[6].Value == value.Labels[6].Value{
							receivedLabel = value.Labels
							receivedDatapoint = value.Datapoint
						}
					}
				}
			}
			assert.True(t, reflect.DeepEqual(c.expected.Labels, receivedLabel), fmt.Sprintf("%s: expected %v got %v", desc, c.expected.Labels, receivedLabel))
			assert.Equal(t, c.expected.Datapoint.Value, receivedDatapoint.Value, fmt.Sprintf("%s: expected value %f got %f", desc, c.expected.Datapoint.Value, receivedDatapoint.Value))
		})
	}

}

func TestDNSRatesConversion(t *testing.T) {
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
		data               []byte
		expectedLabels     []prometheus.Label
		expectedDatapoints []float64
	}{
		"DNSPayloadRatesTotal": {
			data: []byte(`
{
	"policy_dns": {
		"dns": {
        	"rates": {
				"total": {
					"p50": 0,
          			"p90": 1,
          			"p95": 2,
          			"p99": 6
				}
        	}
    	}
	}
}`),
			expectedLabels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "dns_rates_total",
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
						Name:  "quantile",
						Value: "0.5",
					},
					{
						Name:  "__name__",
						Value: "dns_rates_total",
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
						Name:  "quantile",
						Value: "0.9",
					},
					{
						Name:  "__name__",
						Value: "dns_rates_total",
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
						Name:  "quantile",
						Value: "0.95",
					},
					{
						Name:  "__name__",
						Value: "dns_rates_total",
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
						Name:  "quantile",
						Value: "0.99",
					},
				},
				expectedDatapoints: []float64{0, 1, 2, 6},
			},
		}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			data.Data = c.data
			res, err := be.ProcessMetrics(agent, agentID.String(), data)
			require.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s", desc, err))
			var receivedLabel []prometheus.Label
			var receivedDatapoint []float64

			for _, value := range res {
				if c.expectedLabels[0] == value.Labels[0] {
					for _, labels := range value.Labels {
						receivedLabel = append(receivedLabel, labels)
					}
					receivedDatapoint = append(receivedDatapoint, value.Datapoint.Value)
				}
			}
			assert.Equal(t, c.expectedLabels, receivedLabel, fmt.Sprintf("%s: expected %v got %v", desc, c.expectedLabels, receivedLabel))
			assert.Equal(t, c.expectedDatapoints, receivedDatapoint, fmt.Sprintf("%s: expected %v got %v", desc, c.expectedDatapoints, receivedDatapoint))
		})
	}

}
