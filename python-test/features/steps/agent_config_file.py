class FleetAgent:
    def __init__(self):
        pass

    @classmethod
    def config_file_of_agent_tap_pcap(cls, name, token, iface, orb_url, base_orb_url, tls_verify="true"):
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
      auto_provision: true

    api:
      address: {orb_url}
      token: {token}
    mqtt:
      address: tls://{base_orb_url}:8883

            """
        return agent_tap_pcap
