module.exports = {
    '@disabled': false,

    before: (browser) => {
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
      },

      
    'Sink Delete'  : (browser) => {
    const sinks = browser.page.sinks();

    sinks
    .navigate()
    // bug need to remove this pause
    .pause(1000)
    .countSinks(browser)
    .sinkDelete()
    .sinkManagementPage()
    // bug need to remove this pause
    .pause(1000)
    .countSinks(browser)
    .sinkManagementPage()

}}
