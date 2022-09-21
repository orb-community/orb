receivers:
  kafka:
    brokers:
      - kafka1:19092
    topic: otlp_metrics
    protocol_version: 2.0.0

processors:
  batch:

extensions:
  health_check:
    check_collector_pipeline:
      enabled: true
      interval: "1m"
      exporter_failure_threshold: 5

  pprof:
    endpoint: :1888

  zpages:
    endpoint: :55679

  basicauth/client:
    client_auth:
      username: ${USERNAME}
      password: ${PASSWORD}

exporters:
  prometheusremotewrite:
    endpoint: ${PROM_URL}
    auth:
      authenticator: basicauth/client
  logging:
    loglevel: debug

service:
  extensions: [pprof, zpages, health_check, basicauth/client]
  pipelines:
    metrics:
      receivers: [kafka]
      exporters: [prometheusremotewrite]
