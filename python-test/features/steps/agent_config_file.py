from hamcrest import *


class FleetAgent:
    def __init__(self):
        pass

    @classmethod
    def config_file_of_agent_tap_pcap(cls, name, token, iface, orb_url, base_orb_mqtt, tls_verify="true",
                                      auto_provision="true", orb_cloud_mqtt_id=None, orb_cloud_mqtt_key=None,
                                      orb_cloud_mqtt_channel_id=None):
        assert_that(tls_verify, any_of(equal_to("true"), equal_to("false")), "Unexpected value for tls_verify on "
                                                                             "agent pcap config file creation")
        assert_that(auto_provision, any_of(equal_to("true"), equal_to("false")), "Unexpected value for auto_provision "
                                                                                 "on agent pcap config file creation")
        if auto_provision == "true":
            agent_tap_pcap = f"""
    version: "1.0"
    
    visor:
      taps:
        default_pcap:
          input_type: pcap
          config:
            iface: {iface}
            host_spec: "192.168.0.54/32,192.168.0.55/32,127.0.0.1/32"            
    orb:
      backends:
        pktvisor:
          binary: "/usr/local/sbin/pktvisord"
          config_file: /usr/local/orb/{name}.yaml
      tls:
        verify: {tls_verify}
      cloud:
        config:
          agent_name: {name}
          auto_provision: {auto_provision}
    
        api:
          address: {orb_url}
          token: {token}
        mqtt:
          address: {base_orb_mqtt}
    
                """
        else:
            assert_that(orb_cloud_mqtt_id, not_(is_(None)), "orb_cloud_mqtt_id must have a valid value")
            assert_that(orb_cloud_mqtt_channel_id, not_(is_(None)), "orb_cloud_mqtt_channel_id must have a valid value")
            assert_that(orb_cloud_mqtt_key, not_(is_(None)), "orb_cloud_mqtt_key must have a valid value")
            agent_tap_pcap = f"""
        version: "1.0"
        
        visor:
          taps:
            default_pcap:
              input_type: pcap
              config:
                iface: {iface}
                host_spec: "192.168.0.54/32,192.168.0.55/32,127.0.0.1/32"            
        orb:
          backends:
            pktvisor:
              binary: "/usr/local/sbin/pktvisord"
              config_file: /usr/local/orb/{name}.yaml
          tls:
            verify: {tls_verify}
          cloud:
            config:
              auto_provision: {auto_provision}
        
            api:
              address: {orb_url}
            mqtt:
              address: {base_orb_mqtt}
              id: {orb_cloud_mqtt_id}
              key: {orb_cloud_mqtt_key}
              channel_id: {orb_cloud_mqtt_channel_id}
        
                    """
        return agent_tap_pcap
