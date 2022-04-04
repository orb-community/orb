## Scenario: Request password with unregistered email address 
## Steps: 

1- On Orb auth page (`http://localhost/auth/login`) click in **"Forgot Password?"**

2- On Orb request password page (`https://orb.live/auth/request-password`) insert a unregistered email on
"Email address" field

3- Click on **"REQUEST PASSWORD"** button

## Expected Result:

- UI must inform that an error has occurred  
- No email must be sent
