## Scenario: Remove policy using incorrect name 
## Steps:
1 - Create a policy

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}

2 - On policies' page (`orb.live/pages/datasets/policies`) click on remove button
3 - Insert the name of the policy incorrectly on delete modal

## Expected Result:
- "I UNDERSTAND, DELETE THIS POLICY" button must not be enabled
- After user close the deletion modal, policy must not be deleted  
