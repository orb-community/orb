module.exports = {

    
    'Login with correct email and password' : function(browser) {
        var login = browser.page.login()
       login
          .navigate()
          .assert.titleContains('ORB')
          .assert.visible('#orb-pane-div > p:nth-child(2)')
          .assert.containsText('#orb-pane-div > p:nth-child(2)','An Open-Source dynamic edge observability platform')
          .fillLoginForm()
          .submitLogin()
   
         },
    
    
    'Agent Group Page' : function(browser) {
        var agent_groups = browser.page.agent_groups()
        agent_groups
        .assert.visible('li.menu-item:nth-child(2)')
		.click('li.menu-item:nth-child(2)')
        .waitForElementVisible('xpath', '/html/body/ngx-app/ngx-pages/ngx-one-column-layout/nb-layout/div/div/div/nb-sidebar/div/div/nb-menu/ul/li[2]/ul/li[2]/a')
        .click('xpath','/html/body/ngx-app/ngx-pages/ngx-one-column-layout/nb-layout/div/div/div/nb-sidebar/div/div/nb-menu/ul/li[2]/ul/li[2]/a')
        .agentGroupsPage()
        

        },

    'Agent Group Creation - without description' : function(browser){
        var agent_groups = browser.page.agent_groups()
        agent_groups
        .agentGroupsPage()
        .assert.containsText('.page-count', '0 total')
        .click('@agentGroupCreation')
        .agentGroupCreationPage()
        .agentGroupCreation('name', "", "key", "value", "Agent Group successfully created")
        .assert.containsText('.page-count', '1 total')
    },

    'Agent Group Creation - with description' : function(browser){
        var agent_groups = browser.page.agent_groups()
        agent_groups
        .agentGroupsPage()
        .assert.containsText('.page-count', '1 total')
        .click('@agentGroupCreation')
        .agentGroupCreationPage()
        .agentGroupCreation('nam3', "some_description", "key", "value", "Agent Group successfully created")
        .agentGroupsPage()
        .assert.containsText('.page-count', '2 total')
    },


    
}
