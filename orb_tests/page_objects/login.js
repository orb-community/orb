var loginActions = {
  with: function (email, pass) {
    return this
      .navigate()
      .waitForElementVisible('@form', 10000)
      .setValue('@emailInput', email)
      .setValue('@pwdInput', pass)
      .click('@loginButton');
  },

  expectAlert: function (message) {
    return this
      .waitForElementVisible('@alertMessage', 10000)
      .assert.containsText('@alertMessage', message);
  },
}

module.exports = {
  url: '/auth/login',
  commands: [loginActions],
  elements: {
    form: '.pane form',
    emailInput: 'input[name=email]',
    pwdInput: 'input[name=password]',
    loginButton: 'form button',
    alertMessage: '.alert-message',
    requiredAlert: 'p.status-danger',
    orbLogo: '#orb-pane-div img',
    orbCaption: '#orb-pane-div p',
    forgotPwdLink: 'a.forgot-password',
    registerLink: '[aria-label="Register"] a'
  }
}