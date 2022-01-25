## Scenario: Edit sink remote host 
## Steps:
1 - Create a sink

- REST API Method: POST
- endpoint: /sinks
- header: {authorization:token}

2- Edit this sink remote host

- REST API Method: PUT
- endpoint: /sinks/sink_id
- header: {authorization:token}


## Expected Result:
- Request must have status code 200 (ok) and changes must be applied
 