## Scenario: sink has active status while scraping metrics

Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create a sink with valid credentials
4. Create 1 policy
5. Create a dataset linking the group, the sink and one of the policies
6. Wait 1 minute

Expected result:
-
- The container logs contain the message "scraped metrics for policy" referred to the applied policy
- Sink status must be "active"

<br>
-------------------------------------------------

## Advanced Scenario: sink has active status while scraping metrics

Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create a sink
4. Create 1 policy
5. Create a dataset linking the group, the sink and one of the policies
6. Keep sending data for 24 hours


Expected result:
-
- The container logs contain the message "scraped metrics for policy" referred to the applied policy
- Check if even after this time, sink status remains active