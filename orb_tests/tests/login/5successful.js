module.exports = {
  '@disabled': false,

  'Login with registered email and correct password': function (browser) {
    const login = browser.page.login();
    const topbar = browser.page.topbar();
    const email = 'tester@email.com';
    const pwd = '12345678';

    const maximizeWindowCallback = () => {
      console.log('Window maximized');
    };
    browser.maximizeWindow(maximizeWindowCallback);

    login.with(email, pwd);
    topbar.expectLoggedUser(email);
  }
}
