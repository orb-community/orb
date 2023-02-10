<img src="docs/images/ORB-logo-black@3x.png" alt="Orb" width="500"/>
<img src="https://user-images.githubusercontent.com/97463920/218170067-16a95078-6709-4828-b137-9791376b972e.png" alt="Orb UI Preview" width="500"/>


[![Go Report Card](https://goreportcard.com/badge/github.com/ns1labs/orb)](https://goreportcard.com/report/github.com/ns1labs/orb)
[![CodeCov](https://codecov.io/gh/ns1labs/orb/branch/develop/graph/badge.svg)](https://app.codecov.io/gh/ns1labs/orb/tree/develop)

**Orb** is a modern network observability platform built to provide critical visibility into increasingly complex and distributed networks. It can analyze network traffic, run synthetic network probes, and connect the resulting telemetry directly to your existing observability stacks with OpenTelemetry. Orb differentiates from other solutions by pushing analysis close to the traffic sources (reducing inactionable metrics and processing costs), and allows for dynamic reconfiguration of remote agents in real time.

Ready to dive in? See [orb.community](https://orb.community) for [installation instructions](https://orb.community/documentation/install/).

# Why Orb?

## Distributed Deep Network Observability

Orb manages a [fleet](https://orb.community/about/#fleet) of [agents](https://orb.community/about/#agent) deployed across
distributed, hybrid infrastructure:
containers, VMs, servers, routers and switches. The agent taps into traffic streams and extracts real time insights,
resulting in light-weight, actionable metrics.

## Streaming Analysis at the Edge

Based on the [pktvisor observability agent](https://pktvisor.dev), Orb's goal is to push analysis to the edge, where
high resolution data can be analysed in real time without the need to send raw data to a central location for batch
processing.
[Current analysis](https://github.com/orb-community/pktvisor/wiki/Current-Metrics) focuses on L2-L3 Network, DNS, and DHCP
with more analyzers in the works.

## Realtime Agent Orchestration

Orb uses IoT principals to allow the observability agents to connect out to the Orb central control plane, avoiding
firewall problems. Once connected, agents are controlled in real time from the Orb Portal or REST API, orchestrating
observability [policies](https://orb.community/about/#policies) designed to precisely extract the desired insights. Agents
are grouped and addressed based on [tags](https://orb.community/about/#agent-group).

## Flexible Integration With Modern Observability Stacks

Orb was built to integrate with modern observability stacks, supporting [Prometheus](https://prometheus.io/) natively
and designed to support arbitrary [sinks](https://orb.community/about/#sinks) in the future. Collection and sinking of the
metrics from the agents is included; there is no need to run additional data collection pipelines for Orb metrics.

## Portal and REST API Included

Orb includes a modern, responsive UI for managing Agents, Agent Groups, Policies and Sinks. Orb is API first, and all
platform functionality is available for automation via
the [well documented REST API](https://orb.community/api/orb_rest_api/).

## Open Source, Vendor Neutral, Cloud Native

Orb is free, open source software (FOSS) released under MPL. It's a modern microservices application that can be
deployed to any Kubernetes service in private or public cloud. It does not depend on any one vendor to function, thus
avoiding vendor lock-in.

***

# Backed by NS1

**Orb** was born at [NS1 Labs](https://ns1.com/labs), where we're committed to
making [open source, dynamic edge observability a reality](https://ns1.com/blog/orb-a-new-paradigm-for-dynamic-edge-observability)
.

***

* [Installation Instructions](https://orb.community/documentation/install/)
* [View our Wiki](https://github.com/orb-community/orb/wiki) for technical and architectural information
* [File an issue](https://github.com/orb-community/orb/issues/new)
* Follow our [public work board](https://github.com/orb-community/orb/projects/2)
* Start a [Discussion](https://github.com/orb-community/orb/discussions)
* [Join us on Slack](https://join.slack.com/t/netdev-community/shared_invite/zt-1ovv03gwn-S30LtO1qQFvHuvfsEZfjvg)
* Send mail to [info@pktvisor.dev](mailto:info@pktvisor.dev)
