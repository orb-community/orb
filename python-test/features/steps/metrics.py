from hamcrest import *
import requests
from prometheus_client.parser import text_string_to_metric_families
from utils import threading_wait_until


@threading_wait_until
def wait_until_metrics_scraped(local_prometheus_endpoint, expected_metrics, event=None):
    """

    :param local_prometheus_endpoint: endpoint /prometheus
    :param expected_metrics: set of expected metrics to be scraped
    :param event: threading event
    :return: (bool) is metrics is scraped as expected, (set) metrics that are different of expected,
    (set) metrics that exists
    """
    metrics_present = set()
    metrics = requests.get(local_prometheus_endpoint)
    for family in text_string_to_metric_families(metrics.text):
        for sample in family.samples:
            metrics_present.add(sample.name)
    metrics_dif = metrics_present.symmetric_difference(expected_metrics)
    if metrics_dif == set() or metrics_dif == {}:
        event.set()
    return event.is_set(), metrics_dif, metrics_present


def default_enabled_metric_groups_by_handler(handler):
    """

    :param (str) handler: policy handler type
    :return: default metric groups enabled
    """
    assert_that(isinstance(handler, str), equal_to(True), f"Invalid handler type {handler}. Handler must be str type")
    handler = handler.lower()
    assert_that(handler, any_of("dns", "dns-v2", "net", "net-v2", "dhcp", "bgp", "pcap", "flow", "netprobe"),
                "Invalid handler")
    groups_default_enabled = {
        "dns": ["cardinality", "counters", "dns_transaction", "top_qnames", "top_ports"],
        "dns-v2": ["cardinality", "counters", "top_qnames", "quantiles", "top_qtypes", "top_rcodes"],
        "net": ["cardinality", "counters", "top_geo", "top_ips"],
        "net-v2": ["cardinality", "counters", "top_geo", "top_ips", "quantiles"],
        "dhcp": [],
        "bgp": [],
        "pcap": [],
        "flow": ["cardinality", "counters", "by_packets", "by_bytes", "top_ips", "top_ports", "top_ips_ports"],
        "netprobe": ["counters", "histograms"]
    }
    return groups_default_enabled[handler]


