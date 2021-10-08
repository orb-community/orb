var registerActions = {
  with: function (userData) {
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
  url: '/auth/register',
  commands: [registerActions],
  elements: {
    form: '.pane form',
    emailInput: 'input[name=email]',
    pwdInput: 'input[name=password]',
    loginButton: 'form button',
    alertMessage: '.alert-message',
  }
}
