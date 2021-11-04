module.exports = {
  '@disabled': false,

  'Login with invalid regex for email': (browser) => {
    const login = browser.page.login();

    login
      .with('tester@email', '12345678')
      .click('@pwdInput')
      .waitForElementVisible('@requiredAlert', 10000, "Help message about invalid regex for email is visible")
      .assert.containsText('@requiredAlert', 'Email should be the real one!', "Help message contains text 'Email should be the real one!'")
      .assert.not.enabled('@loginButton', "Login button is not enabled");
  }
}
