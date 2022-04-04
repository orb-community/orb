## Scenario: Request registration of an unregistered account with valid password and invalid email 
## Steps: 

1 - Request an account registration using an email without `@server` and password with length greater than or equal to 8


- REST API Method: POST
- endpoint: /users
- body: `{"email":"invalid_email", "password":"password"}`

## Expected Result:

- The request must fail with bad request (error 400)
- No account must be registered