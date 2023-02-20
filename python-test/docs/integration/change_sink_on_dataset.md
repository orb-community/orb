## Scenario: edit sink on dataset

Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create 2 sinks
4. Create 1 policy
5. Create a dataset linking the group, one of the sinks and the policy
6. Wait for scraping metrics for policy
7. Edit the dataset changing the sink

Expected result:
-
- The policy must be applied to the agent (orb-agent API response)
- The container logs contain the message "policy applied successfully" referred to the policy
- The container logs contain the message "scraped metrics for policy" referred to the policy
- Datasets have validity valid
- First applied sink must stop to receive data after the edition
- Second applied sink must start to receive data after the edition