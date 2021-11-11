/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"context"
	"encoding/json"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/mainflux/mainflux"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/backend"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/opentracing/opentracing-go"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

const (
	inputsJson   = `{"pcap":{"1.0":{"filter":{"bpf":{"type":"string","input":"text","label":"Filter Expression","description":"tcpdump compatible filter expression for limiting the traffic examined (with BPF). See https:\/\/www.tcpdump.org\/manpages\/tcpdump.1.html","props":{"example":"udp port 53 and host 127.0.0.1"}}},"config":{"iface":{"type":"string","input":"text","label":"Network Interface","description":"The network interface to capture traffic from","props":{"required":true,"example":"eth0"}},"host_spec":{"type":"string","input":"text","label":"Host Specification","description":"Subnets (comma separated) which should be considered belonging to this host, in CIDR form. Used for ingress\/egress determination, defaults to host attached to the network interface.","props":{"advanced":true,"example":"10.0.1.0\/24,10.0.2.1\/32,2001:db8::\/64"}},"pcap_source":{"type":"string","input":"select","label":"Packet Capture Engine","description":"Packet capture engine to use. Defaults to best for platform.","props":{"advanced":true,"example":"libpcap","options":{"libpcap":"libpcap","af_packet (linux only)":"af_packet"}}}}}}}`
	handlersJson = `{"dns":{"1.0":{"filter":{"exclude_noerror":{"label":"Exclude NOERROR","type":"bool","input":"checkbox","description":"Filter out all NOERROR responses"},"only_rcode":{"label":"Include Only RCODE","type":"number","input":"select","description":"Filter out any queries which are not the given RCODE","props":{"allow_custom_options":true,"options":{"NOERROR":0,"SERVFAIL":2,"NXDOMAIN":3,"REFUSED":5}}},"only_qname_suffix":{"label":"Include Only QName With Suffix","type":"string[]","input":"text","description":"Filter out any queries whose QName does not end in a suffix on the list","props":{"example":".foo.com,.example.com"}}},"config":{},"metrics":{},"metric_groups":{"cardinality":{"label":"Cardinality","description":"Metrics counting the unique number of items in the stream","metrics":[]},"dns_transactions":{"label":"DNS Transactions (Query\/Reply pairs)","description":"Metrics based on tracking queries and their associated replies","metrics":[]},"top_dns_wire":{"label":"Top N Metrics (Various)","description":"Top N metrics across various details from the DNS wire packets","metrics":[]},"top_qnames":{"label":"Top N QNames (All)","description":"Top QNames across all DNS queries in stream","metrics":[]},"top_qnames_by_rcode":{"label":"Top N QNames (Failing RCodes) ","description":"Top QNames across failing result codes","metrics":[]}}}},"net":{"1.0":{"filter":{},"config":{},"metrics":{},"metric_groups":{"ip_cardinality":{"label":"IP Address Cardinality","description":"Unique IP addresses seen in the stream","metrics":[]},"top_geo":{"label":"Top Geo","description":"Top Geo IP and ASN in the stream","metrics":[]},"top_ips":{"label":"Top IPs","description":"Top IP addresses in the stream","metrics":[]}}}},"dhcp":{"1.0":{"filter":{},"config":{},"metrics":{},"metric_groups":{}}}}`
)

type pktvisorBackend struct {
	auth        mainflux.AuthServiceClient
	agentRepo   fleet.AgentRepository
	Backend     string
	Description string
}

type BackendTaps struct {
	Name             string
	InputType        string
	ConfigPredefined []string
	TotalAgents      uint64
}

func (p pktvisorBackend) Metadata() interface{} {
	return struct {
		Backend     string `json:"backend"`
		Description string `json:"description"`
	}{
		Backend:     p.Backend,
		Description: p.Description,
	}
}

func (p pktvisorBackend) MakeHandler(tracer opentracing.Tracer, opts []kithttp.ServerOption, r *bone.Mux) {
	MakePktvisorHandler(tracer, p, opts, r)
}

func (p pktvisorBackend) handlers() (_ types.Metadata, err error) {
	var handlers types.Metadata
	err = json.Unmarshal([]byte(handlersJson), &handlers)
	if err != nil {
		return nil, err
	}
	return handlers, nil
}

func (p pktvisorBackend) inputs() (_ types.Metadata, err error) {
	var handlers types.Metadata
	err = json.Unmarshal([]byte(inputsJson), &handlers)
	if err != nil {
		return nil, err
	}
	return handlers, nil
}

func (p pktvisorBackend) taps(ctx context.Context, ownerID string) ([]types.Metadata, error) {

	taps, err := p.agentRepo.RetrieveAgentMetadataByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	return taps, nil
}

func Register(auth mainflux.AuthServiceClient, agentRepo fleet.AgentRepository) bool {
	backend.Register("pktvisor", &pktvisorBackend{
		Backend:     "pktvisor",
		Description: "pktvisor observability agent from pktvisor.dev",
		auth:        auth,
		agentRepo:   agentRepo,
	})
	return true
}
