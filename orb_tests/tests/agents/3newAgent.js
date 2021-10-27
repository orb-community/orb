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


      'Agent Creation using close button' : function(browser){
        var agent_groups = browser.page.agents()
        agent_groups
        .new()
        .agentCreationPage()
        .agentCreation('newAgent', 'key', 'value', 'Agent successfully created', '@closeCredentialsModal')
    
        
      },

      'Agent Creation using "x" button' : function(browser){
        var agent_groups = browser.page.agents()
        agent_groups
        .new()
        .agentCreationPage()
        .agentCreation('new2Agent', 'key', 'value', 'Agent successfully created', '@close')
    
        
      }
    
    
    }
