## Scenario: Remove agent container

Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create a sink
4. Create a policy
5. Create a dataset linking the group, the sink and the policy
6. Stop and remove orb-agent container

Expected result:
-
- The orb-agent container logs contain:
```
{"level":"info","ts":"time","caller":"pktvisor/pktvisor.go:390","msg":"pktvisor stopping"}
{"l/pktvisor.go:253","msg":"pktvisor stdout","log": "Shutting down"}
{"level":"info","ts":"time","caller":"pktvisor/pktvisor.go:253","msg":"pktvisor stdout","log": "policy [policy_name]": "stopping"}
{"level":"info","ts":"time","caller":"pktvisor/pktvisor.go:253","msg":"pktvisor stdout","log": "policy [policy_name]": "stopping input instance: "}
{"level":"info","ts":"time","caller":"pktvisor/pktvisor.go:253","msg":"pktvisor stdout","log": "policy [policy_name]": "stopping handler instance: "}
{"level":"info","ts":"time","caller":"pktvisor/pktvisor.go:253","msg":"pktvisor stdout","log": "exit with success"}
```
- Logs should not have any error
