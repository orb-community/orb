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
	Live int64 `mapstructure:"live"`
	P50  int64 `mapstructure:"p50"`
	P90  int64 `mapstructure:"p90"`
	P95  int64 `mapstructure:"p95"`
	P99  int64 `mapstructure:"p99"`
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
		Nodata   int64 `mapstructure:"nodata"`
		Noerror  int64 `mapstructure:"noerror"`
		Nxdomain int64 `mapstructure:"nxdomain"`
		Srvfail  int64 `mapstructure:"srvfail"`
		Refused  int64 `mapstructure:"refused"`
		Filtered int64 `mapstructure:"filtered"`
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
	TopGeoLocECS      []NameCount   `mapstructure:"top_geoLoc_ecs"`
	TopASNECS         []NameCount   `mapstructure:"top_asn_ecs"`
	TopQueryECS       []NameCount   `mapstructure:"top_query_ecs"`
	TopQname2         []NameCount   `mapstructure:"top_qname2"`
	TopQname3         []NameCount   `mapstructure:"top_qname3"`
	TopNxdomain       []NameCount   `mapstructure:"top_nxdomain"`
	TopQtype          []NameCount   `mapstructure:"top_qtype"`
	TopRcode          []NameCount   `mapstructure:"top_rcode"`
	TopREFUSED        []NameCount   `mapstructure:"top_refused"`
	TopSizedQnameResp []NameCount   `mapstructure:"top_qname_by_resp_size"`
	TopSRVFAIL        []NameCount   `mapstructure:"top_srvfail"`
	TopNODATA         []NameCount   `mapstructure:"top_nodata"`
	TopUDPPorts       []NameCount   `mapstructure:"top_udp_ports"`
	Period            PeriodPayload `mapstructure:"period"`
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
	Filtered    int64 `mapstructure:"filtered"`
	Protocol    struct {
		TCP struct {
			SYN int64 `mapstructure:"syn"`
		} `mapstructure:"tcp"`
	} `mapstructure:"protocol"`
	PayloadSize Quantiles `mapstructure:"payload_size"`
	Rates       struct {
		Bps_in    Rates `mapstructure:"bps_in"`
		Bps_out   Rates `mapstructure:"bps_out"`
		Pps_in    Rates `mapstructure:"pps_in"`
		Pps_out   Rates `mapstructure:"pps_out"`
		Pps_total Rates `mapstructure:"pps_total"`
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
	Cardinality struct {
		DstIpsOut   int64 `mapstructure:"dst_ips_out"`
		DstPortsOut int64 `mapstructure:"dst_ports_out"`
		SrcIpsIn    int64 `mapstructure:"src_ips_in"`
		SrcPortsIn  int64 `mapstructure:"src_ports_in"`
	} `mapstructure:"cardinality"`
	DeepSamples int64         `mapstructure:"deep_samples"`
	EventRate   Rates         `mapstructure:"event_rate"`
	Filtered    int64         `mapstructure:"filtered"`
	Flows       int64         `mapstructure:"flows"`
	Ipv4        int64         `jmapstructureson:"ipv4"`
	Ipv6        int64         `mapstructure:"ipv6"`
	OtherL4     int64         `mapstructure:"other_l4"`
	PayloadSize Quantiles     `mapstructure:"payload_size"`
	Period      PeriodPayload `mapstructure:"period"`
	Rates       struct {
		Bps Rates `mapstructure:"bps"`
		Pps Rates `mapstructure:"pps"`
	} `mapstructure:"rates"`
	TCP                     int64       `mapstructure:"tcp"`
	TopDstIpsAndPortBytes   []NameCount `mapstructure:"top_dst_ips_and_port_bytes"`
	TopDstIpsAndPortPackets []NameCount `mapstructure:"top_dst_ips_and_port_packets"`
	TopDstIpsBytes          []NameCount `mapstructure:"top_dst_ips_bytes"`
	TopDstIpsPackets        []NameCount `mapstructure:"top_dst_ips_packets"`
	TopDstPortsBytes        []NameCount `mapstructure:"top_dst_ports_bytes"`
	TopDstPortsPackets      []NameCount `mapstructure:"top_dst_ports_packets"`
	TopInIfIndexBytes       []NameCount `mapstructure:"top_in_if_index_bytes"`
	TopInIfIndexPackets     []NameCount `mapstructure:"top_in_if_index_packets"`
	TopOutIfIndexBytes      []NameCount `mapstructure:"top_out_if_index_bytes"`
	TopOutIfIndexPackets    []NameCount `mapstructure:"top_out_if_index_packets"`
	TopSrcIpsAndPortBytes   []NameCount `mapstructure:"top_src_ips_and_port_bytes"`
	TopSrcIpsAndPortPackets []NameCount `mapstructure:"top_src_ips_and_port_packets"`
	TopSrcIpsBytes          []NameCount `mapstructure:"top_src_ips_bytes"`
	TopSrcIpsPackets        []NameCount `mapstructure:"top_src_ips_packets"`
	TopSrcPortsBytes        []NameCount `mapstructure:"top_src_ports_bytes"`
	TopSrcPortsPackets      []NameCount `mapstructure:"top_src_ports_packets"`
	Total                   int64       `mapstructure:"total"`
	Udp                     int64       `mapstructure:"udp"`
}

// StatSnapshot is a snapshot of a given period from pktvisord
type StatSnapshot struct {
	DNS     DNSPayload
	DHCP    DHCPPayload
	Packets PacketPayload
	Pcap    PcapPayload
	Flow    FlowPayload
}
