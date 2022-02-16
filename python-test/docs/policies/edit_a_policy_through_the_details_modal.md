## Scenario: Edit a policy through the details modal 
## Steps:
1 - Create a policy

- REST API Method: POST
- endpoint: /policies/agent/
- header: {authorization:token}

2 - On policies' page (`orb.live/pages/datasets/policies`) click on details button
3 - Click on "edit" button

## Expected Result:
- User should be redirected to this policy's edit page and should be able to make changes