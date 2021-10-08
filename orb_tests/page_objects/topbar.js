var userActions = {
  expectLoggedUser: function (email) {
    return this
      .waitForElementVisible('@userInfo', 10000)
      .assert.containsText('@userInfo', email);
  }
}

module.exports = {
  commands: [userActions],
  elements: {
    userInfo: '.user-name',
  }
}