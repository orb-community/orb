/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

const PktvisorVersion = "4.2.0"

// NameCount represents the count of a unique domain name
type NameCount struct {
	Name     string `mapstructure:"name"`
	Estimate int64  `mapstructure:"estimate"`
}

// Rates represents a histogram of rates at various percentiles
type Rates struct {
	P50 int64 `mapstructure:"p50"`
	P90 int64 `mapstructure:"p90"`
	P95 int64 `mapstructure:"p95"`
	P99 int64 `mapstructure:"p99"`
}

// Quantiles represents a histogram of various percentiles
type Quantiles struct {
	P50 int64 `mapstructure:"p50"`
	P90 int64 `mapstructure:"p90"`
	P95 int64 `mapstructure:"p95"`
	P99 int64 `mapstructure:"p99"`
}

// DHCPPayload contains the information specifically for the DNS protocol
type DHCPPayload struct {
	WirePackets struct {
		Filtered    int64 `mapstructure:"filtered"`
		Total       int64 `mapstructure:"total"`
		DeepSamples int64 `mapstructure:"deep_samples"`
		Discover    int64 `mapstructure:"discover"`
		Offer       int64 `mapstructure:"offer"`
		Request     int64 `mapstructure:"request"`
		Ack         int64 `mapstructure:"ack"`
		Events      int64 `mapstructure:"events"`
	} `mapstructure:"wire_packets"`
	Rates struct {
		Total  Rates `mapstructure:"total"`
		Events Rates `mapstructure:"events"`
	} `mapstructure:"rates"`
	Period PeriodPayload `mapstructure:"period"`
}

// DNSPayload contains the information specifically for the DNS protocol
type DNSPayload struct {
	WirePackets struct {
		Ipv4        int64 `mapstructure:"ipv4"`
		Ipv6        int64 `mapstructure:"ipv6"`
		Queries     int64 `mapstructure:"queries"`
		Replies     int64 `mapstructure:"replies"`
		TCP         int64 `mapstructure:"tcp"`
		Total       int64 `mapstructure:"total"`
		UDP         int64 `mapstructure:"udp"`
		Nodata      int64 `mapstructure:"nodata"`
		Noerror     int64 `mapstructure:"noerror"`
		Nxdomain    int64 `mapstructure:"nxdomain"`
		Srvfail     int64 `mapstructure:"srvfail"`
		Refused     int64 `mapstructure:"refused"`
		Filtered    int64 `mapstructure:"filtered"`
		DeepSamples int64 `mapstructure:"deep_samples"`
		QueryECS    int64 `mapstructure:"query_ecs"`
		Events      int64 `mapstructure:"events"`
	} `mapstructure:"wire_packets"`
	Rates struct {
		Total  Rates `mapstructure:"total"`
		Events Rates `mapstructure:"events"`
	} `mapstructure:"rates"`
	Cardinality struct {
		Qname int64 `mapstructure:"qname"`
	} `mapstructure:"cardinality"`
	Xact struct {
		Counts struct {
			Total    int64 `mapstructure:"total"`
			TimedOut int64 `mapstructure:"timed_out"`
		} `mapstructure:"counts"`
		In struct {
			QuantilesUS Quantiles   `mapstructure:"quantiles_us"`
			TopSlow     []NameCount `mapstructure:"top_slow"`
			Total       int64       `mapstructure:"total"`
		} `mapstructure:"in"`
		Out struct {
			QuantilesUS Quantiles   `mapstructure:"quantiles_us"`
			TopSlow     []NameCount `mapstructure:"top_slow"`
			Total       int64       `mapstructure:"total"`
		} `mapstructure:"out"`
		Ratio struct {
			Quantiles struct {
				P50 float64 `mapstructure:"p50"`
				P90 float64 `mapstructure:"p90"`
				P95 float64 `mapstructure:"p95"`
				P99 float64 `mapstructure:"p99"`
			} `mapstructure:"quantiles"`
		} `mapstructure:"ratio"`
	} `mapstructure:"xact"`
	TopGeoLocECS        []NameCount   `mapstructure:"top_geoLoc_ecs"`
	TopAsnECS           []NameCount   `mapstructure:"top_asn_ecs"`
	TopQueryECS         []NameCount   `mapstructure:"top_query_ecs"`
	TopQname2           []NameCount   `mapstructure:"top_qname2"`
	TopQname3           []NameCount   `mapstructure:"top_qname3"`
	TopNxdomain         []NameCount   `mapstructure:"top_nxdomain"`
	TopQtype            []NameCount   `mapstructure:"top_qtype"`
	TopRcode            []NameCount   `mapstructure:"top_rcode"`
	TopREFUSED          []NameCount   `mapstructure:"top_refused"`
	TopQnameByRespBytes []NameCount   `mapstructure:"top_qname_by_resp_bytes"`
	TopSRVFAIL          []NameCount   `mapstructure:"top_srvfail"`
	TopNODATA           []NameCount   `mapstructure:"top_nodata"`
	TopUDPPorts         []NameCount   `mapstructure:"top_udp_ports"`
	Period              PeriodPayload `mapstructure:"period"`
}

