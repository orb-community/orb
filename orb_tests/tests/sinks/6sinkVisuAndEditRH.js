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


	  
    'Sink Visualization and Edit Remote Host': (browser) => {
    const sinks = browser.page.sinks();

    sinks
    .navigate()
    .sinkVisualization()
	.click('.nb-close')
	.sinkVisualization()
	.click('.sink-edit-button')
	.sinkEditPage()
	.click('@sinkNext')
	.sinkEditAttribute('@sinkRemoteHost', '_new_rm')
	.click('@sinkNext')
	.click('@sinkNext')
	.click('@sinkNext')
	.verify.containsText('span.title', 'Sink successfully updated')
	.sinkCheckEdition('_new_rm')

}}
