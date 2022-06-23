# Orb Agent

Agent - Fleet Synchronization Steps

Success Communication Sequence Diagram
```mermaid
  sequenceDiagram
    Agent->>Fleet: subscribe
    Agent->>Fleet: sendCapabilities
    Agent->>+Fleet: sendGroupMembershipReq 
    Fleet-->>-Agent: sendGroupMembership with GroupIds
    Agent-->>Fleet: Heartbeat
    Agent-->>+Fleet: sendAgentPoliciesReq
    Fleet-->>-Agent: agentPolicies with Policies
    Agent-->>Fleet: Heartbeat
```


Fail Communication Sequence Diagram
```mermaid
  sequenceDiagram
    Agent->>Fleet: subscribe
    Agent->>Fleet: sendCapabilities
    Agent->>+Fleet: sendGroupMembershipReq 
    Agent-->>Fleet: Heartbeat
    Agent-->>+Fleet: sendAgentPoliciesReq
    Fleet-->>-Agent: agentPolicies with Policies
    Agent-->>Fleet: Heartbeat
```

Agent is still without Groups and Policies


With Re-Request Mechanism, general idea

```mermaid
sequenceDiagram
    Agent-)Fleet: subscribe
    Agent-)Fleet: sendCapabilities
    Agent-)Fleet: sendGroupMembershipReq
    activate Fleet
    Agent->>Timer: starts wait timer for response
    activate Timer
    Timer-xAgent: timer runs out
    Agent-)Fleet: sendGroupMembershipReq
    Fleet--)Agent: sendGroupMembership with GroupIds
    deactivate Fleet
    Agent-->>Timer: marks as success
    deactivate Timer
    Agent-->>Fleet: Heartbeat
    Agent-->>Fleet: sendAgentPoliciesReq
    activate Fleet
    Agent->>Timer: starts wait timer for response
    activate Timer
    Fleet-->>Agent: agentPolicies with Policies
    Agent->>Timer: marks as success
    deactivate Timer
    Agent-->>Fleet: Heartbeat
```