// PacketPayload contains information about raw packets regardless of protocol
type PacketPayload struct {
	Cardinality struct {
		DstIpsOut int64 `mapstructure:"dst_ips_out"`
		SrcIpsIn  int64 `mapstructure:"src_ips_in"`
	} `mapstructure:"cardinality"`
	Ipv4        int64 `mapstructure:"ipv4"`
	Ipv6        int64 `mapstructure:"ipv6"`
	TCP         int64 `mapstructure:"tcp"`
	Total       int64 `mapstructure:"total"`
	UDP         int64 `mapstructure:"udp"`
	In          int64 `mapstructure:"in"`
	Out         int64 `mapstructure:"out"`
	UnknownDir  int64 `mapstructure:"unknown_dir"`
	OtherL4     int64 `mapstructure:"other_l4"`
	DeepSamples int64 `mapstructure:"deep_samples"`
	Filtered    int64 `mapstructure:"filtered"`
	Events      int64 `mapstructure:"events"`
	Protocol    struct {
		Tcp struct {
			SYN int64 `mapstructure:"syn"`
		} `mapstructure:"tcp"`
	} `mapstructure:"protocol"`
	PayloadSize Quantiles `mapstructure:"payload_size"`
	Rates       struct {
		BytesIn    Rates `mapstructure:"bytes_in"`
		BytesOut   Rates `mapstructure:"bytes_out"`
		BytesTotal Rates `mapstructure:"bytes_total"`
		PpsIn      Rates `mapstructure:"pps_in"`
		PpsOut     Rates `mapstructure:"pps_out"`
		PpsTotal   Rates `mapstructure:"pps_total"`
		PpsEvents  Rates `mapstructure:"pps_events"`
	} `mapstructure:"rates"`
	TopIpv4   []NameCount   `mapstructure:"top_ipv4"`
	TopIpv6   []NameCount   `mapstructure:"top_ipv6"`
	TopGeoLoc []NameCount   `mapstructure:"top_geoLoc"`
	TopASN    []NameCount   `mapstructure:"top_asn"`
	Period    PeriodPayload `mapstructure:"period"`
}

// PcapPayload contains information about pcap input stream
type PcapPayload struct {
	TcpReassemblyErrors int64 `mapstructure:"tcp_reassembly_errors"`
	IfDrops             int64 `mapstructure:"if_drops"`
	OsDrops             int64 `mapstructure:"os_drops"`
}

// PeriodPayload indicates the period of time for which a snapshot refers to
type PeriodPayload struct {
	StartTS int64 `mapstructure:"start_ts"`
	Length  int64 `mapstructure:"length"`
}

