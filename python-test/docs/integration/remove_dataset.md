## Scenario: remove dataset from an agent

Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create a sink
4. Create 1 policy
5. Create a dataset linking the group, the sink and the policy
7. Remove the dataset

Expected result:
-
- The agent's heartbeat shows that 0 policies are applied
- Container logs should inform that removed policy was stopped and removed
- Container logs that were output after removing dataset does not contain the message "scraped metrics for policy" referred to deleted policy anymore