## Scenario: Insert tags in agents created without tags 
## Steps: 
1 - Create an agent with no pair (key:value) of tags

- REST API Method: POST
- endpoint: /agents
- header: {authorization:token}

2- Edit this agent and insert at least one pair of tag

- REST API Method: PUT
- endpoint: /agents/agent_id
- header: {authorization:token}
 
## Expected Result: 
 
- Request must have status code 200 and tags must be added to the agent