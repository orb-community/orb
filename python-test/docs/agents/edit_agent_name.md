## Scenario: Edit agent name 

## Steps:
1 - Create an agent

- REST API Method: POST
- endpoint: /agents
- header: {authorization:token}

2- Edit this agent name

- REST API Method: PUT
- endpoint: /agents/agent_id
- header: {authorization:token}


## Expected Result:
- Request must have status code 200 (ok) and changes must be applied
