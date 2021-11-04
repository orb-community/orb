var userActions = {
  expectLoggedUser: function (email) {
    return this.waitForElementVisible('@userInfo', 10000, "User info is visible")
      //bug on orb
      // .waitForElementVisible('@userInfo', 10000, "User info is visible")
      .assert.containsText('@userInfo', email, "User is displayed on topbar");
  }
}

module.exports = {
  commands: [userActions],
  elements: {
    userInfo: '.user-name',
  }
}
