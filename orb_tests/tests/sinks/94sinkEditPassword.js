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

    'Sink Edit Password': (browser) => {
    const sinks = browser.page.sinks();


    sinks
    .navigate()
    .sinkEdit()
	.click('@sinkNext')
	.sinkEditAttribute('@sinkPassword', '_n3w_pass')
	.click('@sinkNext')
	.click('@sinkNext')
	.click('@sinkNext')
	.verify.containsText('span.title', 'Sink successfully updated')

}}
