## Scenario: Agent subscription to group after editing agent's tags (agent provisioned before group)
Steps:
-  
1. Provision an agent with tags
2. Create a group with same tags as agent
3. Create another group with different tags
4. Edit agent orb tags to match with second group

Expected result:
-
- Agent heartbeat must show just one group matching
- Agent logs must show that agent is unsubscribed from the first group
- Agent logs must show that agent is subscribed to the second group


## Scenario: Agent subscription to group after editing agent's tags (agent provisioned after group)
Steps:
-  
1. Create a group with tags
2. Provision an agent with same tags
3. Create another group with different tags
4. Edit agent orb tags to match with second group

Expected result:
-
- Agent heartbeat must show just one group matching
- Agent logs must show that agent is unsubscribed from the first group
- Agent logs must show that agent is subscribed to the second group


## Scenario: Agent subscription to group after editing agent's tags (agent provisioned after groups)
Steps:
-  
1. Create a group with tags
2. Create another group with different tags
3. Provision an agent with same tags as first group
4. Edit agent orb tags to match with second group

Expected result:
-
- Agent heartbeat must show just one group matching
- Agent logs must show that agent is unsubscribed from the first group
- Agent logs must show that agent is subscribed to the second group