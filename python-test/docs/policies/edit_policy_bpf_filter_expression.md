## Scenario: Edit policy bpf_filter_expression 

## Steps:
1 - Create a policy

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}

2- Edit this policy bpf_filter_expression

- REST API Method: PUT
- endpoint: /policies/agent/policy_id
- header: {authorization:token}


## Expected Result:
- Request must have status code 200 (ok) and changes must be applied