// FlowPayload contains the information specifically for the Flow protocol
type FlowPayload struct {
	Devices map[string]struct {
		RecordsFiltered         int64       `mapstructure:"records_filtered"`
		RecordsFlows            int64       `mapstructure:"records_flows"`
		TopInInterfacesBytes    []NameCount `mapstructure:"top_in_interfaces_bytes"`
		TopInInterfacesPackets  []NameCount `mapstructure:"top_in_interfaces_packets"`
		TopOutInterfacesBytes   []NameCount `mapstructure:"top_out_interfaces_bytes"`
		TopOutInterfacesPackets []NameCount `mapstructure:"top_out_interfaces_packets"`
		Interfaces              map[string]struct {
			Cardinality struct {
				Conversations int64 `mapstructure:"conversations"`
				DstIpsOut     int64 `mapstructure:"dst_ips_out"`
				DstPortsOut   int64 `mapstructure:"dst_ports_out"`
				SrcIpsIn      int64 `mapstructure:"src_ips_in"`
				SrcPortsIn    int64 `mapstructure:"src_ports_in"`
			} `mapstructure:"cardinality"`
			InIpv4Bytes                int64       `mapstructure:"in_ipv4_bytes"`
			InIpv4Packets              int64       `mapstructure:"in_ipv4_packets"`
			InIpv6Bytes                int64       `mapstructure:"in_ipv6_bytes"`
			InIpv6Packets              int64       `mapstructure:"in_ipv6_packets"`
			InOtherL4Bytes             int64       `mapstructure:"in_other_l4_bytes"`
			InOtherL4Packets           int64       `mapstructure:"in_other_l4_packets"`
			InTcpBytes                 int64       `mapstructure:"in_tcp_bytes"`
			InTcpPackets               int64       `mapstructure:"in_tcp_packets"`
			InUdpBytes                 int64       `mapstructure:"in_udp_bytes"`
			InUdpPackets               int64       `mapstructure:"in_udp_packets"`
			InBytes                    int64       `mapstructure:"in_bytes"`
			InPackets                  int64       `mapstructure:"in_packets"`
			OutIpv4Bytes               int64       `mapstructure:"out_ipv4_bytes"`
			OutIpv4Packets             int64       `mapstructure:"out_ipv4_packets"`
			OutIpv6Bytes               int64       `mapstructure:"out_ipv6_bytes"`
			OutIpv6Packets             int64       `mapstructure:"out_ipv6_packets"`
			OutOtherL4Bytes            int64       `mapstructure:"out_other_l4_bytes"`
			OutOtherL4Packets          int64       `mapstructure:"out_other_l4_packets"`
			OutTcpBytes                int64       `mapstructure:"out_tcp_bytes"`
			OutTcpPackets              int64       `mapstructure:"out_tcp_packets"`
			OutUdpBytes                int64       `mapstructure:"out_udp_bytes"`
			OutUdpPackets              int64       `mapstructure:"out_udp_packets"`
			OutBytes                   int64       `mapstructure:"out_bytes"`
			OutPackets                 int64       `mapstructure:"out_packets"`
			TopInSrcIpsBytes           []NameCount `mapstructure:"top_in_src_ips_bytes"`
			TopInSrcIpsPackets         []NameCount `mapstructure:"top_in_src_ips_packets"`
			TopInSrcPortsBytes         []NameCount `mapstructure:"top_in_src_ports_bytes"`
			TopInSrcPortsPackets       []NameCount `mapstructure:"top_in_src_ports_packets"`
			TopInSrcIpsAndPortBytes    []NameCount `mapstructure:"top_in_src_ips_and_port_bytes"`
			TopInSrcIpsAndPortPackets  []NameCount `mapstructure:"top_in_src_ips_and_port_packets"`
			TopInDstIpsBytes           []NameCount `mapstructure:"top_in_dst_ips_bytes"`
			TopInDstIpsPackets         []NameCount `mapstructure:"top_in_dst_ips_packets"`
			TopInDstPortsBytes         []NameCount `mapstructure:"top_in_dst_ports_bytes"`
			TopInDstPortsPackets       []NameCount `mapstructure:"top_in_dst_ports_packets"`
			TopInDstIpsAndPortBytes    []NameCount `mapstructure:"top_in_dst_ips_and_port_bytes"`
			TopInDstIpsAndPortPackets  []NameCount `mapstructure:"top_in_dst_ips_and_port_packets"`
			TopOutSrcIpsBytes          []NameCount `mapstructure:"top_out_src_ips_bytes"`
			TopOutSrcIpsPackets        []NameCount `mapstructure:"top_out_src_ips_packets"`
			TopOutSrcPortsBytes        []NameCount `mapstructure:"top_out_src_ports_bytes"`
			TopOutSrcPortsPackets      []NameCount `mapstructure:"top_out_src_ports_packets"`
			TopOutSrcIpsAndPortBytes   []NameCount `mapstructure:"top_out_src_ips_and_port_bytes"`
			TopOutSrcIpsAndPortPackets []NameCount `mapstructure:"top_out_src_ips_and_port_packets"`
			TopOutDstIpsBytes          []NameCount `mapstructure:"top_out_dst_ips_bytes"`
			TopOutDstIpsPackets        []NameCount `mapstructure:"top_out_dst_ips_packets"`
			TopOutDstPortsBytes        []NameCount `mapstructure:"top_out_dst_ports_bytes"`
			TopOutDstPortsPackets      []NameCount `mapstructure:"top_out_dst_ports_packets"`
			TopOutDstIpsAndPortBytes   []NameCount `mapstructure:"top_out_dst_ips_and_port_bytes"`
			TopOutDstIpsAndPortPackets []NameCount `mapstructure:"top_out_dst_ips_and_port_packets"`
			TopConversationsBytes      []NameCount `mapstructure:"top_conversations_bytes"`
			TopConversationsPackets    []NameCount `mapstructure:"top_conversations_packets"`
			TopGeoLocBytes             []NameCount `mapstructure:"top_geoLoc_bytes"`
			TopGeoLocPackets           []NameCount `mapstructure:"top_geoLoc_packets"`
			TopAsnBytes                []NameCount `mapstructure:"top_ASN_bytes"`
			TopAsnPackets              []NameCount `mapstructure:"top_ASN_packets"`
		} `mapstructure:"interfaces"`
	} `mapstructure:"devices"`
	Period PeriodPayload `mapstructure:"period"`
}

// StatSnapshot is a snapshot of a given period from pktvisord
type StatSnapshot struct {
	DNS     *DNSPayload    `mapstructure:"DNS,omitempty"`
	DHCP    *DHCPPayload   `mapstructure:"DHCP,omitempty"`
	Packets *PacketPayload `mapstructure:"Packets,omitempty"`
	Pcap    *PcapPayload   `mapstructure:"Pcap,omitempty"`
	Flow    *FlowPayload   `mapstructure:"Flow,omitempty"`
}
