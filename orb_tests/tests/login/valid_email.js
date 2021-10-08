module.exports = {
  '@disabled': false,

  'email should be valid': (browser) => {
    const login = browser.page.login();

    login
      .with('test@email', '12345678')
      .click('@pwdInput')
      .waitForElementVisible('@requiredAlert', 10000)
      .assert.containsText('@requiredAlert', 'Email should be the real one!')
      .assert.not.enabled('@loginButton');
  }
}
