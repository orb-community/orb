module.exports = {
  '@disabled': false,

  'password should contain from 8 to 50 characters': (browser) => {
    const login = browser.page.login();

    login
      .with('wrong@email.com', '1234567')
      .click('@emailInput')
      .waitForElementVisible('@requiredAlert', 10000)
      .assert.containsText('@requiredAlert', 'Password should contain from 8 to 50 characters')
      .assert.not.enabled('@loginButton');
  }
}