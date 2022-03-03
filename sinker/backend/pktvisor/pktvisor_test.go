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

func TestDNSTopKMetricsConversion(t *testing.T) {
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

func prependLabel(labelList []prometheus.Label, label prometheus.Label) []prometheus.Label {
	labelList = append(labelList, prometheus.Label{})
	copy(labelList[1:], labelList)
	labelList[0] = label
	return labelList
}