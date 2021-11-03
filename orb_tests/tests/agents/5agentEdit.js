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
    var agents = browser.page.agents()
    agents
    .navigate()

    .agentsVisualization()
    .click('@edit')
    .click('@back')
    .AgentsPage()
    .agentsVisualization()
    .click('@edit')
    .agentsEdit('agentEdited', 'region', 'br', 'region2', 'usa', 'Agent successfully updated')
    .choose_last_element()
    .agentCheck('agentEdited',  'region: br')
    .choose_last_element()
    .agentCheck('agentEdited',  'region2: usa')
}}
