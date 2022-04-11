## Scenario: apply twice the same policy to agents subscribed to a group

Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create a sink
4. Create 1 policy
5. Create a dataset linking the group, the sink and the policy
6. Create another dataset linking the same group and policy (sink can be the same or a different one)

Expected result:
-
- The policy must be applied to the agent (orb-agent API response) and two datasets must be listed linked to the policy
- The container logs contain the message "policy applied successfully" referred to the policy
- The container logs contain the message "scraped metrics for policy" referred to the policy
- All sinks linked must have active state on response
- Both datasets have validity valid