# How Orb Agent would send info in otlp to the Orb Sink

Orb Agent fetches information from Pktvisor using a receiver pktvisorreceiver in package that implements a customized receiver from opentelemetry.

In the PR [1428](https://github.com/etaques/orb/pull/1428), the orb-agent has now the opentelemetry exporter that will, pass the otlp through MQTT, through the usual channels that orb-sinker receives the information.

In a sequence Diagram as follows
```mermaid
sequenceDiagram
    loop Every minute
    Agent-->>Sinker: Metrics in OTLP format
    Sinker->>Collector: pass through GRPC
    Note right of Collector: Custom Collector with Extensions
    Collector->>Sink: pass using Sink OTEL exporter
    end
```

Let's expand a bit the visio on what will happen in between Collector and Sink.

```mermaid
sequenceDiagram
    Sinker->>Collector: OTLP payload: metrics, policy, agentId
    Collector->>Collector: OTLPReceiver via GRPC
    Note over Collector,Orb: I will reference Orb, as general service here, because of TBDs
    Collector->>Orb: Extension 1, send policy, agentId
    Orb->>Collector: Retrieve Sinks
    Collector->>Sinks: S3OTLPExporter
    Collector->>Sinks: PrometheusOTLPExporter
    Collector->>Sinks: DatadogOTLPExporter
```

The Collector would be a custom implementation of the [Open Telemetry Collector](https://github.com/open-telemetry/opentelemetry-collector-contrib)

Here is a diagram of that

```mermaid
sequenceDiagram
    OTLPReceiver->>CustomExtension1: agentId, policyId
    CustomExtension1->>Orb: Fetch Dataset(s)
    Orb->>CustomExtension1: Respond ds for policyId, agentId
    CustomExtension1->>CustomExtension2: Select sinks from DataSet(s)
    CustomExtension2->>Orb: GetSinkCredentials
    Orb->>CustomExtention2: Each sink credential per policy
    CustomExtension2->>Exporters: metrics and credentials to ship to sink.
```

Here is another diagram of that.

![OTLP-Orb-Collector](./OTLP-Orb-Collector.png)


## Concurrency and Scaling

TDB

Think of strategies of scaling

Sharding could be one

## Pktvisor scraping metrics and how OTLP format can change that

