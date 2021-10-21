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


	  
    'Sink Visualization and Edit Tags' : (browser) => {
    const sinks = browser.page.sinks();

    sinks
    .navigate()
    .sinkVisualization()
	.click('.nb-close')
	.sinkVisualization()
	.click('.sink-edit-button')
	.sinkEditPage()
	.click('@sinkNext')
	.click('@sinkNext')
	.click('.ml-1')
	.verify.attributeEquals('@submit','aria-disabled', 'true')
	.sinkEditTags('@key', '@value', 'new_key', 'new_value')
	.click('@sinkNext')
	.click('@sinkNext')
	.verify.containsText('span.title', 'Sink successfully updated')
	.sinkCheckEdition('new_key')
	.sinkCheckEdition('new_value')

}}
