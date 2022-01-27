## Scenario: Request registration of an unregistered account with valid password and valid email 
## Steps: 

1 - Request an account registration using a valid email and valid password


- REST API Method: POST
- endpoint: /users
- body: `{"email":"email", "password":"password"}`

## Expected Result:

- The request must be processed successfully (status code 201)
- The new account must be registered
- User must be able to access orb using email and password registered