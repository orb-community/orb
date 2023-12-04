import yaml
from utils import create_tags_set
from taps import *
from logger import Logger

log = Logger().logger_instance()

default_path_config_file = "/opt/orb"


class ConfigFiles:
    backend_type = None

    def __init__(self, config_type):
        assert_that(config_type, any_of(equal_to("pktvisor"), equal_to("otel"), equal_to("orb_agent")),
                    f"Unsuported agent backend type: {config_type}")
        self.config_type = config_type

    def pktvisor_config_file(self, iface, tap_name, input_type="pcap", input_tags="3", settings=None):
        assert_that(self.config_type, equal_to("pktvisor"))
        assert_that(input_type, any_of(equal_to("pcap"), equal_to("flow"), equal_to("dnstap"), equal_to("netprobe")),
                    "Unexpect type of input type.")
        if "iface" in settings.keys() and isinstance(settings["iface"], str) and \
                (settings["iface"].lower() == "default" or settings["iface"].lower() == "mocked"):
            settings['iface'] = iface
        tap = Taps()
        if input_type == "pcap":
            tap.add_pcap(tap_name, **settings)
        elif input_type == "flow":
            tap.add_flow(tap_name, **settings)
        elif input_type == "netprobe":
            tap.add_netprobe(tap_name, **settings)
        else:
            tap.add_dnstap(tap_name, **settings)
        if input_tags is not None and input_tags != '0':
            tap_tags = create_tags_set(input_tags, tag_prefix='testtaptag', string_mode='lower')
            tap.add_tag(tap_name, tap_tags)
        pktvisor_backend = \
            {"visor": {
                "taps": tap.taps
            }}
        return pktvisor_backend, tap.taps

    def otel_config_file(self):
        assert_that(self.config_type, equal_to("otel"))
        return dict()

    def orb_agent_config_file(self, backend_type, config_file, auto_provision, port, backend_settings, cloud_settings,
                              orb_url, base_orb_mqtt, tls_verify=True):
        assert_that(self.config_type, equal_to("orb_agent"))
        agent = AgentConfigs(tls_verify)
        agent.add_backend(backend_type, config_file, port, **backend_settings)
        agent.set_cloud(auto_provision, orb_url, base_orb_mqtt, **cloud_settings)
        log.debug(f"Orb Agent Config File: {agent.config}")
        return agent.config


class AgentConfigs:
    def __init__(self, tls_verify=True):
        assert_that(tls_verify, any_of(equal_to(True), equal_to(False)), "Unexpected value for tls_verify on "
                                                                         "agent config file creation")
        self.config = {
            "backends": {},
            "cloud": {},
            "tls": {"verify": tls_verify}
        }

    def add_backend(self, backend_type, config_file, port, **settings):
        backend = {backend_type: {"config_file": config_file}}
        backend[backend_type].update(settings)
        if backend_type == "pktvisor":
            assert_that(settings, has_key("binary"), "Missing binary for pktvisor backend")
            backend[backend_type].update({"api_port": port})
        elif backend_type == "otel":
            backend[backend_type].update({"otlp_port": port})
        self.config["backends"].update(backend)
        return self.config

    def remove_backend(self, backend):
        self.config["backends"].pop(backend, None)
        return self.config

    def set_cloud(self, auto_provision, orb_url, base_orb_mqtt, **settings):
        assert_that(auto_provision, any_of(equal_to(True), equal_to(False)), "Unexpected value for auto_provision "
                                                                             "on agent config file creation")
        log.debug(f"Setting for agent cloud: {settings}")

        cloud_config = {"auto_provision": auto_provision}
        cloud_api = {"address": orb_url}
        cloud_mqtt = {"address": base_orb_mqtt}

        if auto_provision:
            assert_that(settings, has_key("name"), "Missing agent name for auto provision agent")
            assert_that(settings, has_key("token"), "Missing token for auto provision agent")
            cloud_config.update({"agent_name": settings["name"]})
            cloud_api.update({"token": settings["token"]})
        else:
            assert_that(settings, has_key("orb_cloud_mqtt_id"), "Missing id for cloud mqtt")
            assert_that(settings, has_key("orb_cloud_mqtt_key"), "Missing key for cloud mqtt")
            assert_that(settings, has_key("orb_cloud_mqtt_channel_id"), "Missing channel id for cloud mqtt")
            mqtt_configs = {"id": settings["orb_cloud_mqtt_id"], "key": settings["orb_cloud_mqtt_key"],
                            "channel_id": settings["orb_cloud_mqtt_channel_id"]}
            cloud_mqtt.update(mqtt_configs)

        cloud = {"api": cloud_api, "mqtt": cloud_mqtt, "config": cloud_config}

        self.config["cloud"].update(cloud)
        return self.config


