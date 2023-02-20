## Scenario: Check agent details 
## Steps: 
1 - Create an agent

- REST API Method: POST
- endpoint: /agents
- header: {authorization:token}

2 - Get an agent

- REST API Method: GET
- endpoint: /agents/agent_id

## Expected Result: 
- Status code must be 200 and an agent name, channel id, ts_created, status and tags must be returned on response
  * If an agent container was never provisioned, status must be `new`
  * If an agent container is running, status must be `online`
  * If an agent container is stopped/removed, status must be `offline`
