## Scenario: Save agent without tag 
## Steps: 
1 - Create an agent with at least one pair (key:value) of tags

- REST API Method: POST
- endpoint: /agents
- header: {authorization:token}

2- Edit this agent tag and remove all pairs

- REST API Method: PUT
- endpoint: /agents/agent_id
- header: {authorization:token}

## Expected Result: 
- Request must have status code 200 and all tags must be removed from the agent 
