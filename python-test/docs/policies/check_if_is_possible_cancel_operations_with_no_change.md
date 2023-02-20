## Scenario: Check if is possible cancel operations with no change 
## Steps:
1 - Create a policy

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}

2 - On policies' page (`orb.live/pages/datasets/policies`) click on edit button

3 - Change policy's name

4 - Change policy's description and click "next"

5 - Change policy's tap configuration options and filter and click "next"

6 - Change policy's handler

7 - Click "back" until return to policies' page

## Expected Result:
- No changes must have been applied to the policy
