module.exports = {   
  '@disabled': false,

  'password is required': (browser) => {
    const login = browser.page.login();

    login
      .with('wrong@email.com', '')
      .click('@emailInput')
      .waitForElementVisible('@requiredAlert', 10000)
      .assert.containsText('@requiredAlert', 'Password is required!')
      .assert.not.enabled('@loginButton');
  }
}