module.exports = {
  '@disabled': false,
  
  'layout components should be consistent with spec': (browser) => {
    const forgotPwdLink = browser.launch_url + '/auth/request-password';
    const registerLink = browser.launch_url + '/auth/register';

    const login = browser.page.login();

    login
      .navigate()
      .waitForElementVisible('@orbLogo', 10000, "Orb logo is being displayed")
      .assert.containsText('@orbCaption', 'An Open-Source dynamic edge observability platform', "Message 'An Open-Source dynamic edge observability platform' is being displayed")
      .assert.containsText('@forgotPwdLink', 'Forgot Password?', "'Forgot Password' option is being displayed")
      .assert.attributeEquals('@forgotPwdLink', 'href', forgotPwdLink, "'Forgot Password' is clickable")
      .assert.containsText('@registerLink', 'Register', "'Register' option is being displayed")
      .assert.attributeEquals('@registerLink', 'href', registerLink, "'Register' is clickable")
      .assert.not.enabled('@loginButton', "Login button is not enabled");
  }
}
