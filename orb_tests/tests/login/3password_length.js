module.exports = {
  '@disabled': false,

  'Login with an invalid password legth ': (browser) => {
    const login = browser.page.login();

    login
      .with('wrong@email.com', '1234567')
      .click('@emailInput')
      .waitForElementVisible('@requiredAlert', 10000, "Help message about password length is visible")
      .assert.containsText('@requiredAlert', 'Password should contain from 8 to 50 characters', "Help messagem about password lenght contains text 'Password should contain from 8 to 50 characters'")
      .assert.not.enabled('@loginButton', "Login button is not enabled");
  }
}
