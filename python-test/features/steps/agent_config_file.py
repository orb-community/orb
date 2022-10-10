import yaml
from utils import create_tags_set

from taps import *


class FleetAgent:
    def __init__(self):
        pass

    @classmethod
    def config_file_of_orb_agent(cls, name, token, iface, orb_url, base_orb_mqtt, tap_name, tls_verify=True,
                                 auto_provision=True, orb_cloud_mqtt_id=None, orb_cloud_mqtt_key=None,
                                 orb_cloud_mqtt_channel_id=None, input_type="pcap", input_tags='3', settings=None):
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
                            "config_file": f"/usr/local/orb/{name}.yaml"
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
                            "config_file": f"/usr/local/orb/{name}.yaml"
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
        agent = yaml.dump(agent)
        return agent, tap.taps
