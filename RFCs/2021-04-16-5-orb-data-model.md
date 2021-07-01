## Orb Data Model

Orb manages pktvisor configuration in a central control plane. The only configuration that remains at the edge with the
agent are the Tap configuration (ns1labs/pktvisor#75) and edge Tags configuration (below) because they are host
specific.

### Tags and Group Configurations

Orb needs the ability to address the agents that it is controlling. It does this by matching Groups with Tags.

#### Tags

orb-agent is told on startup what its Tags are: these are arbitrary labels which typically represent information such as
region, pop, or node type.

`tags.yaml`

```yaml
version: "1.0"

orb:
  tags:
    region: EU
    pop: ams02
    node_type: dns
```

#### tags on orb-agent start up

```shell
$ orb-agent --config tags.yaml
```

#### combining tags and taps on orb-agent start up

Since both Taps and Tags are necessary for orb-agent start up, you can pass both in via two separate config files:

```shell
$ orb-agent --config taps.yaml --config tags.yaml
```

Or instead combine them into a single file:

`orb-agent.yaml`

```yaml
version: "1.0"

pktvisor:
  taps:
    anycast:
      type: pcap
      config:
        iface: eth0
orb:
  tags:
    region: EU
    pop: ams02
    node_type: dns
```

```shell
$ orb-agent --config orb-agent.yaml
```

### Orb Groups

Groups are named configurations of arbitrary labels which can match against the Tags of the agents available in the Orb
ecosystem. They may be thought of as groups of agents. These names are referenced in Orb Policies.
pktvisord does not read this configuration or use this data; it is used only by orb-agent. This schema is found only in
the control plane, not on the command line or in files.

```yaml
version: "1.0"

orb:
  groups:
    all_dns:
      node_type: dns
    eu_dns:
      region: EU
      node_type: dns
```

### Orb Sinks

Orb includes a metric collection system. Sinks specify where to send the summarized metric data. pktvisord does not read
this configuration or use this data; it is used only by orb-agent. This schema is found only in the control plane, not
on the command line or in files.

```yaml
version: "1.0"

orb:
  sinks:
    default_prometheus:
      type: prometheus_exporter
      address: 0.0.0.0:9598
      default_namespace: service
    my_s3:
      type: aws_s3
      bucket: my-bucket
      compression: gzip
      region: us-east-1
```

### Orb Policies

An Orb policy ties together a Group, an agent Collection Policy, and one or more Sinks. pktvisord does not read this
configuration or use this data; it is used only by orb-agent. This schema is found only in the control plane, not on the
command line or in files.

orb-agent will be made aware of the collection policy if it is in the policy's group. In case of a match, orb-agent will
attempt to apply the collection policy to its pktvisord, and update the control plane about success or failure.

```yaml
version: "1.0"

orb:
  policy:
    group: eu_dns
    agent_policy: anycast_dns
    sinks:
      - default_prometheus
```

