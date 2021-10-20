module.exports = {
  '@disabled': false,

  'Login with wrong combination of email and password': (browser) => {
    const login = browser.page.login();

    login
      .with('wrong@email.com', 'any-pass')
      .expectAlert('Login/Email combination is not correct, please try again.', "Alert message is visible and contains text 'Login/Email combination is not correct, please try again.'");
  },


  'Login with wrong email and registered password': (browser) => {
    const login = browser.page.login();

    login
      .with('testerr@email.com', '12345678')
      .expectAlert('Login/Email combination is not correct, please try again.', "Alert message is visible and contains text 'Login/Email combination is not correct, please try again.'");
  },


  'Login with registered email and wrong password': (browser) => {
    const login = browser.page.login();

    login
      .with('tester@email.com', '123456789')
      .expectAlert('Login/Email combination is not correct, please try again.', "Alert message is visible and contains text 'Login/Email combination is not correct, please try again.'");
  }
}
