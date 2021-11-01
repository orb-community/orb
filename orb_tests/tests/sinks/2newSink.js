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
      
    'Sink Creation' : (browser) => {
    const sinks = browser.page.sinks();

    sinks
    .navigate()
    //.countSinks()
    .verify.visible('.appearance-ghost')
    .verify.attributeEquals('.appearance-ghost','aria-disabled', 'false')
    .click('.appearance-ghost')
    .sinkCreation('some_name', 'some_description', 'remote_host', 'tester', 'password', 'key', 'value', 'Sink successfully created')
    .sinkManagementPage()
    //bug need to remove this pause
    .pause(2000)
    .countSinks(browser)
    	

}}
