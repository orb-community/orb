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

'Agent Group Edit': function(browser){
    var agent_groups = browser.page.agent_groups()
    agent_groups
    .navigate()

    .agentGroupVisualization()
    .click('@edit')
    .click('@back')
    .listView()
    .agentGroupVisualization()
    .click('@edit')
    .agentGroupsEdit('agentGroup_Edited', 'edited_desc', 'editedKey', 'editedValue', 'secondKey', 'secondValue', 'Agent Group successfully updated')
    .demo_Test()
    .agentGroupCheck('agentGroup_Edited', 'edited_desc', 'editedKey', 'editedValue')
    .demo_Test()
    .agentGroupCheck('agentGroup_Edited', 'edited_desc', 'secondKey', 'secondValue')
}}
