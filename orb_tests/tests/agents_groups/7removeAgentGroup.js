
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
    
    'Agent groups delete' : function(browser){

        var agent_groups = browser.page.agent_groups()
        //bug - incomplete reload. need to remove this navigate
        //bug - delay.  need to remove this paude
        agent_groups
        .navigate()
        .pause(1000)
        .countAgentGroups(browser)
        .agentGroupsDelete()
        .listView()
        .agentGroupsDelete()
        .countAgentGroups(browser)
        //bug need to insert loop


    }}
