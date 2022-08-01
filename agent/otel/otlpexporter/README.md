#Integration tests for Open Telemetry Pktvisor Agent Receiver

- [ ] POC with Authenticate using KeyCloak to test authentication methods
- [ ] Understand KeyCloak to see if we can bring to the sinker service
- [ ] Implement similar auth module within otlpexporter in agent 
- [ ] Make agent connect to the POC sinker service
- [ ] Bring new sinker OTEL service to Orb 
- [ ] Present POC

```mermaid
sequenceDiagram
    participant pktvisor
    participant OrbAgent
    participant AgentOtelReceiver
    participant AgentOtelExporter
    participant OrbSinker
    participant SinkerOtelReceiver
    participant SinkerOtelExporter
    participant Prometheus
    
    pktvisor-->>OrbAgent: metrics via HTTP
    OrbAgent-->>AgentOtelReceiver: metrics via neutral Otel format
    AgentOtelReceiver-->>AgentOtelExporter: pass the config and needs to authenticate with OTLP
    OrbAgent-->>OrbSinker: new transport over gRPC secured with API Key
    OrbSinker-->>SinkerOtelReceiver: send metrics from gRPC and fetch info from the OTEL sink
    SinkerOtelReceiver-->>SinkerOtelExporter: connects to TSDB from config and OTEL sink URL
    SinkerOtelExporter-->>Prometheus: from config
```

