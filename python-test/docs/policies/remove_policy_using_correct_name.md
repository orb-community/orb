## Scenario: Remove policy using correct name 
## Steps:
1 - Create a policy

- REST API Method: POST
- endpoint:  /policies/agent/
- header: {authorization:token}

2 - On policies' page (`orb.live/pages/datasets/policies`) click on remove button
3 - Insert the name of the policy correctly on delete modal
4 - Confirm the operation by clicking on "I UNDERSTAND, DELETE THIS POLICY" button

## Expected Result:
- Policy must be deleted 

 