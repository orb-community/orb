import yaml
from utils import return_port_to_run_docker_container

from taps import *


class FleetAgent:
    def __init__(self):
        pass

    @classmethod
    def config_file_of_orb_agent(cls, name, token, iface, orb_url, base_orb_mqtt, tap_name, tls_verify="true",
                                 auto_provision="true", orb_cloud_mqtt_id=None, orb_cloud_mqtt_key=None,
                                 orb_cloud_mqtt_channel_id=None, input_type="pcap", settings=None):
        assert_that(tls_verify, any_of(equal_to("true"), equal_to("false")), "Unexpected value for tls_verify on "
                                                                             "agent pcap config file creation")
        assert_that(auto_provision, any_of(equal_to("true"), equal_to("false")), "Unexpected value for auto_provision "
                                                                                 "on agent pcap config file creation")
        assert_that(input_type, any_of(equal_to("pcap"), equal_to("flow"), equal_to("dnstap")),
                    "Unexpect type of input type.")
        if "iface" in settings.keys() and settings["iface"] == "default":
            settings['iface'] = iface
        if input_type == "pcap":
            tap = Taps.pcap(tap_name, input_type, settings)
        elif input_type == "flow":
            tap = Taps.flow(tap_name, input_type, settings)
        else:
            tap = Taps.dnstap(tap_name, input_type, settings)
        if auto_provision == "true":
            agent = {
                "version": "1.0",
                "visor": {
                    "taps": tap
                },
                "orb": {
                    "backends": {
                        "pktvisor": {
                            "binary": "/usr/local/sbin/pktvisord",
                            "config_file": f"/usr/local/orb/{name}.yaml"
                        }
                    },
                    "tls": {
                        "verify": {
                            "tls_verify": tls_verify
                        }
                    },
                    "cloud": {
                        "config": {
                            "auto_provision":  auto_provision,
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
        else:
            assert_that(orb_cloud_mqtt_id, not_(is_(None)), "orb_cloud_mqtt_id must have a valid value")
            assert_that(orb_cloud_mqtt_channel_id, not_(is_(None)), "orb_cloud_mqtt_channel_id must have a valid value")
            assert_that(orb_cloud_mqtt_key, not_(is_(None)), "orb_cloud_mqtt_key must have a valid value")
            agent = {
                "version": "1.0",
                "visor": {
                    "taps": tap
                },
                "orb": {
                    "backends": {
                        "pktvisor": {
                            "binary": "/usr/local/sbin/pktvisord",
                            "config_file": f"/usr/local/orb/{name}.yaml"
                        }
                    },
                    "tls": {
                        "verify": {
                            "tls_verify": tls_verify
                        }
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
        agent = yaml.dump(agent)
        return agent, tap
