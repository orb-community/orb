## Scenario: Check if email and password are required fields 

### Sub Scenarios: 
## I - Check if email is a required field

### Steps: 

1 - Request an account registration without email field

- REST API Method: POST
- endpoint: /users
- body: `{"password":"password", "metadata":{"company":"company","fullName":"name"}}`

### Expected Result: 

- The request must fail with bad request (error 400) and no account must be registered

## II - Check if password is a required field

### Steps:

1- Request an account registration without password field

- REST API Method: POST
- endpoint: /users
- body: `{"email":"email", "metadata":{"company":"company","fullName":"name"}}`

### Expected Result:

- The request must fail with bad request (error 400) and no account must be registered


## III - Check if password and email are required fields

### Steps:

1 - Request an account registration using just metadata

- REST API Method: POST
- endpoint: /users
- body: `{"metadata":{"company":"company","fullName":"name"}}`

### Expected Result:

- The request must fail with bad request (error 400) and no account must be registered