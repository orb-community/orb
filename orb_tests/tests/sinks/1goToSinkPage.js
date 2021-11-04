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
      
    'Go to sink page' : (browser) => {
    const sinks = browser.page.sinks();

    sinks
    .assert.visible('li.menu-item:nth-child(4)', "Sink Management is visible on ORB menu")
    .click('li.menu-item:nth-child(4)')
    .sinkManagementPage()

    	

}}
