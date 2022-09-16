receivers:
  kafka:
    brokers:
      - kafka1:19092
    protocol_version: 2.0.0

processors:
  batch:

extensions:
  health_check: {}

  pprof:
    endpoint: :1888

  zpages:
    endpoint: :55679

  basicauth/client:
    client_auth:
      username: $USERNAME
      password: $PASSWORD

exporters:
  prometheusremotewrite:
    endpoint: $URL
    auth:
      authenticator: basicauth/client
  logging:

service:
  extensions: [pprof, zpages, health_check, basicauth/client]
  pipelines:
    metrics:
      receivers: [kafka]
      processors: [batch]
      exporters: [logging, prometheusremotewrite]
