module.exports = {
  '@disabled': false,

  'successful login': function (browser) {
    const login = browser.page.login();
    const topbar = browser.page.topbar();
    const email = 'imelo@daitan.com';
    const pwd = '12345678';

    const maximizeWindowCallback = () => {
      console.log('Window maximized');
    };
    browser.maximizeWindow(maximizeWindowCallback);

    login.with(email, pwd);
    topbar.expectLoggedUser(email);
  }
}
