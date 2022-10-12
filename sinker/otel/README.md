# Sink-Collector Sequence Diagrams

## On Create Sink

The new service, which will trigger K8s New pods, is called Maestro for now

## Create

```mermaid
sequenceDiagram
    autoNumber 1
    User->>Sinks: Creates a new sink
    activate RedisConfigYaml
    Sinks->>RedisConfigYaml: Creates a new YAML configuration entry in topic
    User->>Sinks: Assigns sink in Dataset, data will start flow
    autoNumber 1
    critical Agent scraped metrics flowing through Sinker with Sink ID 222
        Sinker->>Maestro: Checks Otel-Collector deployed for Sink 222
        Maestro->>Sinker: Responds that needs to create
        Sinker->>RedisConfigYaml: Asks for ConfigYaml
        activate RedisEntryCollector
        Sinker->>RedisEntryCollector: Creates new deployment yaml entry with TTL with duration of {TBD}
        RedisEntryCollector-->>Maestro: Listens to new deployment entry 
        Maestro->>Kubernetes: Create otel-collector pod with Sink 222 config YAML
    end 
    deactivate RedisConfigYaml
    deactivate RedisEntryCollector
```

## Update

```mermaid
sequenceDiagram
    autoNumber 1
    activate RedisConfigYaml
    User->>Sinks: Updates Sink Information
    Sinks->>RedisConfigYaml: Updates YAML configuration entry
    autoNumber 1
    critical Metrics flowing
        Maestro-->>RedisConfigYaml: Listens to changes on YAML
        Maestro->>RedisEntryCollector: Generates a new deployment with updated configuration
        activate RedisEntryCollector
        Maestro->>Kubernetes: Updates deployment with new values     
    end
    deactivate RedisConfigYaml
    deactivate RedisEntryCollector
```

## Delete

```mermaid
sequenceDiagram
    autoNumber 1
    activate RedisConfigYaml
    activate RedisEntryCollector
    User->>Sinks: Deletes a sink
    Sinks->>RedisConfigYaml: Removes YAML configuration entry
    deactivate RedisConfigYaml
    autoNumber 1
    critical Metrics
    option Time to Live expires
        RedisEntryCollector-->>Maestro: Entry TTL expired
        deactivate RedisEntryCollector
    option Maestro listens to ConfigYaml
        Maestro-->>RedisConfigYaml: Listens to topic's entry removal
        Maestro->>RedisEntryCollector: Removes entry
    end
    Maestro->>Kubernetes: Removes pod, and inactivate deployment YAML
```
