## Scenario: Request registration of an unregistered account with invalid password and valid email 
## Steps: 

1 - Request an account registration using a valid email and password with length less than 8


- REST API Method: POST
- endpoint: /users
- body: `{"email":"email", "password":"invalid_password"}`

## Expected Result: 

- The request must fail with bad request (error 400) and response message must be "password does not meet the requirements"
- No account must be registered