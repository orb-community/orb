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


      'Agent Delete' : function(browser){
        var agents = browser.page.agents()
        agents
        .navigate()
        .AgentsPage()
        .countAgent(browser)
        .agentsDelete()        
        .AgentsPage()
        .countAgent(browser)

      }}
