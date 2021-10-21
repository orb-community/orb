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
    .verify.containsText('.page-count', '1 total')
	.sinkDelete()
	.verify.containsText('.page-count', '0 total')
	.sinkManagementPage()

}}