class FleetAgent:
    def __init__(self):
        pass

    @classmethod
    def config_file_of_orb_agent(cls, agent_name, backend_type, auto_provision, port, backend_settings,
                                 cloud_settings, orb_url, base_orb_mqtt, tls_verify=True, backend_file=None,
                                 config_file="default", overwrite_default=False):
        assert_that(tls_verify, any_of(equal_to(True), equal_to(False)), "Unexpected value for tls_verify on "
                                                                         "agent pcap config file creation")
        assert_that(overwrite_default, any_of(equal_to(True), equal_to(False)),
                    "Unexpected value for overwrite_default")
        if overwrite_default is True:
            name_agent_file = "agent"
        else:
            name_agent_file = agent_name
        if config_file == "default":
            config_file = f"{default_path_config_file}/{name_agent_file}.yaml"
        orb_agent_file = ConfigFiles("orb_agent").orb_agent_config_file(backend_type, config_file, auto_provision, port,
                                                                        backend_settings, cloud_settings,
                                                                        orb_url, base_orb_mqtt, tls_verify=True)
        if backend_file is not None:
            assert_that(isinstance(backend_file, dict), equal_to(True), f"Invalid backend file: {backend_file}")
            orb_agent_file.update(backend_file)
        agent = yaml.dump(orb_agent_file)
        log.debug(f"Orb Agent Config File: {agent}")
        return agent

    # @classmethod
    # def config_file_of_orb_agent(cls, name, token, orb_url, base_orb_mqtt, tls_verify=True, auto_provision=True,
    #                              orb_cloud_mqtt_id=None, orb_cloud_mqtt_key=None, orb_cloud_mqtt_channel_id=None,
    #                              overwrite_default=False):
    #     assert_that(tls_verify, any_of(equal_to(True), equal_to(False)), "Unexpected value for tls_verify on "
    #                                                                      "agent pcap config file creation")
    #
    #     assert_that(auto_provision, any_of(equal_to(True), equal_to(False)), "Unexpected value for auto_provision "
    #                                                                          "on agent pcap config file creation")
    #     if overwrite_default is True:
    #         pkt_name_file = "agent"
    #     else:
    #         pkt_name_file = name
    #     if auto_provision:
    #         agent = {
    #             "orb": {
    #                 "backends": {
    #                     "pktvisor": {
    #                         "binary": "/usr/local/sbin/pktvisord",
    #                         "config_file": f"{default_path_config_file}/{pkt_name_file}.yaml"
    #                     }
    #                 },
    #                 "tls": {
    #                     "verify": tls_verify
    #                 },
    #                 "cloud": {
    #                     "config": {
    #                         "auto_provision": auto_provision,
    #                         "agent_name": name
    #                     },
    #                     "api": {
    #                         "address": orb_url,
    #                         "token": token
    #                     },
    #                     "mqtt": {
    #                         "address": base_orb_mqtt
    #                     }
    #                 }
    #             }
    #         }
    #     else:
    #         assert_that(orb_cloud_mqtt_id, not_(is_(None)), "orb_cloud_mqtt_id must have a valid value")
    #         assert_that(orb_cloud_mqtt_channel_id, not_(is_(None)), "orb_cloud_mqtt_channel_id must have a valid value")
    #         assert_that(orb_cloud_mqtt_key, not_(is_(None)), "orb_cloud_mqtt_key must have a valid value")
    #         agent = {
    #             "version": "1.0",
    #             "orb": {
    #                 "backends": {
    #                     "pktvisor": {
    #                         "binary": "/usr/local/sbin/pktvisord",
    #                         "config_file": f"{default_path_config_file}/{pkt_name_file}.yaml"
    #                     }
    #                 },
    #                 "tls": {
    #                     "verify": tls_verify
    #
    #                 },
    #                 "cloud": {
    #                     "config": {
    #                         "auto_provision": auto_provision
    #                     },
    #                     "api": {
    #                         "address": orb_url
    #                     },
    #                     "mqtt": {
    #                         "address": base_orb_mqtt,
    #                         "id": orb_cloud_mqtt_id,
    #                         "key": orb_cloud_mqtt_key,
    #                         "channel_id": orb_cloud_mqtt_channel_id
    #                     }
    #                 }
    #             }
    #         }
    #     agent = yaml.dump(agent)
    #     # return agent, tap.taps
    #     return agent, None  # todo arrumar
    #
    # @classmethod
    # def config_file_of_orb_agent_with_otel_backend(cls, name, token, orb_url, base_orb_mqtt, tls_verify=True,
    #                                                auto_provision=True, orb_cloud_mqtt_id=None, orb_cloud_mqtt_key=None,
    #                                                orb_cloud_mqtt_channel_id=None, overwrite_default=False):
    #     assert_that(tls_verify, any_of(equal_to(True), equal_to(False)), "Unexpected value for tls_verify on "
    #                                                                      "agent pcap config file creation")
    #
    #     assert_that(auto_provision, any_of(equal_to(True), equal_to(False)), "Unexpected value for auto_provision "
    #                                                                          "on agent pcap config file creation")
    #     if overwrite_default is True:
    #         otel_name_file = "agent"
    #     else:
    #         otel_name_file = name
    #     if auto_provision:
    #         agent = {
    #             "version": "1.0",
    #             "orb": {
    #                 "backends": {
    #                     "otel": {
    #                         "config_file": f"{default_path_config_file}/{otel_name_file}.yaml"
    #                     }
    #                 },
    #                 "tls": {
    #                     "verify": tls_verify
    #                 },
    #                 "cloud": {
    #                     "config": {
    #                         "auto_provision": auto_provision,
    #                         "agent_name": name
    #                     },
    #                     "api": {
    #                         "address": orb_url,
    #                         "token": token
    #                     },
    #                     "mqtt": {
    #                         "address": base_orb_mqtt
    #                     }
    #                 }
    #             }
    #         }
    #     else:
    #         assert_that(orb_cloud_mqtt_id, not_(is_(None)), "orb_cloud_mqtt_id must have a valid value")
    #         assert_that(orb_cloud_mqtt_channel_id, not_(is_(None)), "orb_cloud_mqtt_channel_id must have a valid value")
    #         assert_that(orb_cloud_mqtt_key, not_(is_(None)), "orb_cloud_mqtt_key must have a valid value")
    #         agent = {
    #             "version": "1.0",
    #             "orb": {
    #                 "backends": {
    #                     "otel": {
    #                         "config_file": f"{default_path_config_file}/{otel_name_file}.yaml"
    #                     }
    #                 },
    #                 "tls": {
    #                     "verify": tls_verify
    #
    #                 },
    #                 "cloud": {
    #                     "config": {
    #                         "auto_provision": auto_provision
    #                     },
    #                     "api": {
    #                         "address": orb_url
    #                     },
    #                     "mqtt": {
    #                         "address": base_orb_mqtt,
    #                         "id": orb_cloud_mqtt_id,
    #                         "key": orb_cloud_mqtt_key,
    #                         "channel_id": orb_cloud_mqtt_channel_id
    #                     }
    #                 }
    #             }
    #         }
    #     agent = yaml.dump(agent)
    #     return agent

# amanda = ConfigFiles("orb_agent").orb_agent_config_file("otel", "path/to_config",
#                                                         False,
#                                                         10853, {}, {"orb_cloud_mqtt_id": "amanda",
#                                                                     "orb_cloud_mqtt_key": "amanda_token",
#                                                                     "orb_cloud_mqtt_channel_id": "eu"}, "https:orb.com",
#                                                         "https:mqtt.com",
#                                                         tls_verify=True)
