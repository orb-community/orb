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

    'Agent Group Creation - with description' : function(browser){
      var agent_groups = browser.page.agent_groups()
      agent_groups
      .navigate()
      .new()
      //bug need to insert count agent groups
      .agentGroupCreation('nam3', "some_description", "key", "value", "Agent Group successfully created")
      //bug - incomplete reload. need to remove this navigate
      //bug - delay.  need to remove this paude
      .navigate()
      .pause(1000)
      .agentGroupVisualization()
      .agentGroupCheck('nam3', "some_description", "key", "value")
      .countAgentGroups(browser)


  }}
