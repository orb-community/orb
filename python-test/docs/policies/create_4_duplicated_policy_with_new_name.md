## Scenario: Create 4 duplicated policy with new name 

## Steps:
1 - Create a policy

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}

2 - Duplicate this policy 4 times inserting new name

- REST API Method: POST
- endpoint: /policies/agent/{policy_id}/duplicate
- header: {authorization:token}
- body: {"name": "new name"}


## Expected Result:

- All requests must have status code 201 (created) and the policies must be created
