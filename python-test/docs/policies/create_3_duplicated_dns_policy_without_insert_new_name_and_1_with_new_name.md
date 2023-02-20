## Scenario: Create 3 duplicated dns policy without insert new name and 1 with new name 

## Steps:
1 - Create a policy

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}

2 - Duplicate this policy 3 times

- REST API Method: POST
- endpoint: /policies/agent/{policy_id}/duplicate
- header: {authorization:token}

3 - Duplicate this policy 1 more time inserting new name

- REST API Method: POST
- endpoint: /policies/agent/{policy_id}/duplicate
- header: {authorization:token}
- body: {"name": "new name"}

## Expected Result:

- All requests must have status code 201 (created) and the policies must be created
