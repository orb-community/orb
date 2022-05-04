# OTLP/MQTT Exporter


Alt 1. Test Orb Agent with Open Telemetry
use this config in localconifg/config.yaml
```yaml
otel:
  enable: true
```

Exports traces and/or metrics via MQTT using [OTLP](
https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/otlp.md)
format.

Supported pipeline types: traces, metrics, logs

:warning: OTLP logs format is currently marked as "Beta" and may change in
incompatible ways.

The following settings are required:

- `endpoint` (no default): The target base URL to send data to (e.g.: https://example.com:4318).
  To send each signal a corresponding path will be added to this base URL, i.e. for traces
  "/v1/traces" will appended, for metrics "/v1/metrics" will be appended, for logs
  "/v1/logs" will be appended. 

The following settings can be optionally configured:

- `traces_channel` (no default): The target URL to send trace data to (e.g.: https://example.com:4318/v1/traces).
   If this setting is present the `endpoint` setting is ignored for traces.
- `metrics_channel` (no default): The target URL to send metric data to (e.g.: https://example.com:4318/v1/metrics).
   If this setting is present the `endpoint` setting is ignored for metrics.
- `logs_channel` (no default): The target URL to send log data to (e.g.: https://example.com:4318/v1/logs).
   If this setting is present the `endpoint` setting is ignored logs.
- `tls`: see [TLS Configuration Settings](../../config/configtls/README.md) for the full set of available options.

[//]: # (-Not sure yet if this will apply `read_buffer_size` &#40;default = 0&#41;: ReadBufferSize for MQTT client.)

[//]: # (-Not sure yet if this will apply `write_buffer_size` &#40;default = 512 * 1024&#41;: WriteBufferSize for HTTP client.)

Example:

```yaml
exporters:
  otlpmqtt:
    endpoint: https://example.com:4318/v1/traces
```

By default `gzip` compression is enabled. See [compression comparison](../../config/configgrpc/README.md#compression-comparison) for details benchmark information. To disable, configure as follows:

```yaml
exporters:
  otlpmqtt:
    ...
    compression: none
```

The full list of settings exposed for this exporter are documented [here](./config.go)
with detailed sample configurations [here](./testdata/config.yaml).
