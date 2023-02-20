## Scenario: Check if is possible cancel operations with no change 
## Steps:
1 - Create a sink

- REST API Method: POST
- endpoint: /sinks
- header: {authorization:token}

2 - On sinks' page (`orb.live/pages/sinks`) click on edit button

3 - Change sinks' name

4 - Change sink's description and click "next"

5 - Change sink's remote host

6 - Change sink's username

7 - Change sink's password

8 - Click "back" until return to sinks' page

## Expected Result:
- No changes must have been applied to the sink
