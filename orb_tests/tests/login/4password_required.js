module.exports = {   
  '@disabled': false,

  'Login without password': (browser) => {
    const login = browser.page.login();

    login
      .with('tester@email.com', '')
      .click('@emailInput')
      .waitForElementVisible('@requiredAlert', 10000, "Password request help message is visible")
      .assert.containsText('@requiredAlert', 'Password is required!', "Help message is 'Password is required!'")
      .assert.not.enabled('@loginButton', "Login button is not enabled");
  }
}
