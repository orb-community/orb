## Scenario: Remove agent

Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create a sink
4. Create a policy
5. Create a dataset linking the group, the sink and the policy
6. Remove agent from orb

Expected result:
-
- Orb-agent logs should not have any error
- Group must match 0 agents