## Scenario: remove one of multiple datasets with same policy

Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create a sink
4. Create 1 policy
5. Create a dataset linking the group, the sink and the policy
6. Create another dataset linking the same group and policy (sink can be the same or a different one)
7. Remove one of the datasets

Expected result:
-
- The agent's heartbeat shows that 1 policies are applied
- The orb agent container logs that were output after removing dataset contain the message "scraped metrics for policy" referred to the applied policy