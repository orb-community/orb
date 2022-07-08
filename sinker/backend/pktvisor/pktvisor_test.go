package pktvisor_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/pb"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinker/backend"
	"github.com/ns1labs/orb/sinker/backend/pktvisor"
	"github.com/ns1labs/orb/sinker/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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

	var agent = &pb.AgentInfoRes{
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

	commonLabels := []prometheus.Label{
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
	}

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
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dhcp_wire_packets_filtered"})),
				Datapoint: prometheus.Datapoint{
					Value: 10,
				},
			},
		},
		"DHCPPayloadWirePacketsTotal": {
			data: []byte(`
{
	"policy_dhcp": {
        "dhcp": {
            "wire_packets": {
                "total": 10
            }
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dhcp_wire_packets_total"})),
				Datapoint: prometheus.Datapoint{
					Value: 10,
				},
			},
		},
		"DHCPPayloadWirePacketsDeepSamples": {
			data: []byte(`
{
	"policy_dhcp": {
        "dhcp": {
            "wire_packets": {
                "deep_samples": 10
            }
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dhcp_wire_packets_deep_samples"})),
				Datapoint: prometheus.Datapoint{
					Value: 10,
				},
			},
		},
		"DHCPPayloadWirePacketsDiscover": {
			data: []byte(`
{
	"policy_dhcp": {
        "dhcp": {
            "wire_packets": {
                "discover": 10
            }
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dhcp_wire_packets_discover"})),
				Datapoint: prometheus.Datapoint{
					Value: 10,
				},
			},
		},
		"DHCPPayloadWirePacketsOffer": {
			data: []byte(`
{
	"policy_dhcp": {
        "dhcp": {
            "wire_packets": {
                "offer": 10
            }
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dhcp_wire_packets_offer"})),
				Datapoint: prometheus.Datapoint{
					Value: 10,
				},
			},
		},
		"DHCPPayloadWirePacketsRequest": {
			data: []byte(`
{
	"policy_dhcp": {
        "dhcp": {
            "wire_packets": {
                "request": 10
            }
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dhcp_wire_packets_request"})),
				Datapoint: prometheus.Datapoint{
					Value: 10,
				},
			},
		},
		"DHCPPayloadWirePacketsAck": {
			data: []byte(`
{
	"policy_dhcp": {
        "dhcp": {
            "wire_packets": {
                "ack": 10
            }
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dhcp_wire_packets_ack"})),
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

	var agent = &pb.AgentInfoRes{
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

	var agent = &pb.AgentInfoRes{
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

	var agent = &pb.AgentInfoRes{
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

	commonLabels := []prometheus.Label{
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
	}

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
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "pcap_tcp_reassembly_errors",
				})),
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
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "pcap_if_drops",
				})),
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
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "pcap_os_drops",
				})),
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

	var agent = &pb.AgentInfoRes{
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

	commonLabels := []prometheus.Label{
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
	}

	cases := map[string]struct {
		data     []byte
		expected prometheus.TimeSeries
	}{
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
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_cardinality_qname",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 4,
				},
			},
		},
		"DNSPayloadTopNxdomain": {
			data: []byte(`
{
	"policy_dns": {
		"dns": {
        	"top_nxdomain": [
				{
	          		"estimate": 186,
         			"name": "89.187.189.231"				
        		}
			]
    	}
	}
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_top_nxdomain",
				}), prometheus.Label{
					Name:  "qname",
					Value: "89.187.189.231",
				}),
				Datapoint: prometheus.Datapoint{
					Value: 186,
				},
			},
		},
		"DNSPayloadTopRefused": {
			data: []byte(`
{
	"policy_dns": {
		"dns": {
        	"top_refused": [
				{
	          		"estimate": 186,
         			"name": "89.187.189.231"				
        		}
			]
    	}
	}
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_top_refused",
				}), prometheus.Label{
					Name:  "qname",
					Value: "89.187.189.231",
				}),
				Datapoint: prometheus.Datapoint{
					Value: 186,
				},
			},
		},
		"DNSPayloadTopSrvfail": {
			data: []byte(`
{
	"policy_dns": {
		"dns": {
        	"top_srvfail": [
				{
	          		"estimate": 186,
         			"name": "89.187.189.231"				
        		}
			]
    	}
	}
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_top_srvfail",
				}), prometheus.Label{
					Name:  "qname",
					Value: "89.187.189.231",
				}),
				Datapoint: prometheus.Datapoint{
					Value: 186,
				},
			},
		},
		"DNSPayloadTopNodata": {
			data: []byte(`
{
	"policy_dns": {
		"dns": {
        	"top_nodata": [
				{
	          		"estimate": 186,
         			"name": "89.187.189.231"				
        		}
			]
    	}
	}
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_top_nodata",
				}), prometheus.Label{
					Name:  "qname",
					Value: "89.187.189.231",
				}),
				Datapoint: prometheus.Datapoint{
					Value: 186,
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
						if c.expected.Labels[6].Value == value.Labels[6].Value {
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

	var agent = &pb.AgentInfoRes{
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

	commonLabels := []prometheus.Label{
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
	}

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
			expectedLabels: labelQuantiles(commonLabels, prometheus.Label{
				Name:  "__name__",
				Value: "dns_rates_total",
			}),
			expectedDatapoints: []float64{0, 1, 2, 6},
		},
		"PacketsPayloadRatesPpsIn": {
			data: []byte(`
{
	"policy_dns": {
		"packets": {
        	"rates": {
				"pps_in": {
					"p50": 0,
          			"p90": 1,
        			"p95": 2,
        			"p99": 6
				}
        	}
    	}
	}
}`),
			expectedLabels: labelQuantiles(commonLabels, prometheus.Label{
				Name:  "__name__",
				Value: "packets_rates_pps_in",
			}),
			expectedDatapoints: []float64{0, 1, 2, 6},
		},
		"PacketsPayloadRatesPpsTotal": {
			data: []byte(`
{
	"policy_dns": {
		"packets": {
        	"rates": {
				"pps_total": {
					"p50": 0,
          			"p90": 1,
        			"p95": 2,
        			"p99": 6
				}
        	}
    	}
	}
}`),
			expectedLabels: labelQuantiles(commonLabels, prometheus.Label{
				Name:  "__name__",
				Value: "packets_rates_pps_total",
			}),
			expectedDatapoints: []float64{0, 1, 2, 6},
		},
		"PacketsPayloadRatesPpsOut": {
			data: []byte(`
{
	"policy_dns": {
		"packets": {
        	"rates": {
				"pps_out": {
					"p50": 0,
          			"p90": 1,
        			"p95": 2,
        			"p99": 6
				}
        	}
    	}
	}
}`),
			expectedLabels: labelQuantiles(commonLabels, prometheus.Label{
				Name:  "__name__",
				Value: "packets_rates_pps_out",
			}),
			expectedDatapoints: []float64{0, 1, 2, 6},
		},
		"DHCPPayloadRates": {
			data: []byte(`
{
	"policy_dhcp": {
		"dhcp": {
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
			expectedLabels: labelQuantiles(commonLabels, prometheus.Label{
				Name:  "__name__",
				Value: "dhcp_rates_total",
			}),
			expectedDatapoints: []float64{0, 1, 2, 6},
		},
		"FlowPayloadRatesBPS": {
			data: []byte(`
				{
					"policy_flow": {
						"flow": {
							"rates": {
							  "bps": {
								"p50": 2,
								"p90": 3,
								"p95": 76811346,
								"p99": 76811347
							  }
							}
						}
					}
				}`),
			expectedLabels: labelQuantiles(commonLabels, prometheus.Label{
				Name:  "__name__",
				Value: "flow_rates_bps",
			}),
			expectedDatapoints: []float64{2, 3, 76811346, 76811347},
		},
		"FlowPayloadRatesPPS": {
			data: []byte(`
				{
					"policy_flow": {
						"flow": {
							"rates": {
							  "pps": {
								"p50": 2,
								"p90": 3,
								"p95": 76811346,
								"p99": 76811347
							  }
							}
						}
					}
				}`),
			expectedLabels: labelQuantiles(commonLabels, prometheus.Label{
				Name:  "__name__",
				Value: "flow_rates_pps",
			}),
			expectedDatapoints: []float64{2, 3, 76811346, 76811347},
		},
		"FlowPayloadSize": {
			data: []byte(`
				{
					"policy_flow": {
						"flow": {
							"payload_size": {
								"p50": 2,
								"p90": 3,
								"p95": 76811346,
								"p99": 76811347
							}
						}
					}
				}`),
			expectedLabels: labelQuantiles(commonLabels, prometheus.Label{
				Name:  "__name__",
				Value: "flow_payload_size",
			}),
			expectedDatapoints: []float64{2, 3, 76811346, 76811347},
		},
		"FlowEventRate": {
			data: []byte(`
				{
					"policy_flow": {
						"flow": {
							"event_rate": {
								"p50": 2,
								"p90": 3,
								"p95": 76811346,
								"p99": 76811347
							}
						}
					}
				}`),
			expectedLabels: labelQuantiles(commonLabels, prometheus.Label{
				Name:  "__name__",
				Value: "flow_event_rate",
			}),
			expectedDatapoints: []float64{2, 3, 76811346, 76811347},
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

			assert.ElementsMatch(t, c.expectedLabels, receivedLabel, fmt.Sprintf("%s: expected %v got %v", desc, c.expectedLabels, receivedLabel))
			assert.ElementsMatch(t, c.expectedDatapoints, receivedDatapoint, fmt.Sprintf("%s: expected %v got %v", desc, c.expectedDatapoints, receivedDatapoint))
		})
	}

}

func TestDNSTopKMetricsConversion(t *testing.T) {
	var logger *zap.Logger
	pktvisor.Register(logger)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	agentID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var agent = &pb.AgentInfoRes{
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
		"PacketPayloadToqQName2": {
			data: []byte(`
{
    "policy_dns": {
		"dns": {
        	"top_qname2": [
				{
          	  	  "estimate": 8,
          	  	  "name": ".google.com"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "dns_top_qname2",
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
						Name:  "qname",
						Value: ".google.com",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"PacketPayloadToqQName3": {
			data: []byte(`
{
    "policy_dns": {
		"dns": {
        	"top_qname3": [
				{
          	  	  "estimate": 6,
          	  	  "name": ".l.google.com"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "dns_top_qname3",
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
						Name:  "qname",
						Value: ".l.google.com",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 6,
				},
			},
		},
		"PacketPayloadTopQueryECS": {
			data: []byte(`
{
    "policy_dns": {
		"dns": {
        	"top_query_ecs": [
				{
          	  	  "estimate": 6,
          	  	  "name": "2001:470:1f0b:1600::"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "dns_top_query_ecs",
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
						Name:  "ecs",
						Value: "2001:470:1f0b:1600::",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 6,
				},
			},
		},
		"PacketPayloadToqQType": {
			data: []byte(`
{
    "policy_dns": {
		"dns": {
        	"top_qtype": [
				{
          	  	  "estimate": 6,
          	  	  "name": "HTTPS"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "dns_top_qtype",
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
						Name:  "qtype",
						Value: "HTTPS",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 6,
				},
			},
		},
		"PacketPayloadTopUDPPorts": {
			data: []byte(`
{
    "policy_dns": {
		"dns": {
      		"top_udp_ports": [
			  {
				"estimate": 2,
          		"name": "39783"
			  }
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "dns_top_udp_ports",
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
						Name:  "port",
						Value: "39783",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 2,
				},
			},
		},
		"PacketPayloadTopRCode": {
			data: []byte(`
{
    "policy_dns": {
		"dns": {
        	"top_rcode": [
				{
          	  	  "estimate": 8,
          	  	  "name": "NOERROR"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "dns_top_rcode",
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
						Name:  "rcode",
						Value: "NOERROR",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
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

func TestDNSWirePacketsConversion(t *testing.T) {
	var logger *zap.Logger
	pktvisor.Register(logger)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	agentID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var agent = &pb.AgentInfoRes{
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

	commonLabels := []prometheus.Label{
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
	}

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
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_wire_packets_ipv4",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 1,
				},
			},
		},
		"DNSPayloadWirePacketsIpv6": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
            "wire_packets": {
				"ipv6": 14
			}
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_wire_packets_ipv6",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 14,
				},
			},
		},
		"DNSPayloadWirePacketsNodata": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
            "wire_packets": {
				"nodata": 8
			}
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_wire_packets_nodata",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"DNSPayloadWirePacketsNoerror": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
            "wire_packets": {
				"noerror": 8
			}
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_wire_packets_noerror",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"DNSPayloadWirePacketsNxdomain": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
            "wire_packets": {
				"nxdomain": 6
			}
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_wire_packets_nxdomain",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 6,
				},
			},
		},
		"DNSPayloadWirePacketsQueries": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
            "wire_packets": {
				"queries": 7
			}
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_wire_packets_queries",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 7,
				},
			},
		},
		"DNSPayloadWirePacketsRefused": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
            "wire_packets": {
				"refused": 8
			}
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_wire_packets_refused",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"DNSPayloadWirePacketsFiltered": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
            "wire_packets": {
				"filtered": 8
			}
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_wire_packets_filtered",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"DNSPayloadWirePacketsReplies": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
            "wire_packets": {
				"replies": 8
			}
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_wire_packets_replies",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"DNSPayloadWirePacketsSrvfail": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
            "wire_packets": {
				"srvfail": 9
			}
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_wire_packets_srvfail",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 9,
				},
			},
		},
		"DNSPayloadWirePacketsTcp": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
            "wire_packets": {
				"tcp": 9
			}
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_wire_packets_tcp",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 9,
				},
			},
		},
		"DNSPayloadWirePacketsTotal": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
            "wire_packets": {
				"total": 9
			}
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_wire_packets_total",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 9,
				},
			},
		},
		"DNSPayloadWirePacketsUdp": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
            "wire_packets": {
				"udp": 9
			}
        }
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_wire_packets_udp",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 9,
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

func TestDNSXactConversion(t *testing.T) {
	var logger *zap.Logger
	pktvisor.Register(logger)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	agentID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var agent = &pb.AgentInfoRes{
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

	commonLabels := []prometheus.Label{
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
	}

	cases := map[string]struct {
		data     []byte
		expected prometheus.TimeSeries
	}{
		"DNSPayloadXactCountsTimedOut": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
			"xact": {
	            "counts": {
					"timed_out": 1
				}
        	}
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_xact_counts_timed_out",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 1,
				},
			},
		},
		"DNSPayloadXactCountsTotal": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
			"xact": {
	            "counts": {
					"total": 8
				}
        	}
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_xact_counts_total",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"DNSPayloadXactInTotal": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
			"xact": {
	            "in": {
					"total": 8
				}
        	}
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_xact_in_total",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"DNSPayloadXactInTopSlow": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
			"xact": {
	            "in": {
					"top_slow": [
						{
							"estimate": 111,
							"name": "23.43.252.68"						
						}
					]
				}
        	}
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_xact_in_top_slow",
				}), prometheus.Label{
					Name:  "qname",
					Value: "23.43.252.68",
				}),
				Datapoint: prometheus.Datapoint{
					Value: 111,
				},
			},
		},
		"DNSPayloadXactOutTopSlow": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
			"xact": {
	            "out": {
					"top_slow": [
						{
							"estimate": 111,
							"name": "23.43.252.68"						
						}
					]
				}
        	}
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_xact_out_top_slow",
				}), prometheus.Label{
					Name:  "qname",
					Value: "23.43.252.68",
				}),
				Datapoint: prometheus.Datapoint{
					Value: 111,
				},
			},
		},
		"DNSPayloadXactOutTotal": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
			"xact": {
	            "out": {
					"total": 8
				}
        	}
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_xact_out_total",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 8,
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

func TestPacketsConversion(t *testing.T) {
	var logger *zap.Logger
	pktvisor.Register(logger)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	agentID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var agent = &pb.AgentInfoRes{
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

	commonLabels := []prometheus.Label{
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
	}

	cases := map[string]struct {
		data     []byte
		expected prometheus.TimeSeries
	}{
		"DNSPayloadPacketsCardinalityDst": {
			data: []byte(`
{
	"policy_dns": {
        "packets": {
			"cardinality": {
	            "dst_ips_out": 41
        	}
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_cardinality_dst_ips_out",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 41,
				},
			},
		},
		"DNSPayloadPacketsCardinalitySrc": {
			data: []byte(`
{
	"policy_dns": {
        "packets": {
			"cardinality": {
	            "src_ips_in": 43
        	}
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_cardinality_src_ips_in",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 43,
				},
			},
		},
		"DNSPayloadPacketsDeepSamples": {
			data: []byte(`
{
	"policy_dns": {
        "packets": {
			"deep_samples": 3139
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_deep_samples",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 3139,
				},
			},
		},
		"DNSPayloadPacketsIn": {
			data: []byte(`
{
	"policy_dns": {
        "packets": {
			"in": 1422
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_in",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 1422,
				},
			},
		},
		"DNSPayloadPacketsIpv4": {
			data: []byte(`
{
	"policy_dns": {
        "packets": {
			"ipv4": 2506
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_ipv4",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 2506,
				},
			},
		},
		"DNSPayloadPacketsIpv6": {
			data: []byte(`
{
	"policy_dns": {
        "packets": {
			"ipv6": 2506
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_ipv6",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 2506,
				},
			},
		},
		"DNSPayloadPacketsOtherL4": {
			data: []byte(`
{
	"policy_dns": {
        "packets": {
			"other_l4": 637
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_other_l4",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 637,
				},
			},
		},
		"DNSPayloadPacketsFiltered": {
			data: []byte(`
{
	"policy_dns": {
        "packets": {
			"filtered": 637
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_filtered",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 637,
				},
			},
		},
		"DNSPayloadPacketsOut": {
			data: []byte(`
{
	"policy_dns": {
        "packets": {
			"out": 1083
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_out",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 1083,
				},
			},
		},
		"DNSPayloadPacketsTcp": {
			data: []byte(`
{
	"policy_dns": {
        "packets": {
			"tcp": 549
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_tcp",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 549,
				},
			},
		},
		"DNSPayloadPacketsTotal": {
			data: []byte(`
{
	"policy_dns": {
        "packets": {
			"total": 3139
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_total",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 3139,
				},
			},
		},
		"DNSPayloadPacketsUdp": {
			data: []byte(`
{
	"policy_dns": {
        "packets": {
			"udp": 1953
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_udp",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 1953,
				},
			},
		},
		"DNSPayloadPacketsTopIpv4": {
			data: []byte(`
{
	"policy_dns": {
        "packets": {
			"top_ipv4": [
				{
					"estimate": 996,
					"name": "103.6.85.201"					
				}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_top_ipv4",
				}), prometheus.Label{
					Name:  "ipv4",
					Value: "103.6.85.201",
				}),
				Datapoint: prometheus.Datapoint{
					Value: 996,
				},
			},
		},
		"DNSPayloadPacketsTopIpv6": {
			data: []byte(`
{
	"policy_dns": {
        "packets": {
			"top_ipv6": [
				{
					"estimate": 996,
					"name": "103.6.85.201"					
				}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_top_ipv6",
				}), prometheus.Label{
					Name:  "ipv6",
					Value: "103.6.85.201",
				}),
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

func TestPeriodConversion(t *testing.T) {
	var logger *zap.Logger
	pktvisor.Register(logger)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	agentID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var agent = &pb.AgentInfoRes{
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

	commonLabels := []prometheus.Label{
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
	}

	cases := map[string]struct {
		data            []byte
		expectedLength  prometheus.TimeSeries
		expectedStartTs prometheus.TimeSeries
	}{
		"DNSPayloadPeriod": {
			data: []byte(`
{
	"policy_dns": {
        "dns": {
			"period": {
		        "length": 60,
        		"start_ts": 1624888107
        	}
		}
    }
}`),
			expectedLength: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_period_length",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 60,
				},
			},
			expectedStartTs: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dns_period_start_ts",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 1624888107,
				},
			},
		},
		"PacketsPayloadPeriod": {
			data: []byte(`
{
	"policy_packets": {
        "packets": {
			"period": {
		        "length": 60,
        		"start_ts": 1624888107
        	}
		}
    }
}`),
			expectedLength: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_period_length",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 60,
				},
			},
			expectedStartTs: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "packets_period_start_ts",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 1624888107,
				},
			},
		},
		"DHCPPayloadPeriod": {
			data: []byte(`
{
	"policy_dhcp": {
        "dhcp": {
			"period": {
		        "length": 60,
        		"start_ts": 1624888107
        	}
		}
    }
}`),
			expectedLength: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dhcp_period_length",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 60,
				},
			},
			expectedStartTs: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "dhcp_period_start_ts",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 1624888107,
				},
			},
		},
		"FlowPayloadPeriod": {
			data: []byte(`
{
	"policy_flow": {
        "flow": {
			"period": {
		        "length": 60,
        		"start_ts": 1624888107
        	}
		}
    }
}`),
			expectedLength: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "flow_period_length",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 60,
				},
			},
			expectedStartTs: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "flow_period_start_ts",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 1624888107,
				},
			},
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			data.Data = c.data
			res, err := be.ProcessMetrics(agent, agentID.String(), data)
			require.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s", desc, err))
			var receivedLabelStartTs []prometheus.Label
			var receivedDatapointStartTs prometheus.Datapoint
			var receivedLabelLength []prometheus.Label
			var receivedDatapointLength prometheus.Datapoint
			for _, value := range res {
				if c.expectedLength.Labels[0] == value.Labels[0] {
					receivedLabelLength = value.Labels
					receivedDatapointLength = value.Datapoint
				} else if c.expectedStartTs.Labels[0] == value.Labels[0] {
					receivedLabelStartTs = value.Labels
					receivedDatapointStartTs = value.Datapoint
				}
			}
			assert.True(t, reflect.DeepEqual(c.expectedLength.Labels, receivedLabelLength), fmt.Sprintf("%s: expected %v got %v", desc, c.expectedLength.Labels, receivedLabelLength))
			assert.Equal(t, c.expectedLength.Datapoint.Value, receivedDatapointLength.Value, fmt.Sprintf("%s: expected value %f got %f", desc, c.expectedLength.Datapoint.Value, receivedDatapointLength.Value))
			assert.True(t, reflect.DeepEqual(c.expectedStartTs.Labels, receivedLabelStartTs), fmt.Sprintf("%s: expected %v got %v", desc, c.expectedStartTs.Labels, receivedLabelStartTs))
			assert.Equal(t, c.expectedStartTs.Datapoint.Value, receivedDatapointStartTs.Value, fmt.Sprintf("%s: expected value %f got %f", desc, c.expectedStartTs.Datapoint.Value, receivedDatapointStartTs.Value))

		})
	}
}

func TestFlowCardinalityConversion(t *testing.T) {
	var logger *zap.Logger
	pktvisor.Register(logger)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	agentID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var agent = &pb.AgentInfoRes{
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

	commonLabels := []prometheus.Label{
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
	}

	cases := map[string]struct {
		data     []byte
		expected prometheus.TimeSeries
	}{
		"FlowPayloadCardinalityDstIpsOut": {
			data: []byte(`
				{
					"policy_flow": {
					  "flow": {
						"cardinality": {
						  "dst_ips_out": 4
						}
					  }
					}
				}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "flow_cardinality_dst_ips_out",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 4,
				},
			},
		},
		"FlowPayloadCardinalityDstPortsOut": {
			data: []byte(`
				{
					"policy_flow": {
					  "flow": {
						"cardinality": {
						  "dst_ports_out": 31,
						  "src_ips_in": 4,
						  "src_ports_in": 31
						}
					  }
					}
				}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "flow_cardinality_dst_ports_out",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 31,
				},
			},
		},
		"FlowPayloadCardinalitySrcIpsIn": {
			data: []byte(`
				{
					"policy_flow": {
					  "flow": {
						"cardinality": {
						  "src_ips_in": 4,
						  "src_ports_in": 31
						}
					  }
					}
				}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "flow_cardinality_src_ips_in",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 4,
				},
			},
		},
		"FlowPayloadCardinalitySrcPortsIn": {
			data: []byte(`
				{
					"policy_flow": {
					  "flow": {
						"cardinality": {
						  "src_ports_in": 31
						}
					  }
					}
				}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "flow_cardinality_src_ports_in",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 31,
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
						if c.expected.Labels[6].Value == value.Labels[6].Value {
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

func TestFlowConversion(t *testing.T) {
	var logger *zap.Logger
	pktvisor.Register(logger)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	agentID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var agent = &pb.AgentInfoRes{
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

	commonLabels := []prometheus.Label{
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
	}

	cases := map[string]struct {
		data     []byte
		expected prometheus.TimeSeries
	}{
		"FlowPayloadDeepSamples": {
			data: []byte(`
				{
					"policy_flow": {
						"flow": {
							"deep_samples": 9279
						}
					}
				}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "flow_deep_samples",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 9279,
				},
			},
		},
		"FlowPayloadFiltered": {
			data: []byte(`
				{
					"policy_flow": {
						"flow": {
							"filtered": 8
						}
					}
				}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "flow_filtered",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowPayloadFlows": {
			data: []byte(`
				{
					"policy_flow": {
						"flow": {
							"flows": 8
						}
					}
				}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "flow_flows",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowPayloadIpv4": {
			data: []byte(`
				{
					"policy_flow": {
						"flow": {
							"ipv4": 52785
						}
					}
				}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "flow_ipv4",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 52785,
				},
			},
		},
		"FlowPayloadIpv6": {
			data: []byte(`
				{
					"policy_flow": {
						"flow": {
							"ipv6": 52785
						}
					}
				}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "flow_ipv6",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 52785,
				},
			},
		},
		"FlowPayloadOtherL4": {
			data: []byte(`
				{
					"policy_flow": {
						"flow": {
							"other_l4": 52785
						}
					}
				}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "flow_other_l4",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 52785,
				},
			},
		},
		"FlowPayloadTCP": {
			data: []byte(`
				{
					"policy_flow": {
						"flow": {
							"tcp": 52785
						}
					}
				}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "flow_tcp",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 52785,
				},
			},
		},
		"FlowPayloadTotal": {
			data: []byte(`
				{
					"policy_flow": {
						"flow": {
							"total": 52785
						}
					}
				}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "flow_total",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 52785,
				},
			},
		},
		"FlowPayloadUdp": {
			data: []byte(`
				{
					"policy_flow": {
						"flow": {
							"udp": 52785
						}
					}
				}`),
			expected: prometheus.TimeSeries{
				Labels: append(prependLabel(commonLabels, prometheus.Label{
					Name:  "__name__",
					Value: "flow_udp",
				})),
				Datapoint: prometheus.Datapoint{
					Value: 52785,
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

func TestFlowTopKMetricsConversion(t *testing.T) {
	var logger *zap.Logger
	pktvisor.Register(logger)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	agentID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var agent = &pb.AgentInfoRes{
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
		"FlowTopDstIpsAndPortBytes": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_dst_ips_and_port_bytes": [
				{
          	  	  "estimate": 8,
          	  	  "name": "10.4.2.2:5000"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_dst_ips_and_port_bytes",
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
						Name:  "ip_port",
						Value: "10.4.2.2:5000",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowTopDstIpsAndPortPackets": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_dst_ips_and_port_packets": [
				{
          	  	  "estimate": 8,
          	  	  "name": "10.4.2.2:5000"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_dst_ips_and_port_packets",
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
						Name:  "ip_port",
						Value: "10.4.2.2:5000",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowTopDstIpsBytes": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_dst_ips_bytes": [
				{
          	  	  "estimate": 8,
          	  	  "name": "10.4.2.2"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_dst_ips_bytes",
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
						Name:  "ip",
						Value: "10.4.2.2",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowTopDstIpsPackets": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_dst_ips_packets": [
				{
          	  	  "estimate": 8,
          	  	  "name": "10.4.2.2"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_dst_ips_packets",
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
						Name:  "ip",
						Value: "10.4.2.2",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowTopDstPortsBytes": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_dst_ports_bytes": [
				{
          	  	  "estimate": 8,
          	  	  "name": "5000"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_dst_ports_bytes",
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
						Name:  "port",
						Value: "5000",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowTopDstPortsPackets": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_dst_ports_packets": [
				{
          	  	  "estimate": 8,
          	  	  "name": "5000"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_dst_ports_packets",
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
						Name:  "port",
						Value: "5000",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowTopInIfIndexBytes": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_in_if_index_bytes": [
				{
          	  	  "estimate": 8,
          	  	  "name": "300"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_in_if_index_bytes",
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
						Name:  "index",
						Value: "300",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowTopInIfIndexPackets": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_in_if_index_packets": [
				{
          	  	  "estimate": 8,
          	  	  "name": "300"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_in_if_index_packets",
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
						Name:  "index",
						Value: "300",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowTopOutIfIndexBytes": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_out_if_index_bytes": [
				{
          	  	  "estimate": 8,
          	  	  "name": "200"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_out_if_index_bytes",
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
						Name:  "index",
						Value: "200",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowTopOutIfIndexPackets": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_out_if_index_packets": [
				{
          	  	  "estimate": 8,
          	  	  "name": "200"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_out_if_index_packets",
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
						Name:  "index",
						Value: "200",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		}, "FlowTopSrcIpsAndPortBytes": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_src_ips_and_port_bytes": [
				{
          	  	  "estimate": 8,
          	  	  "name": "10.4.2.2:5000"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_src_ips_and_port_bytes",
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
						Name:  "ip_port",
						Value: "10.4.2.2:5000",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowTopSrcIpsAndPortPackets": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_src_ips_and_port_packets": [
				{
          	  	  "estimate": 8,
          	  	  "name": "10.4.2.2:5000"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_src_ips_and_port_packets",
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
						Name:  "ip_port",
						Value: "10.4.2.2:5000",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowTopSrcIpsBytes": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_src_ips_bytes": [
				{
          	  	  "estimate": 8,
          	  	  "name": "10.4.2.2"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_src_ips_bytes",
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
						Name:  "ip",
						Value: "10.4.2.2",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowTopSrcIpsPackets": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_src_ips_packets": [
				{
          	  	  "estimate": 8,
          	  	  "name": "10.4.2.2"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_src_ips_packets",
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
						Name:  "ip",
						Value: "10.4.2.2",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowTopSrcPortsBytes": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_src_ports_bytes": [
				{
          	  	  "estimate": 8,
          	  	  "name": "4500"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_src_ports_bytes",
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
						Name:  "port",
						Value: "4500",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
				},
			},
		},
		"FlowTopSrcPortsPackets": {
			data: []byte(`
{
    "policy_flow": {
		"flow": {
        	"top_src_ports_packets": [
				{
          	  	  "estimate": 8,
          	  	  "name": "4500"
        		}
			]
		}
    }
}`),
			expected: prometheus.TimeSeries{
				Labels: []prometheus.Label{
					{
						Name:  "__name__",
						Value: "flow_top_src_ports_packets",
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
						Name:  "port",
						Value: "4500",
					},
				},
				Datapoint: prometheus.Datapoint{
					Value: 8,
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

func TestAgentTagsConversion(t *testing.T) {
	var logger *zap.Logger
	pktvisor.Register(logger)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	agentID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var agent = &pb.AgentInfoRes{
		OwnerID:   ownerID.String(),
		AgentName: "agent-test",
		AgentTags: types.Tags{"testkey": "testvalue", "testkey2": "testvalue2"},
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
		"Example metrics": {
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
						Name:  "testkey",
						Value: "testvalue",
					},
					{
						Name:  "testkey2",
						Value: "testvalue2",
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
			assert.ElementsMatch(t, c.expected.Labels, receivedLabel, fmt.Sprintf("%s: expected %v got %v", desc, c.expected.Labels, receivedLabel))
			assert.Equal(t, c.expected.Datapoint.Value, receivedDatapoint.Value, fmt.Sprintf("%s: expected value %f got %f", desc, c.expected.Datapoint.Value, receivedDatapoint.Value))
		})
	}

}

func TestTagsConversion(t *testing.T) {
	var logger *zap.Logger
	pktvisor.Register(logger)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	agentID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	var agent = &pb.AgentInfoRes{
		OwnerID:   ownerID.String(),
		AgentName: "agent-test",
		AgentTags: types.Tags{"test": "true"},
		OrbTags:   types.Tags{"test2": "true2"},
	}

	var sameTagKeyAgent = &pb.AgentInfoRes{
		OwnerID:   ownerID.String(),
		AgentName: "agent-test",
		AgentTags: types.Tags{"test": "true"},
		OrbTags:   types.Tags{"test": "true2"},
	}

	data := fleet.AgentMetricsRPCPayload{
		PolicyID:   policyID.String(),
		PolicyName: "policy-test",
		Datasets:   nil,
		Format:     "json",
		BEVersion:  "1.0",
		Data: []byte(`
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
	}

	be := backend.GetBackend("pktvisor")

	commonLabels := []prometheus.Label{
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
	}

	cases := map[string]struct {
		agent    *pb.AgentInfoRes
		expected prometheus.TimeSeries
	}{
		"Different agent tags and orb tag": {
			agent: agent,
			expected: prometheus.TimeSeries{
				Labels: append(commonLabels, prometheus.Label{
					Name:  "test",
					Value: "true",
				}, prometheus.Label{
					Name:  "test2",
					Value: "true2",
				}),
			},
		},
		"Same key agent tags and orb tag": {
			agent: sameTagKeyAgent,
			expected: prometheus.TimeSeries{
				Labels: append(commonLabels, prometheus.Label{
					Name:  "test",
					Value: "true2",
				}),
			},
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			res, err := be.ProcessMetrics(c.agent, agentID.String(), data)
			require.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s", desc, err))
			var receivedLabel []prometheus.Label
			for _, value := range res {
				if commonLabels[0].Value == value.Labels[0].Value {
					receivedLabel = value.Labels
				}
			}
			assert.ElementsMatch(t, c.expected.Labels, receivedLabel, fmt.Sprintf("%s: expected %v got %v", desc, c.expected.Labels, receivedLabel))
		})
	}

}

func prependLabel(labelList []prometheus.Label, label prometheus.Label) []prometheus.Label {
	labelList = append(labelList, prometheus.Label{})
	copy(labelList[1:], labelList)
	labelList[0] = label
	return labelList
}

func labelQuantiles(labelList []prometheus.Label, label prometheus.Label) []prometheus.Label {
	for i := 0; i < 28; i += 7 {
		labelList = append(labelList[:i+1], labelList[i:]...)
		labelList[i] = label
	}
	return labelList
}
