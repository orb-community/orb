module.exports = {
  '@disabled': false,

  'email is required': (browser) => {
    const login = browser.page.login();

    login
      .with('', '12345678')
      .click('@pwdInput')
      .waitForElementVisible('@requiredAlert', 10000)
      .assert.containsText('@requiredAlert', 'Email is required!')
      .assert.not.enabled('@loginButton');
  }
}
