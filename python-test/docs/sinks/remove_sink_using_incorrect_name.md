## Scenario: Remove sink using incorrect name 
## Steps:
1 - Create a sink

- REST API Method: POST
- endpoint: /sinks
- header: {authorization:token}

2 - On agent groups' page (`orb.live/pages/sinks`) click on remove button
3 - Insert the name of the sink correctly on delete modal

## Expected Result:
- Sink must be deleted 
- "I UNDERSTAND, DELETE THIS SINK" button must not be enabled
- After user close the deletion modal, sink must not be deleted  