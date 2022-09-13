receivers:
  kafka:
    brokers:
      - kafka1:9092
    protocol_version: 2.0.0


processors:
  batch:

extensions:
  health_check:
  pprof:
    endpoint: :1888
  zpages:
    endpoint: :55679
  basicauth/exporter:
    client_auth:
      username: $USERNAME
      password: $PASSWORD

exporters:
  prometheusremotewrite:
    endpoint: $URL
    auth:
      authenticator: basicauth/exporter
  logging:

service:
  extensions: [pprof, zpages, health_check]
  pipelines:
    metrics:
      receivers: [kafka]
      processors: [batch]
      exporters: [logging, prometheusremotewrite]
