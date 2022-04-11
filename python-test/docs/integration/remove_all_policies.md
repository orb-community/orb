## Scenario: remove all policies from an agent

Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create a sink
4. Create 2 policies
5. Create a dataset linking the group, the sink and one of the policies
6. Create another dataset linking the same group, sink and the other policy
7. Remove both policies

Expected result:
-
- The agent's heartbeat shows that 0 policies are applied
- Container logs should inform that removed policy was stopped and removed
- Container logs that were output after removing policies does not contain the message "scraped metrics for policy" referred to deleted policies anymore
- Datasets became "invalid"