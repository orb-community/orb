## Scenario: Login with invalid credentials 

### Sub Scenarios:
## I - Login with invalid email

## Steps: 

1- Request authentication token using unregistered email and some registered password

- REST API Method: POST
- endpoint: /tokens
- body: `{"email": "invalid_email", "password": "password"}`


## Expected Result: 
 - The request must fail with forbidden (error 403) and no token must be generated

## II - Login with invalid password

## Steps:

1- Request authentication token using registered email and wrong password

- REST API Method: POST
- endpoint: /tokens
- body: `{"email": "email", "password": "wrong_password"}`


## Expected Result:
- The request must fail with forbidden (error 403) and no token must be generated

## III - Login with invalid email and invalid password

## Steps:

1- Request authentication token using unregistered email and unregistered password

- REST API Method: POST
- endpoint: /tokens
- body: `{"email": "invalid_email", "password": "invalid_password"}`


## Expected Result:

- The request must fail with forbidden (error 403) and no token must be generated