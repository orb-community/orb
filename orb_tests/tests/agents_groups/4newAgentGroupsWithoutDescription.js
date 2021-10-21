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

'Agent Group Creation - without description' : function(browser){
    var agent_groups = browser.page.agent_groups()
    agent_groups
    .navigate()
    .listView()
    .new()
    .agentGroupCreationPage()
    .agentGroupCreation('name', "", "key", "value", "Agent Group successfully created")
    //bug - incomplete reload. need to remove this navigate
    //bug - delay.  need to remove this paude
    .navigate()
    .pause(1000)
    .countAgentGroups(browser)
}}
