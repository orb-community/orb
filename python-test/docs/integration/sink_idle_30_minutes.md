## Scenario: sink has idle status after 30 minutes without data

Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create a sink with invalid credentials
4. Create 1 policy
5. Create a dataset linking the group, the sink and one of the policies
6. Wait 1 minute
7. Remove the dataset to which sink is linked
8. Wait 30 minutes

Expected result:
-
- Sink status must be "idle"