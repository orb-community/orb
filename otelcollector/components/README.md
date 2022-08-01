# OtelCol Components and How Orb uses them


## Sinker OTLP
In the Sinker, where we receive metrics from the orb-agents in the handleMetrics function

```mermaid
sequenceDiagram
    autonumber
    activate Sinker
    loop handleMetrics
        Sinker->>Orb: Fetches Sinks data from Orb
        Sinker->>OTLPReceiver: Passes the OTLP through a custom receiver
        OTLPReceiver->>AttributeProcessor: Add sink data to OTLP package
        AttributeProcessor->>OTLPExporter: Exports OTLP with custom attrs to OtelCol
    end
    deactivate Sinker

```

## OtelCol

Here is how the orb otel-collector should work 

```mermaid
sequenceDiagram
    autonumber
    loop handleMetrics
        OTLPReceiver->>GroupByAttributeProcessor: separate from MFOwnerID and AgentID 
        GroupByAttributeProcessor->>RoutingProcessor: route to sinks within OTLP
        RoutingProcessor->>Exporters: export data
    end
```

To test and check performance, these processors could be activated after Routing

```go
import (
	"go.opentelemetry.io/collector/processor/batchprocessor"
        "go.opentelemetry.io/collector/processor/memorylimiterprocessor"
)
_ = []component.Factories{
    // current version and stability for metrics [ 0.56.0 , beta ]
    batchprocessor.NewFactory(),
    // current version and stability for metrics [ 0.56.0 , beta ]
    memorylimiterprocessor.NewFactory(),
}
```
