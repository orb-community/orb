## Scenario: Edit a sink through the details modal 
## Steps:
1 - Create a sink

- REST API Method: POST
- endpoint: /sinks
- header: {authorization:token}

2 - On sinks' page (`orb.live/pages/sinks`) click on details button
3 - Click on "edit" button

## Expected Result:
- User should be redirected to this sink's edit page and should be able to make changes