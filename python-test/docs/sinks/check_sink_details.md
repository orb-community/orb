## Scenario: Check sink details 
## Steps:
1 - Create a sink

- REST API Method: POST
- endpoint: /sinks
- header: {authorization:token}

2 - Get a sink

- REST API Method: GET
- endpoint: /sinks/sink_id

## Expected Result:
- Status code must be 200 and sink name, description, service type, remote host, status, username and tags must be returned on response

  * If a sink never received data, status must be `new`
  * If a sink is receiving data, status must be `active`
  * If a sink has not received data for more than 30 minutes, status must be `idle`