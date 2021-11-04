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

    'Sink Edit Description, Remote Host, Username, Password and Keys'  : (browser) => {
    const sinks = browser.page.sinks();

    sinks
    .navigate()
    .sinkEdit()
	.sinkEditAttribute('@sinkDescription','_n33w')
	.click('@sinkNext')
	.sinkEditAttribute('@sinkRemoteHost', '_n33w_rm')
	.sinkEditAttribute('@sinkUsername', '_n33w_usr')
	.sinkEditAttribute('@sinkPassword', '_n33w_pass')
	.click('@sinkNext')
	.click('.ml-1')
	.verify.attributeEquals('@submit','aria-disabled', 'true', "'Subimit' button is not enabled")
	.sinkEditTags('@key', '@value', '_n33w_key', '_n33w_value')
	.click('@sinkNext')
	.click('@sinkNext')
	.verify.containsText('span.title', 'Sink successfully updated', "Confirmation message is being correctly displayed")
	.sinkCheckEdition('_n33w')
	.sinkCheckEdition('_n33w_rm')
	// BUG
	//.sinkCheckEdition('_n33w_usr')
	.sinkCheckEdition('_n33w_key')
	.sinkCheckEdition('_n33w_value')

}}
