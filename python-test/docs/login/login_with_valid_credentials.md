## Scenario: Login with valid credentials 

## Steps:

1- Request authentication token using registered email referred password

- REST API Method: POST
- endpoint: /tokens
- body: `{"email": "email", "password": "password"}`


 
## Expected Result: 

- Status code must be 200 and a token must be returned on response