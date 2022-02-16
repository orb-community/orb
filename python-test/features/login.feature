@login
  Feature: login tests

  Scenario Outline: Request registration of a registered email using <password_status> password
      Given there is a registered account
      When request referred account registration using registered email, <password_status> password, <username> user name and <company> company name
      Examples:
        |  password_status  | username | company |
        |      registered   |  Tester  |    NS1  |
        |      registered   |   None   |    None |
        |      registered   |  Tester  |    None |
        |      registered   |   None   |    NS1  |
        |    unregistered   |  Tester  |    NS1  |
        |    unregistered   |   None   |    None |
        |    unregistered   |  Tester  |    None |
        |    unregistered   |   None   |    NS1  |
      Then account register should not be changed

    Scenario Outline: Login with invalid credentials
      Given there is a registered account
      When the Orb user request an authentication token using <email_status> email and <password_status> password
      Examples:
        | email_status | password_status |
        | incorrect   | incorrect         |
        | incorrect   | correct         |
        | correct   | incorrect         |
      Then user should not be able to authenticate

  Scenario Outline: Check if email is a required field
      When user request account registration <email> email, <password> password, <username> user name and <company> company name
      Examples:
        |    email   | password | username | company |
        |   without  |    with  |   with   |   with  |
        |   without  |    with  |   with   | without |
        |   without  |    with  | without  | without |
        |   without  |    with  | without  |   with  |
        |   without  | without  | without  |   with  |
        |   without  | without  | without  | without |
        |   without  | without  |   with   |   with  |
        |   without  | without  |   with   | without |
      Then user should not be able to authenticate


  Scenario Outline: Check if password is a required field
      When user request account registration <email> email, <password> password, <username> user name and <company> company name
      Examples:
        |  password  |   email  | username | company |
        |   without  |   with   |   with   |   with  |
        |   without  |   with   |   with   | without |
        |   without  |   with   | without  | without |
        |   without  |   with   | without  |   with  |
        |   without  | without  | without  |   with  |
        |   without  | without  | without  | without |
        |   without  | without  |   with   |   with  |
        |   without  | without  |   with   | without |
      Then user should not be able to authenticate
