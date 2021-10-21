module.exports = {
'@disabled': false,
  
'Create an account' : function(browser) {
	var login = browser.page.login()
    var accountRegister = browser.page.accountRegister()
    const registerLink = browser.launch_url + '/pages';

	const maximizeWindowCallback = () => {
		console.log('Window maximized');
	  };
	  browser.maximizeWindow(maximizeWindowCallback);

	login.navigate()
	accountRegister.orbRegister()
	.verify.urlEquals(registerLink)


	}}
