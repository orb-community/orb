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

    'Duplicate sink creation' : (browser) => {
    const sinks = browser.page.sinks();

    sinks
    .navigate()
    .sinkManagementPage()
    .verify.visible('.appearance-ghost', "'New Sink' button is visible")
    .verify.attributeEquals('.appearance-ghost','aria-disabled', 'false', "'New Sink' button is enabled")
    .click('.appearance-ghost')
    .sinkCreation('some_name', 'some_description', 'remote_host', 'tester', 'password', 'key', 'value', 'Failed to create Sink')
    .waitForElementVisible('@previous')
    .click('@previous')
    .waitForElementVisible('@back')
    .click('@back')
    .waitForElementVisible('@cancel')
    .click('@cancel')
    .waitForElementVisible('@cancel')
    .click('@cancel')
    .sinkManagementPage()

}}
