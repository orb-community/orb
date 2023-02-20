## Scenario: Check if remote host, username and password are required to create a sink 

---------------------------------------------------------

## Without remote host


## Steps: 
1 - Create a sink without remote host

- REST API Method: POST
- endpoint: /sinks
- header: {authorization:token}
 
## Expected Result: 

- Request must fail with status code 400 (bad request) 

--------------------------------------------------------

## Without username


## Steps:
1 - Create a sink without username

- REST API Method: POST
- endpoint: /sinks
- header: {authorization:token}

## Expected Result:

- Request must fail with status code 400 (bad request) 

--------------------------------------------------------

## Without password


## Steps:
1 - Create a sink without password

- REST API Method: POST
- endpoint: /sinks
- header: {authorization:token}

## Expected Result:

- Request must fail with status code 400 (bad request) 