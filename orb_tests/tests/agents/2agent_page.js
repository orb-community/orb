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
    

    'agents page info': function(browser) {
   
    const agents = browser.page.agents();
  
    agents
    .navigate()
    .AgentsPage()
    },

    'Count of agents': function(browser){
    const agents = browser.page.agents();
    
    agents
    //bug - need to remove this time
    .navigate()
    .pause(1000)
    .getText('.page-count',  function(result){
        //console.log('Value is:', result.value);
        if (result.value == "0 total") {
            browser.expect.elements('datatable-row-wrapper').count.to.equal(parseInt(result.value))
          
            browser.verify.containsText('.sink-info-accent', 'There are no Agents yet.', "Info message of Agents count is correctly displayed")
            browser.verify.containsText('.empty-row', 'No data to display', "View table info message is correctly displayed")
        } else {
            browser.expect.elements('datatable-row-wrapper').count.to.equal(parseInt(result.value))
            browser.verify.containsText('.justify-content-between > div:nth-child(1)', 'You have', "Beginning of info message is correctly displayed")
            browser.verify.containsText('.justify-content-between > div:nth-child(1)', parseInt(result.value), "Number of Agents is correctly displayed")
            browser.verify.containsText('.justify-content-between > div:nth-child(1)', 'agents deployed in', "End of info message is correctly displayed")
        }
    })

       // .pause(2000)
    }
  }
  