## Scenario: Remove sink using correct name 
## Steps:
1 - Create a sink

- REST API Method: POST
- endpoint: /sinks
- header: {authorization:token}

2 - On agent groups' page (`orb.live/pages/sinks`) click on remove button
3 - Insert the name of the sink correctly on delete modal
4 - Confirm the operation by clicking on "I UNDERSTAND, DELETE THIS SINK" button

## Expected Result:
- Sink must be deleted 

