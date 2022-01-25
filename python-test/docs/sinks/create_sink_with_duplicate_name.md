## Scenario: Create sink with duplicate name 
## Steps:
1 - Create a sink

- REST API Method: POST
- endpoint: /sinks
- header: {authorization:token}

2 - Create another sink using the same sink name

## Expected Result:
- First request must have status code 201 (created) and one sink must be created on orb
- Second request must fail with status code 409 (conflict) and no other sink must be created (make sure that first sink has not been modified)

