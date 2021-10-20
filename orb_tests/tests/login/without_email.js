module.exports = {
  '@disabled': false,

  'Login without email': (browser) => {
    const login = browser.page.login();

    login
      .with('', '12345678')
      .click('@pwdInput')
      .waitForElementVisible('@requiredAlert', 10000, "Email request help message is visible")
      .assert.containsText("@requiredAlert", "Email is required!", "Help message is 'Email is required'")
      .assert.not.enabled('@loginButton', "Login button is not enabled");
  }
}
