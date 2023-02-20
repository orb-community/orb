## Scenario: Request registration of a registered account using unregistered password, username and company 
## Steps: 

1 - Request an account registration using an already registered email, unregistered password and fullname and company field filled


- REST API Method: POST
- endpoint: /users
- body: `{"email":"already_registered_email", "password":"unregistered_password", "metadata":{"company":"company","fullName":"name"}}`

## Expected Result:
 
- The request must fail with conflict (error 409), response message must be "email already taken"
- No changes should be made to the previously registered account
  (name, company and password must be the ones registered for the first time and the new password should not give access to the account)
