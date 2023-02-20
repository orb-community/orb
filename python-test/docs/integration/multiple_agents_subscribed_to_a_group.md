## Scenario: multiple agents subscribed to a group

Steps:
-  
1. Provision an agent with tags
2. Provision another agent with same tags (use var env `ORB_BACKENDS_PKTVISOR_API_PORT` to change pktvisor port)
3. Create a group with same tags as agents
4. Create a sink
5. Create 1 policy
6. Create a dataset linking the group, the sink and the policy

Expected result:
-
- The policy must be applied to both agents (orb-agent API response)
- The container logs contain the message "policy applied successfully" referred to the policy
- The container logs contain the message "scraped metrics for policy" referred to the policy
- Referred sink must have active state on response
- Dataset must have validity valid