## Scenario: Check if total sinks on sinks' page is correct 
## Steps:
1 - Create multiple sinks

- REST API Method: POST
- endpoint: /sinks
- header: {authorization:token}

2 - Get all existing sinks

- REST API Method: GET
- endpoint: /sinks

3 - On sinks' page (`orb.live/pages/sinks`) check the total number of sinks at the end of the sinks table

4 - Count the number of existing sinks

## Expected Result:
- Total sinks on API response, sinks page and the real number must be the same 

 
