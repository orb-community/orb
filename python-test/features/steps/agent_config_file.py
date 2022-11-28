import yaml
from utils import create_tags_set

from taps import *

default_path_config_file = "/opt/orb"


class FleetAgent:
    def __init__(self):
        pass

    @classmethod
    def config_file_of_orb_agent(cls, name, token, iface, orb_url, base_orb_mqtt, tap_name, tls_verify=True,
                                 auto_provision=True, orb_cloud_mqtt_id=None, orb_cloud_mqtt_key=None,
                                 orb_cloud_mqtt_channel_id=None, input_type="pcap", input_tags='3', settings=None,
                                 include_otel_env_var=False, enable_otel=False, overwrite_default=False):
        if isinstance(include_otel_env_var, str):
            assert_that(include_otel_env_var.lower(), any_of("true", "false"), "Unexpected value for "
                                                                               "'include_otel_env_var'.")
            include_otel_env_var = eval(include_otel_env_var.title())
        else:
            assert_that(include_otel_env_var, any_of(False, True), "Unexpected value for 'include_otel_env_var'")
        if isinstance(enable_otel, str):
            assert_that(enable_otel.lower(), any_of("true", "false"), "Unexpected value for "
                                                                      "'enable_otel'.")
            enable_otel = eval(enable_otel.title())
        else:
            assert_that(enable_otel, any_of(False, True), "Unexpected value for 'enable_otel'")

            assert_that(tls_verify, any_of(equal_to(True), equal_to(False)), "Unexpected value for tls_verify on "
                                                                             "agent pcap config file creation")

        assert_that(auto_provision, any_of(equal_to(True), equal_to(False)), "Unexpected value for auto_provision "
                                                                             "on agent pcap config file creation")
        assert_that(input_type, any_of(equal_to("pcap"), equal_to("flow"), equal_to("dnstap")),
                    "Unexpect type of input type.")
        if "iface" in settings.keys() and settings["iface"] == "default":
            settings['iface'] = iface
        tap = Taps()
        if input_type == "pcap":
            tap.add_pcap(tap_name, **settings)
        elif input_type == "flow":
            tap.add_flow(tap_name, **settings)
        else:
            tap.add_dnstap(tap_name, **settings)
        if input_tags is not None and input_tags != '0':
            tap_tags = create_tags_set(input_tags, tag_prefix='testtaptag', string_mode='lower')
            tap.add_tag(tap_name, tap_tags)
        if overwrite_default is True:
            pkt_name_file = "agent"
        else:
            pkt_name_file = name
        if auto_provision:
            agent = {
                "version": "1.0",
                "visor": {
                    "taps": tap.taps
                },
                "orb": {
                    "backends": {
                        "pktvisor": {
                            "binary": "/usr/local/sbin/pktvisord",
                            "config_file": f"{default_path_config_file}/{pkt_name_file}.yaml"
                        }
                    },
                    "tls": {
                        "verify": tls_verify
                    },
                    "cloud": {
                        "config": {
                            "auto_provision": auto_provision,
                            "agent_name": name
                        },
                        "api": {
                            "address": orb_url,
                            "token": token
                        },
                        "mqtt": {
                            "address": base_orb_mqtt
                        }
                    }
                }
            }
            if include_otel_env_var is True:
                agent['orb']['otel'] = {"enable": enable_otel}
        else:
            assert_that(orb_cloud_mqtt_id, not_(is_(None)), "orb_cloud_mqtt_id must have a valid value")
            assert_that(orb_cloud_mqtt_channel_id, not_(is_(None)), "orb_cloud_mqtt_channel_id must have a valid value")
            assert_that(orb_cloud_mqtt_key, not_(is_(None)), "orb_cloud_mqtt_key must have a valid value")
            agent = {
                "version": "1.0",
                "visor": {
                    "taps": tap.taps
                },
                "orb": {
                    "backends": {
                        "pktvisor": {
                            "binary": "/usr/local/sbin/pktvisord",
                            "config_file": f"{default_path_config_file}/{pkt_name_file}.yaml"
                        }
                    },
                    "tls": {
                        "verify": tls_verify

                    },
                    "cloud": {
                        "config": {
                            "auto_provision": auto_provision
                        },
                        "api": {
                            "address": orb_url
                        },
                        "mqtt": {
                            "address": base_orb_mqtt,
                            "id": orb_cloud_mqtt_id,
                            "key": orb_cloud_mqtt_key,
                            "channel_id": orb_cloud_mqtt_channel_id
                        }
                    }
                }
            }
            if include_otel_env_var is True:
                agent['orb']['otel'] = {"enable": enable_otel}
        agent = yaml.dump(agent)
        return agent, tap.taps
