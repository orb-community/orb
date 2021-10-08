module.exports = {
  '@disabled': false,

  'wrong combination of login and password': (browser) => {
    const login = browser.page.login();

    login
      .with('wrong@email.com', '12345678')
      .expectAlert('Login/Email combination is not correct, please try again.');
  }
}
