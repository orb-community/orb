export const POLICY_OTEL_CONFIG_YAML =
`receivers:
  httpcheck:
    targets:
      - endpoint: http://orb.community
        method: GET
      - endpoint: https://orb.live
        method: GET
    collection_interval: 60s
extensions:
exporters:
service:
  pipelines:
    metrics:
      exporters:
      receivers:
        - httpcheck
`;

export const POLICY_PKTVISOR_CONFIG_YAML =
`handlers:
  modules:
    default_dns:
      type: dns
    default_net:
      type: net
input:
  input_type: pcap
  tap: default_pcap
kind: collection
`;

export const POLICY_OTEL_CONFIG_JSON =
`{
  "receivers": {
    "httpcheck": {
      "targets": [
        {
          "endpoint": "http://orb.community",
          "method": "GET"
        },
        {
          "endpoint": "https://orb.live",
          "method": "GET"
        }
      ],
      "collection_interval": "60s"
    }
  },
  "extensions": null,
  "exporters": null,
  "service": {
    "pipelines": {
      "metrics": {
        "exporters": null,
        "receivers": [
          "httpcheck"
        ]
      }
    }
  }
}`;

export const POLICY_PKTVISOR_CONFIG_JSON =
`{
  "handlers": {
    "modules": {
      "default_dns": {
        "type": "dns"
      },
      "default_net": {
        "type": "net"
      }
    }
  },
  "input": {
    "input_type": "pcap",
    "tap": "default_pcap"
  },
  "kind": "collection"
}`;

// yet to be defined
export const POLICY_DIODE_CONFIG_YAML = ``;

export const POLICY_DIODE_CONFIG_JSON =
`{

}`;

export const SINK_OTLP_CONFIG_YAML =
`authentication:
  type: basicauth
  password: ""
  username: ""
exporter:
  endpoint: ""
`;
export const SINK_PROMETHEUS_CONFIG_YAML =
`authentication:
  type: basicauth
  password: ""
  username: ""
exporter:
  remote_host: ""
`;


