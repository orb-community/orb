module.exports = {
  '@disabled': false,
  
  'layout components should be consistent with spec': (browser) => {
    const forgotPwdLink = browser.launch_url + '/auth/request-password';
    const registerLink = browser.launch_url + '/auth/register';

    const login = browser.page.login();

    login
      .navigate()
      .waitForElementVisible('@orbLogo', 10000)
      .assert.containsText('@orbCaption', 'An Open-Source dynamic edge observability platform')
      .assert.containsText('@forgotPwdLink', 'Forgot Password?')
      .assert.attributeEquals('@forgotPwdLink', 'href', forgotPwdLink)
      .assert.containsText('@registerLink', 'Register')
      .assert.attributeEquals('@registerLink', 'href', registerLink)
      .assert.not.enabled('@loginButton');
  }
}
