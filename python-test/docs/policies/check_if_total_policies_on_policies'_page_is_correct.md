## Scenario: Check if total policies on policies' page is correct
## Steps:
1 - Create multiple policies

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}

2 - Get all existing policies

- REST API Method: GET
- endpoint: /policies/agent/

3 - On policies' page (`orb.live/pages/datasets/policies`) check the total number of policies at the end of the policies table

4 - Count the number of existing policies

## Expected Result:
- Total policies on API response, policies page and the real number must be the same 

