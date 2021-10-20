module.exports = {
  beforeEach: (browser) => {
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

  'New Agent Group' : function(browser) {
    const agentGroup = browser.page.agent_groups();
    const data = {
      step1: {
        name: 'group1',
        description: 'group1 description',
      },
      step2: [
        {
          key: 'region',
          value: 'br',
        },
        {
          key: 'backend',
          value: 'visor',
        },
        {
          key: 'pop',
          value: 'pop03',
        }
      ]
    };

    agentGroup.new()
      .waitForElementVisible('@newHeading', 10000)
      .verify.containsText('@activeBreadcrumb', 'New Agent Group')
      .verify.containsText('@newHeading', 'New Agent Group')
      .verify.containsText({ selector: '@stepLabel', index: 0 }, 'Agent Group Details')
      .verify.containsText({ selector: '@stepCaption', index: 0 }, 'This is how you will be able to easily identify your Agent Group')
      .verify.containsText({ selector: '@stepLabel', index: 1 }, 'Agent Group Tags')
      .verify.containsText({ selector: '@stepCaption', index: 1 }, 'Set the tags that will be used to group Agents')
      .verify.containsText({ selector: '@stepLabel', index: 2 }, 'Review & Confirm')
      .verify.not.enabled('@next')
      .verify.enabled('@back')
      .verify.containsText({ selector: '@detailsLabels', index: 0 }, 'Agent Group Name*')
      .setValue('@newNameInput', data.step1.name)
      .verify.enabled('@next')
      .verify.containsText({ selector: '@detailsLabels', index: 1 }, 'Agent Group Description')
      .setValue('@newDescriptionInput', data.step1.description)
      .click('@next')
      
      .verify.containsText({selector: '@tagLabels', index: 0}, 'Key*')
      .verify.containsText({selector: '@tagLabels', index: 2}, 'Value*')
      .verify.not.enabled('@next')
      .verify.enabled('button.status-primary:nth-child(1)')
    
    data.step2.forEach((tag, i) => {
      agentGroup.setValue('@keyInput', tag.key)
      .setValue('@valueInput', tag.value)
      .click('@addTagButton')
      .waitForElementVisible({ selector: '@tagChip', index: 0 }, 10000)
      .verify.containsText({ selector: '@tagChip', index: 0 }, `${tag.key}: ${tag.value}`)
    });
    
    agentGroup.verify.enabled('@next')
    .click('@next')


  }
};
