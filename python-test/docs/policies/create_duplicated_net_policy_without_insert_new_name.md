## Scenario: Create duplicated net policy without insert new name 

## Steps: 
1 - Create a net policy

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}

2 - Duplicate this policy

- REST API Method: POST
- endpoint: /policies/agent/{policy_id}/duplicate
- header: {authorization:token}


## Expected Result:

- 3 request must have status code 201 (created) and the policy must be created
- From the 4th order, it must fail by conflict
