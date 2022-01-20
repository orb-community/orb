## Scenario: remove one of multiple datasets

Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create a sink
4. Create 2 policies
5. Create a dataset linking the group, the sink and one of the policies
6. Create another dataset linking the same group, sink and the other policy
7. Remove one of the datasets

Expected result:
-
- The agent's heartbeat shows that 1 policies are applied
- The orb-agent container logs should inform that removed policy was stopped and removed
- The orb-agent container logs that were output after removing dataset contain the message "scraped metrics for policy" referred to applied policy
- The orb-agent container logs that were output after removing dataset does not contain the message "scraped metrics for policy" referred to deleted policy anymore