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
    

'Go to Agent page from Orb Home': function(browser){
    const agentsPage = browser.launch_url + '/pages/fleet/agents'
    browser
    .verify.visible('li.menu-item:nth-child(2)', "Fleet Management is visible on ORB menu")
    .verify.containsText('[title="Fleet Management"]', 'Fleet Management', "Fleet management is correctly writen")
    .click('li.menu-item:nth-child(2)')
    .waitForElementVisible('xpath', '/html/body/ngx-app/ngx-pages/ngx-one-column-layout/nb-layout/div/div/div/nb-sidebar/div/div/nb-menu/ul/li[2]/ul/li[2]/a', "Agents is visible on ORB menu")
    .click('xpath','/html/body/ngx-app/ngx-pages/ngx-one-column-layout/nb-layout/div/div/div/nb-sidebar/div/div/nb-menu/ul/li[2]/ul/li[1]/a')
    .verify.urlEquals(agentsPage)
    browser.page.agents().AgentsPage()
  }}