def expected_metrics_by_handlers_and_groups(handler, groups_enabled, groups_disabled):
    """

    :param (str) handler: name of policy handler
    :param (str) groups_enabled: string with comma separated metric groups enables
    :param (str) groups_disabled: string with comma separated metric groups disabled
    :return: set of expected metrics generated
    """

    default_enabled_metric_groups = default_enabled_metric_groups_by_handler(handler)
    groups_enabled = groups_enabled.split(", ")
    groups_enabled = list(set(groups_enabled + default_enabled_metric_groups))
    groups_disabled = groups_disabled.split(", ")
    if all(metric_group in groups_disabled for metric_group in groups_enabled) or "all" in groups_disabled:
        if len(default_enabled_metric_groups) > 0:
            return set()
    if isinstance(handler, str) and handler.lower() == "dns":
        metric_groups = {
            "dns_rates_total",
            "dns_rates_total_sum",
            "dns_rates_total_count",
            "dns_rates_events",
            "dns_rates_events_sum",
            "dns_rates_events_count",
            "dns_top_qtype",
            "dns_top_rcode",
            "dns_wire_packets_deep_samples",
            "dns_wire_packets_events"
        }
        if ("cardinality" in groups_enabled and "cardinality" not in groups_disabled) or \
                ("all" in groups_enabled and "cardinality" not in groups_disabled):
            metric_groups.add("dns_cardinality_qname")

        if ("top_ecs" in groups_enabled and "top_ecs" not in groups_disabled) or \
                ("all" in groups_enabled and "top_ecs" not in groups_disabled):
            metric_groups.add("dns_top_asn_ecs")
            metric_groups.add("dns_top_geoLoc_ecs")
            metric_groups.add("dns_top_query_ecs")
            if "counters" in groups_enabled and "counters" not in groups_disabled:
                metric_groups.add("dns_wire_packets_query_ecs")

        if ("top_qnames_details" in groups_enabled and "top_qnames_details" not in groups_disabled and
            "top_qnames" in groups_enabled and "top_qnames" not in groups_disabled) or \
                ("all" in groups_enabled and "top_qnames_details" not in groups_disabled and "all" in groups_enabled and
                 "top_qnames" not in groups_disabled):
            metric_groups.add("dns_top_qname_by_resp_bytes")

        if ("counters" in groups_enabled and "counters" not in groups_disabled) or \
                ("all" in groups_enabled and "counters" not in groups_disabled):
            metric_groups.add("dns_wire_packets_filtered")
            metric_groups.add("dns_wire_packets_ipv4")
            metric_groups.add("dns_wire_packets_ipv6")
            metric_groups.add("dns_wire_packets_nodata")
            metric_groups.add("dns_wire_packets_noerror")
            metric_groups.add("dns_wire_packets_nxdomain")
            metric_groups.add("dns_wire_packets_queries")
            metric_groups.add("dns_wire_packets_refused")
            metric_groups.add("dns_wire_packets_replies")
            metric_groups.add("dns_wire_packets_srvfail")
            metric_groups.add("dns_wire_packets_tcp")
            metric_groups.add("dns_wire_packets_total")
            metric_groups.add("dns_wire_packets_udp")
            if ("top_ecs" in groups_enabled and "top_ecs" not in groups_disabled) or \
                    ("all" in groups_enabled and "top_ecs" not in groups_disabled):
                metric_groups.add("dns_wire_packets_query_ecs")

        if ("dns_transaction" in groups_enabled and "dns_transaction" not in groups_disabled) or \
                ("all" in groups_enabled and "dns_transaction" not in groups_disabled):
            metric_groups.add("dns_xact_counts_timed_out")
            metric_groups.add("dns_xact_counts_total")
            metric_groups.add("dns_xact_in_quantiles_us")
            metric_groups.add("dns_xact_in_quantiles_us_sum")
            metric_groups.add("dns_xact_in_quantiles_us_count")
            metric_groups.add("dns_xact_in_total")
            metric_groups.add("dns_xact_out_quantiles_us")
            metric_groups.add("dns_xact_out_quantiles_us_sum")
            metric_groups.add("dns_xact_out_quantiles_us_count")
            # todo find a way to test slow metrics
            # metric_groups.add("dns_xact_out_top_slow")
            # metric_groups.add("dns_xact_in_top_slow")
            metric_groups.add("dns_xact_out_total")
            metric_groups.add("dns_xact_ratio_quantiles")
            metric_groups.add("dns_xact_ratio_quantiles_sum")
            metric_groups.add("dns_xact_ratio_quantiles_count")

        if ("histograms" in groups_enabled and "histograms" not in groups_disabled) or \
                ("all" in groups_enabled and "histograms" not in groups_disabled):
            metric_groups.add("dns_xact_in_histogram_us_bucket")
            metric_groups.add("dns_xact_in_histogram_us_count")
            metric_groups.add("dns_xact_out_histogram_us_bucket")
            metric_groups.add("dns_xact_out_histogram_us_count")

        if ("top_qnames" in groups_enabled and "top_qnames" not in groups_disabled) or \
                ("all" in groups_enabled and "top_qnames" not in groups_disabled):
            metric_groups.add("dns_top_nodata")
            metric_groups.add("dns_top_nxdomain")
            metric_groups.add("dns_top_qname2")
            metric_groups.add("dns_top_qname3")
            metric_groups.add("dns_top_refused")
            metric_groups.add("dns_top_srvfail")
            if ("top_qnames_details" in groups_enabled and "top_qnames_details" not in groups_disabled) or \
                    ("all" in groups_enabled and "top_qnames_details" not in groups_disabled):
                metric_groups.add("dns_top_qname_by_resp_bytes")
                metric_groups.add("dns_top_noerror")

        if ("top_ports" in groups_enabled and "top_ports" not in groups_disabled) or \
                ("all" in groups_enabled and "top_ports" not in groups_disabled):
            metric_groups.add("dns_top_udp_ports")

    elif isinstance(handler, str) and handler.lower() == "dns-v2":
        metric_groups = {
            "dns_observed_packets",
            "dns_deep_sampled_packets",
            "dns_rates_observed_pps_sum",
            "dns_rates_observed_pps_count",
            "dns_rates_observed_pps"
        }
        if ("cardinality" in groups_enabled and "cardinality" not in groups_disabled) or \
                ("all" in groups_enabled and "cardinality" not in groups_disabled):
            metric_groups.add("dns_cardinality_qname")
        if ("counters" in groups_enabled and "counters" not in groups_disabled) or \
                ("all" in groups_enabled and "counters" not in groups_disabled):
            metric_groups.add("dns_xacts")
            metric_groups.add("dns_udp_xacts")
            metric_groups.add("dns_timeout_queries")
            metric_groups.add("dns_tcp_xacts")
            metric_groups.add("dns_srvfail_xacts")
            metric_groups.add("dns_refused_xacts")
            metric_groups.add("dns_orphan_responses")
            metric_groups.add("dns_nxdomain_xacts")
            metric_groups.add("dns_noerror_xacts")
            metric_groups.add("dns_nodata_xacts")
            metric_groups.add("dns_ipv6_xacts")
            metric_groups.add("dns_ipv4_xacts")
            metric_groups.add("dns_filtered_packets")
            metric_groups.add("dns_ecs_xacts")
            metric_groups.add("dns_dot_xacts")
            metric_groups.add("dns_doq_xacts")
            metric_groups.add("dns_doh_xacts")
            metric_groups.add("dns_dnscrypt_udp_xacts")
            metric_groups.add("dns_dnscrypt_tcp_xacts")
            metric_groups.add("dns_checking_disabled_xacts")
            metric_groups.add("dns_authoritative_answer_xacts")
            metric_groups.add("dns_authenticated_data_xacts")
        if ("quantiles" in groups_enabled and "quantiles" not in groups_disabled) or \
                ("all" in groups_enabled and "quantiles" not in groups_disabled):
            metric_groups.add("dns_xact_rates_sum")
            metric_groups.add("dns_xact_rates_count")
            metric_groups.add("dns_xact_rates")
        if ("top_ecs" in groups_enabled and "top_ecs" not in groups_disabled) or \
                ("all" in groups_enabled and "top_ecs" not in groups_disabled):
            metric_groups.add("dns_top_geo_loc_ecs_xacts")
            metric_groups.add("dns_top_ecs_xacts")
            metric_groups.add("dns_top_asn_ecs_xacts")
        if ("top_ports" in groups_enabled and "top_ports" not in groups_disabled) or \
                ("all" in groups_enabled and "top_ports" not in groups_disabled):
            metric_groups.add("dns_top_udp_ports_xacts")
        if ("top_qnames" in groups_enabled and "top_qnames" not in groups_disabled) or \
                ("all" in groups_enabled and "top_qnames" not in groups_disabled):
            metric_groups.add("dns_top_qname3_xacts")
            metric_groups.add("dns_top_qname2_xacts")
        if ("top_qtypes" in groups_enabled and "top_qtypes" not in groups_disabled) or \
                ("all" in groups_enabled and "top_qtypes" not in groups_disabled):
            metric_groups.add("dns_top_qtype_xacts")
        if ("top_rcodes" in groups_enabled and "top_rcodes" not in groups_disabled) or \
                ("all" in groups_enabled and "top_rcodes" not in groups_disabled):
            metric_groups.add("dns_top_srvfail_xacts")
            metric_groups.add("dns_top_refused_xacts")
            metric_groups.add("dns_top_rcode_xacts")
            metric_groups.add("dns_top_nxdomain_xacts")
            metric_groups.add("dns_top_noerror_xacts")
            metric_groups.add("dns_top_nodata_xacts")
        if ("top_size" in groups_enabled and "top_size" not in groups_disabled) or \
                ("all" in groups_enabled and "top_size" not in groups_disabled):
            metric_groups.add("dns_top_response_bytes")
            metric_groups.add("dns_response_query_size_ratio_sum")
            metric_groups.add("dns_response_query_size_ratio_count")
            metric_groups.add("dns_response_query_size_ratio")
        if ("xact_times" in groups_enabled and "xact_times" not in groups_disabled) or \
                ("all" in groups_enabled and "xact_times" not in groups_disabled):
            metric_groups.add("dns_xact_time_us_sum")
            metric_groups.add("dns_xact_time_us_count")
            metric_groups.add("dns_xact_time_us")
            metric_groups.add("dns_xact_histogram_us_bucket")
            metric_groups.add("dns_xact_histogram_us_count")
            # todo find a way to test slow metrics
            # metric_groups.add("dns_top_slow_xacts")
    elif isinstance(handler, str) and handler.lower() == "net":
        metric_groups = {
            "packets_deep_samples",
            "packets_events",
            "packets_payload_size",
            "packets_payload_size_sum",
            "packets_payload_size_count",
            "payload_rates_bytes_in",
            "payload_rates_bytes_in_sum",
            "payload_rates_bytes_in_count",
            "payload_rates_bytes_out",
            "payload_rates_bytes_out_sum",
            "payload_rates_bytes_out_count",
            "payload_rates_bytes_total",
            "payload_rates_bytes_total_sum",
            "payload_rates_bytes_total_count",
            "packets_rates_pps_events",
            "packets_rates_pps_events_sum",
            "packets_rates_pps_events_count",
            "packets_rates_pps_in",
            "packets_rates_pps_in_sum",
            "packets_rates_pps_in_count",
            "packets_rates_pps_out",
            "packets_rates_pps_out_sum",
            "packets_rates_pps_out_count",
            "packets_rates_pps_total",
            "packets_rates_pps_total_sum",
            "packets_rates_pps_total_count"
        }
        if ("cardinality" in groups_enabled and "cardinality" not in groups_disabled) or \
                ("all" in groups_enabled and "cardinality" not in groups_disabled):
            metric_groups.add("packets_cardinality_src_ips_in")
            metric_groups.add("packets_cardinality_dst_ips_out")

        if ("counters" in groups_enabled and "counters" not in groups_disabled) or \
                ("all" in groups_enabled and "counters" not in groups_disabled):
            metric_groups.add("packets_filtered")
            metric_groups.add("packets_in")
            metric_groups.add("packets_ipv4")
            metric_groups.add("packets_ipv6")
            metric_groups.add("packets_other_l4")
            metric_groups.add("packets_out")
            metric_groups.add("packets_protocol_tcp_syn")
            metric_groups.add("packets_tcp")
            metric_groups.add("packets_total")
            metric_groups.add("packets_udp")
            metric_groups.add("packets_unknown_dir")

        if ("top_geo" in groups_enabled and "top_geo" not in groups_disabled) or \
                ("all" in groups_enabled and "top_geo" not in groups_disabled):
            metric_groups.add("packets_top_ASN")
            metric_groups.add("packets_top_geoLoc")

        if ("top_ips" in groups_enabled and "top_ips" not in groups_disabled) or \
                ("all" in groups_enabled and "top_ips" not in groups_disabled):
            metric_groups.add("packets_top_ipv4")
            metric_groups.add("packets_top_ipv6")

    elif isinstance(handler, str) and handler.lower() == "net-v2":
        metric_groups = {
            "net_deep_sampled_packets",
            "net_observed_packets",
            "net_rates_observed_pps",
            "net_rates_observed_pps_count",
            "net_rates_observed_pps_sum"
        }
        if ("cardinality" in groups_enabled and "cardinality" not in groups_disabled) or \
                ("all" in groups_enabled and "cardinality" not in groups_disabled):
            metric_groups.add("net_cardinality_ips")
        if ("counters" in groups_enabled and "counters" not in groups_disabled) or \
                ("all" in groups_enabled and "counters" not in groups_disabled):
            metric_groups.add("net_filtered_packets")
            metric_groups.add("net_ipv4_packets")
            metric_groups.add("net_ipv6_packets")
            metric_groups.add("net_other_l4_packets")
            metric_groups.add("net_tcp_packets")
            metric_groups.add("net_tcp_syn_packets")
            metric_groups.add("net_total_packets")
            metric_groups.add("net_udp_packets")

        if ("top_geo" in groups_enabled and "top_geo" not in groups_disabled) or \
                ("all" in groups_enabled and "top_geo" not in groups_disabled):
            metric_groups.add("net_top_asn_packets")
            metric_groups.add("net_top_geo_loc_packets")

        if ("top_ips" in groups_enabled and "top_ips" not in groups_disabled) or \
                ("all" in groups_enabled and "top_ips" not in groups_disabled):
            metric_groups.add("net_top_ipv4_packets")
            metric_groups.add("net_top_ipv6_packets")
        if ("quantiles" in groups_enabled and "quantiles" not in groups_disabled) or \
                ("all" in groups_enabled and "quantiles" not in groups_disabled):
            metric_groups.add("net_payload_size_bytes")
            metric_groups.add("net_payload_size_bytes_count")
            metric_groups.add("net_payload_size_bytes_sum")
            metric_groups.add("net_rates_bps")
            metric_groups.add("net_rates_bps_count")
            metric_groups.add("net_rates_bps_sum")
            metric_groups.add("net_rates_pps")
            metric_groups.add("net_rates_pps_count")
            metric_groups.add("net_rates_pps_sum")

    elif isinstance(handler, str) and handler.lower() == "dhcp":
        metric_groups = {
            "dhcp_rates_total",
            "dhcp_rates_total_sum",
            "dhcp_rates_total_count",
            "dhcp_rates_events",
            "dhcp_rates_events_sum",
            "dhcp_rates_events_count",
            "dhcp_top_servers",
            "dhcp_top_clients",
            "dhcp_wire_packets_ack",
            "dhcp_wire_packets_advertise",
            "dhcp_wire_packets_deep_samples",
            "dhcp_wire_packets_discover",
            "dhcp_wire_packets_events",
            "dhcp_wire_packets_filtered",
            "dhcp_wire_packets_offer",
            "dhcp_wire_packets_reply",
            "dhcp_wire_packets_request",
            "dhcp_wire_packets_request_v6",
            "dhcp_wire_packets_solicit",
            "dhcp_wire_packets_total",
        }

    elif isinstance(handler, str) and handler.lower() == "bgp":
        metric_groups = {
            "bgp_rates_events",
            "bgp_rates_events_count",
            "bgp_rates_events_sum",
            "bgp_rates_total",
            "bgp_rates_total_count",
            "bgp_rates_total_sum",
            "bgp_wire_packets_deep_samples",
            "bgp_wire_packets_events",
            "bgp_wire_packets_filtered",
            "bgp_wire_packets_keepalive",
            "bgp_wire_packets_notification",
            "bgp_wire_packets_update",
            "bgp_wire_packets_open",
            "bgp_wire_packets_routerefresh",
            "bgp_wire_packets_total"
        }

    elif isinstance(handler, str) and handler.lower() == "pcap":
        metric_groups = {
            "pcap_if_drops",
            "pcap_os_drops",
            "pcap_tcp_reassembly_errors"
        }

    elif isinstance(handler, str) and handler.lower() == "flow":
        metric_groups = set()
        if ("cardinality" in groups_enabled and "cardinality" not in groups_disabled) or \
                ("all" in groups_enabled and "cardinality" not in groups_disabled):
            metric_groups.add("flow_cardinality_dst_ips_out")
            metric_groups.add("flow_cardinality_dst_ports_out")
            metric_groups.add("flow_cardinality_src_ips_in")
            metric_groups.add("flow_cardinality_src_ports_in")
            if ("conversations" in groups_enabled and "conversations" not in groups_disabled) or \
                    ("all" in groups_enabled and "conversations" not in groups_disabled):
                metric_groups.add("flow_cardinality_conversations")
        if ("top_tos" in groups_enabled and "top_tos" not in groups_disabled) or \
                ("all" in groups_enabled and "top_tos" not in groups_disabled):
            if ("by_bytes" in groups_enabled and "by_bytes" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_bytes" not in groups_disabled):
                metric_groups.add("flow_top_in_dscp_packets")
                metric_groups.add("flow_top_out_dscp_packets")
            if ("by_packets" in groups_enabled and "by_packets" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_packets" not in groups_disabled):
                metric_groups.add("flow_top_in_dscp_bytes")
                metric_groups.add("flow_top_out_dscp_bytes")
        if ("counters" in groups_enabled and "counters" not in groups_disabled) or \
                ("all" in groups_enabled and "counters" not in groups_disabled):
            metric_groups.add("flow_records_filtered")
            metric_groups.add("flow_records_flows")
            if ("by_bytes" in groups_enabled and "by_bytes" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_bytes" not in groups_disabled):
                metric_groups.add("flow_in_bytes")
                metric_groups.add("flow_in_ipv4_bytes")
                metric_groups.add("flow_in_ipv6_bytes")
                metric_groups.add("flow_in_other_l4_bytes")
                metric_groups.add("flow_in_tcp_bytes")
                metric_groups.add("flow_in_udp_bytes")
                metric_groups.add("flow_out_bytes")
                metric_groups.add("flow_out_ipv4_bytes")
                metric_groups.add("flow_out_ipv6_bytes")
                metric_groups.add("flow_out_other_l4_bytes")
                metric_groups.add("flow_out_tcp_bytes")
                metric_groups.add("flow_out_udp_bytes")
            if ("by_packets" in groups_enabled and "by_packets" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_packets" not in groups_disabled):
                metric_groups.add("flow_in_ipv4_packets")
                metric_groups.add("flow_in_ipv6_packets")
                metric_groups.add("flow_in_other_l4_packets")
                metric_groups.add("flow_in_packets")
                metric_groups.add("flow_in_tcp_packets")
                metric_groups.add("flow_in_udp_packets")
                metric_groups.add("flow_out_ipv4_packets")
                metric_groups.add("flow_out_ipv6_packets")
                metric_groups.add("flow_out_other_l4_packets")
                metric_groups.add("flow_out_packets")
                metric_groups.add("flow_out_tcp_packets")
                metric_groups.add("flow_out_udp_packets")
        if ("top_ips" in groups_enabled and "top_ips" not in groups_disabled) or \
                ("all" in groups_enabled and "top_ips" not in groups_disabled):
            if ("by_bytes" in groups_enabled and "by_bytes" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_bytes" not in groups_disabled):
                metric_groups.add("flow_top_in_dst_ips_bytes")
                metric_groups.add("flow_top_in_src_ips_bytes")
                metric_groups.add("flow_top_out_dst_ips_bytes")
                metric_groups.add("flow_top_out_src_ips_bytes")
            if ("by_packets" in groups_enabled and "by_packets" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_packets" not in groups_disabled):
                metric_groups.add("flow_top_in_dst_ips_packets")
                metric_groups.add("flow_top_in_src_ips_packets")
                metric_groups.add("flow_top_out_dst_ips_packets")
                metric_groups.add("flow_top_out_src_ips_packets")
        if ("top_ports" in groups_enabled and "top_ports" not in groups_disabled) or \
                ("all" in groups_enabled and "top_ports" not in groups_disabled):
            if ("by_bytes" in groups_enabled and "by_bytes" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_bytes" not in groups_disabled):
                metric_groups.add("flow_top_in_dst_ports_bytes")
                metric_groups.add("flow_top_in_src_ports_bytes")
                metric_groups.add("flow_top_out_dst_ports_bytes")
                metric_groups.add("flow_top_out_src_ports_bytes")
            if ("by_packets" in groups_enabled and "by_packets" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_packets" not in groups_disabled):
                metric_groups.add("flow_top_in_dst_ports_packets")
                metric_groups.add("flow_top_in_src_ports_packets")
                metric_groups.add("flow_top_out_dst_ports_packets")
                metric_groups.add("flow_top_out_src_ports_packets")
        if ("top_ips_ports" in groups_enabled and "top_ips_ports" not in groups_disabled) or \
                ("all" in groups_enabled and "top_ips_ports" not in groups_disabled):
            if ("by_bytes" in groups_enabled and "by_bytes" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_bytes" not in groups_disabled):
                metric_groups.add("flow_top_in_dst_ip_ports_bytes")
                metric_groups.add("flow_top_in_src_ip_ports_bytes")
                metric_groups.add("flow_top_out_dst_ip_ports_bytes")
                metric_groups.add("flow_top_out_src_ip_ports_bytes")
            if ("by_packets" in groups_enabled and "by_packets" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_packets" not in groups_disabled):
                metric_groups.add("flow_top_in_dst_ip_ports_packets")
                metric_groups.add("flow_top_in_src_ip_ports_packets")
                metric_groups.add("flow_top_out_dst_ip_ports_packets")
                metric_groups.add("flow_top_out_src_ip_ports_packets")
        if ("top_geo" in groups_enabled and "top_geo" not in groups_disabled) or \
                ("all" in groups_enabled and "top_geo" not in groups_disabled):
            if ("by_bytes" in groups_enabled and "by_bytes" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_bytes" not in groups_disabled):
                metric_groups.add("flow_top_asn_bytes")
                metric_groups.add("flow_top_geo_loc_bytes")
            if ("by_packets" in groups_enabled and "by_packets" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_packets" not in groups_disabled):
                metric_groups.add("flow_top_geo_loc_packets")
                metric_groups.add("flow_top_asn_packets")
        if ("top_interfaces" in groups_enabled and "top_interfaces" not in groups_disabled) or \
                ("all" in groups_enabled and "top_interfaces" not in groups_disabled):
            if ("by_bytes" in groups_enabled and "by_bytes" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_bytes" not in groups_disabled):
                metric_groups.add("flow_top_in_interfaces_bytes")
                metric_groups.add("flow_top_out_interfaces_bytes")
            if ("by_packets" in groups_enabled and "by_packets" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_packets" not in groups_disabled):
                metric_groups.add("flow_top_in_interfaces_packets")
                metric_groups.add("flow_top_out_interfaces_packets")
        if ("conversations" in groups_enabled and "conversations" not in groups_disabled) or \
                ("all" in groups_enabled and "conversations" not in groups_disabled):
            if ("cardinality" in groups_enabled and "cardinality" not in groups_disabled) or \
                    ("all" in groups_enabled and "cardinality" not in groups_disabled):
                metric_groups.add("flow_cardinality_conversations")
            if ("by_bytes" in groups_enabled and "by_bytes" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_bytes" not in groups_disabled):
                metric_groups.add("flow_top_conversations_bytes")
            if ("by_packets" in groups_enabled and "by_packets" not in groups_disabled) or \
                    ("all" in groups_enabled and "by_packets" not in groups_disabled):
                metric_groups.add("flow_top_conversations_packets")

    elif isinstance(handler, str) and handler.lower() == "netprobe":
        metric_groups = set()
        if ("quantiles" in groups_enabled and "quantiles" not in groups_disabled) or \
                ("all" in groups_enabled and "quantiles" not in groups_disabled):
            metric_groups.add("netprobe_response_quantiles_us")
            metric_groups.add("netprobe_response_quantiles_us_sum")
            metric_groups.add("netprobe_response_quantiles_us_count")
            if ("counters" in groups_enabled and "counters" not in groups_disabled) or \
                    ("all" in groups_enabled and "counters" not in groups_disabled):
                metric_groups.add("netprobe_response_max_us")
                metric_groups.add("netprobe_response_min_us")

        if ("counters" in groups_enabled and "counters" not in groups_disabled) or \
                ("all" in groups_enabled and "counters" not in groups_disabled):
            metric_groups.add("netprobe_attempts")
            metric_groups.add("netprobe_dns_lookup_failures")
            metric_groups.add("netprobe_packets_timeout")
            metric_groups.add("netprobe_successes")
            if (("quantiles" in groups_enabled and "quantiles" not in groups_disabled) or
                ("histograms" in groups_enabled and "histograms" not in groups_disabled)) or \
                    (("all" in groups_enabled and "quantiles" not in groups_disabled) or
                     ("all" in groups_enabled and "histograms" not in groups_disabled)):
                metric_groups.add("netprobe_response_max_us")
                metric_groups.add("netprobe_response_min_us")
        if ("histograms" in groups_enabled and "histograms" not in groups_disabled) or \
                ("all" in groups_enabled and "histograms" not in groups_disabled):
            metric_groups.add("netprobe_response_histogram_us_bucket")
            metric_groups.add("netprobe_response_histogram_us_count")
            if ("counters" in groups_enabled and "counters" not in groups_disabled) or \
                    ("all" in groups_enabled and "counters" not in groups_disabled):
                metric_groups.add("netprobe_response_max_us")
                metric_groups.add("netprobe_response_min_us")
    else:
        raise f"{handler} is not a valid handler"
    return metric_groups
