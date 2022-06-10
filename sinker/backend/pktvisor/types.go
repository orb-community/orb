/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

const PktvisorVersion = "3.3.0"

// NameStringCount represents the count of a unique string domain name
type NameStringCount struct {
	Name     string `mapstructure:"name"`
	Estimate int64  `mapstructure:"estimate"`
}

// NameIntCount represents the count of a unique int domain name
type NameIntCount struct {
	Name     int64 `mapstructure:"name"`
	Estimate int64 `mapstructure:"estimate"`
}

// Rates represents a histogram of rates at various percentiles
type Rates struct {
	Live int64 `mapstructure:"live"`
	P50  int64 `mapstructure:"p50"`
	P90  int64 `mapstructure:"p90"`
	P95  int64 `mapstructure:"p95"`
	P99  int64 `mapstructure:"p99"`
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
	} `mapstructure:"wire_packets"`
	Rates struct {
		Total Rates `mapstructure:"total"`
	} `mapstructure:"rates"`
	Period PeriodPayload `mapstructure:"period"`
}

// DNSPayload contains the information specifically for the DNS protocol
type DNSPayload struct {
	WirePackets struct {
		Ipv4     int64 `mapstructure:"ipv4"`
		Ipv6     int64 `mapstructure:"ipv6"`
		Queries  int64 `mapstructure:"queries"`
		Replies  int64 `mapstructure:"replies"`
		TCP      int64 `mapstructure:"tcp"`
		Total    int64 `mapstructure:"total"`
		UDP      int64 `mapstructure:"udp"`
		Noerror  int64 `mapstructure:"noerror"`
		Nxdomain int64 `mapstructure:"nxdomain"`
		Srvfail  int64 `mapstructure:"srvfail"`
		Refused  int64 `mapstructure:"refused"`
	} `mapstructure:"wire_packets"`
	Rates struct {
		Total Rates `mapstructure:"total"`
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
			QuantilesUS struct {
				P50 int64 `mapstructure:"p50"`
				P90 int64 `mapstructure:"p90"`
				P95 int64 `mapstructure:"p95"`
				P99 int64 `mapstructure:"p99"`
			} `mapstructure:"quantiles_us"`
			TopSlow []NameStringCount `mapstructure:"top_slow"`
			Total   int64             `mapstructure:"total"`
		} `mapstructure:"in"`
		Out struct {
			QuantilesUS struct {
				P50 int64 `mapstructure:"p50"`
				P90 int64 `mapstructure:"p90"`
				P95 int64 `mapstructure:"p95"`
				P99 int64 `mapstructure:"p99"`
			} `mapstructure:"quantiles_us"`
			TopSlow []NameStringCount `mapstructure:"top_slow"`
			Total   int64             `mapstructure:"total"`
		} `mapstructure:"out"`
	} `mapstructure:"xact"`
	TopQname2   []NameStringCount `mapstructure:"top_qname2"`
	TopQname3   []NameStringCount `mapstructure:"top_qname3"`
	TopNxdomain []NameStringCount `mapstructure:"top_nxdomain"`
	TopQtype    []NameStringCount `mapstructure:"top_qtype"`
	TopRcode    []NameStringCount `mapstructure:"top_rcode"`
	TopREFUSED  []NameStringCount `mapstructure:"top_refused"`
	TopSRVFAIL  []NameStringCount `mapstructure:"top_srvfail"`
	TopUDPPorts []NameStringCount `mapstructure:"top_udp_ports"`
	Period      PeriodPayload     `mapstructure:"period"`
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
	OtherL4     int64 `mapstructure:"other_l4"`
	DeepSamples int64 `mapstructure:"deep_samples"`
	Rates       struct {
		Pps_in struct {
			Live int64 `mapstructure:"live"`
			P50  int64 `mapstructure:"p50"`
			P90  int64 `mapstructure:"p90"`
			P95  int64 `mapstructure:"p95"`
			P99  int64 `mapstructure:"p99"`
		} `mapstructure:"pps_in"`
		Pps_out struct {
			Live int64 `mapstructure:"live"`
			P50  int64 `mapstructure:"p50"`
			P90  int64 `mapstructure:"p90"`
			P95  int64 `mapstructure:"p95"`
			P99  int64 `mapstructure:"p99"`
		} `mapstructure:"pps_out"`
		Pps_total struct {
			Live int64 `mapstructure:"live"`
			P50  int64 `mapstructure:"p50"`
			P90  int64 `mapstructure:"p90"`
			P95  int64 `mapstructure:"p95"`
			P99  int64 `mapstructure:"p99"`
		} `mapstructure:"pps_total"`
	} `mapstructure:"rates"`
	TopIpv4   []NameStringCount `mapstructure:"top_ipv4"`
	TopIpv6   []NameStringCount `mapstructure:"top_ipv6"`
	TopGeoLoc []NameStringCount `mapstructure:"top_geoLoc"`
	TopASN    []NameStringCount `mapstructure:"top_asn"`
	Period    PeriodPayload     `mapstructure:"period"`
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
	Cardinality struct {
		DstIpsOut   int64 `mapstructure:"dst_ips_out"`
		DstPortsOut int64 `mapstructure:"dst_ports_out"`
		SrcIpsIn    int64 `mapstructure:"src_ips_in"`
		SrcPortsIn  int64 `mapstructure:"src_ports_in"`
	} `mapstructure:"cardinality"`
	DeepSamples int64 `mapstructure:"deep_samples"`
	EventRate   Rates `mapstructure:"event_rate"`
	Filtered    int64 `mapstructure:"filtered"`
	Flows       int64 `mapstructure:"flows"`
	Ipv4        int64 `jmapstructureson:"ipv4"`
	Ipv6        int64 `mapstructure:"ipv6"`
	OtherL4     int64 `mapstructure:"other_l4"`
	PayloadSize struct {
		P50 int64 `mapstructure:"p50"`
		P90 int64 `mapstructure:"p90"`
		P95 int64 `mapstructure:"p95"`
		P99 int64 `mapstructure:"p99"`
	} `mapstructure:"payload_size"`
	Period PeriodPayload `mapstructure:"period"`
	Rates  struct {
		Bps Rates `mapstructure:"bps"`
		Pps Rates `mapstructure:"pps"`
	} `mapstructure:"rates"`
	TCP                  int64             `mapstructure:"tcp"`
	TopDstIpsBytes       []NameStringCount `mapstructure:"top_dst_ips_bytes"`
	TopDstIpsPackets     []NameStringCount `mapstructure:"top_dst_ips_packets"`
	TopDstPortsBytes     []NameIntCount    `mapstructure:"top_dst_ports_bytes"`
	TopDstPortsPackets   []NameIntCount    `mapstructure:"top_dst_ports_packets"`
	TopInIfIndexBytes    []NameIntCount    `mapstructure:"top_in_if_index_bytes"`
	TopInIfIndexPackets  []NameIntCount    `mapstructure:"top_in_if_index_packets"`
	TopOutIfIndexBytes   []NameIntCount    `mapstructure:"top_out_if_index_bytes"`
	TopOutIfIndexPackets []NameIntCount    `mapstructure:"top_out_if_index_packets"`
	TopSrcIpsBytes       []NameStringCount `mapstructure:"top_src_ips_bytes"`
	TopSrcIpsPackets     []NameStringCount `mapstructure:"top_src_ips_packets"`
	TopSrcPortsBytes     []NameIntCount    `mapstructure:"top_src_ports_bytes"`
	TopSrcPortsPackets   []NameIntCount    `mapstructure:"top_src_ports_packets"`
	Total                int64             `mapstructure:"total"`
	Udp                  int64             `mapstructure:"udp"`
}

// StatSnapshot is a snapshot of a given period from pktvisord
type StatSnapshot struct {
	DNS     DNSPayload
	DHCP    DHCPPayload
	Packets PacketPayload
	Pcap    PcapPayload
	Flow    FlowPayload
}
