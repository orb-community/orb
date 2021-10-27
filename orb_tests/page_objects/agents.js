var loginActions = {
    AgentsPage: function () {
        return this
        .waitForElementVisible('@path', "Agents path is visible")
        .verify.containsText('@agentPath', 'Fleet Management', "Agents is inherited from Fleet Management")
        .verify.containsText('@view', 'Agents List', "Agent view is named 'Agents List'")
        .verify.containsText('@header', "All Agents", "Agents Header is 'All Agents'")
        .waitForElementVisible('.flex-column', "Agent Groups table view is visible")
        .waitForElementVisible('@table', "Agent table view is visible")
        .waitForElementVisible("@new", "New Agent button is visible")
        .waitForElementVisible("@filter", "Filter type is visible")
        .waitForElementVisible("@search", "Search by filter is visible")

    },

    new: function() {
      return this.navigate()
      .waitForElementVisible('@newButton', 10000,  "New Agent button is visible")
      .click('@newButton')
    },
    

    agentCreationPage: function() {
      return this.waitForElementVisible('@pathNew', "Agents path is visible")
      .verify.containsText('@agentPath', 'Fleet Management', "Agents is inherited from Fleet Management")
      .verify.containsText('@pathNew', 'Agents List', "Agent view is named 'Agents List'")
      .verify.containsText('@pathNew', "New Agent", "Agents Header is 'New Agent'")
      .verify.containsText('@headerNew', 'New Agent', "Header is 'New Agent'")
      .verify.containsText('@agentDetails', 'Agent Details', "Header is 'Agent Details'")
      .verify.containsText('.header', 'This is how you will be able to easily identify your Agent', "Help text about name is correctly written")
      .verify.containsText('.header', 'Orb Tags', "'Agent Tags' is being displayed")
      .verify.containsText('.header', 'Set the tags that will be used to filter your Agent', "Help text about tags is correctly written")
      .verify.containsText('.header', 'Review & Confirm', "'Review & Confirm' is being displayed")
      .verify.containsText('.step-content', 'Agent Name*', "'Agent Name*' field is being displayed")
      .verify.attributeEquals('@next','aria-disabled', 'true', "'Next' button is not enabled")
      .verify.attributeEquals('@back','aria-disabled', 'false', "'Back' button is enabled")
  },


    agentCreation: function(name, key, value, verify, closeOption='@close') {
      return this.setValue('@newNameInput', name)
      .verify.attributeEquals('@next','aria-disabled', 'false', "'Next' button is enabled")
      // .setValue('@newDescriptionInput', description)
      .click('@next')
      .verify.attributeEquals('@next','aria-disabled', 'true', "'Next' button is not enabled")
      .verify.attributeEquals('button.status-primary:nth-child(1)','aria-disabled', 'false', "'Back' button is enabled")
      .verify.attributeEquals('@addTag','aria-disabled', 'true', "'Add tags' button is not enabled")
      .setValue('@key', key)
      .setValue('@value', value)
      .verify.attributeEquals('@addTag','aria-disabled', 'false', "'Add tags' button is enabled")
      .click('@addTag')
      .verify.attributeEquals('@next','aria-disabled', 'false', "'Next' button is enabled")
      .verify.attributeEquals('button.status-primary:nth-child(1)','aria-disabled', 'false', "'Back' button is enabled")
      .click('@next')
      .verify.containsText('@pathNew', "New Agent", "Agents Header is 'New Agent'")
      .verify.attributeEquals('@back','aria-disabled', 'false', "'Back' button is enabled")
      .verify.attributeEquals('@next','aria-disabled', 'false', "'Next' button is enabled")
      .click('@next')
      .verify.visible('[class="cdk-overlay-pane"]', "Agent Credentials modal is visible")
      .verify.containsText("@agentCredentialsHeader", "Agent Credentials", "Agent Credentials modal's header is 'Agent Credentials'")
      .verify.containsText("@agentCredencialsBody", "Make sure to copy the Agent Key now. You won’t be able to see it again!", "Agent Credentials help text is 'Make sure to copy the Agent Key now. You won’t be able to see it again!'")
      .verify.containsText("@agentCredencialsBody", "Agent Key", "Agent key field's name is 'Agent Key'")
      .verify.containsText("@agentCredencialsBody", "Provisioning Command", "Provisioning Command field's name is 'Provisioning Command'")
      .verify.visible("@agentKey", "Agent Key is being  displayed")
      .verify.visible("@agentProvisioningCommand", "Agent Provisioning Command is being displayed")
      .click('@copyKey')
      .click('@copyProvisioningCommand')
      // options: '@closeCredentialsModal' or '@close'
      .click(closeOption)
      .verify.containsText('span.title', verify, "Confirmation message is correctly displayed")


},

agentsDelete: function() {
  return this.verify.attributeEquals('button.orb-action-hover:nth-child(3)', 'aria-disabled', 'false', "'Remove' agent button is enabled")
  .click('button.orb-action-hover:nth-child(3)')
  .verify.attributeEquals('@deleteAgent','aria-disabled', 'true', "'Confirm agent delete button is not enabled")
  .verify.visible('@deleteAgentModal', "Delete agent modal is visible")
  .verify.containsText('ngx-agent-delete-component.ng-star-inserted > nb-card:nth-child(1) > nb-card-header:nth-child(1)', 'Delete Agent Confirmation', "Header of help text about delete confirmation is correctly written")
  .verify.containsText('@deleteAgentModal', 'Are you sure you want to delete this Agent? This action cannot be undone.*', "Help text about delete confirmation is correctly written")
  .verify.containsText('@deleteAgentModal', '*To confirm, type your Agent name exactly as it appears', "End of help text about delete confirmation is correctly written")
  .getAttribute('.input-full-width', 'placeholder',  function(result) {this.setValue('.input-full-width', result.value) })
  .verify.attributeEquals('@deleteAgent','aria-disabled', 'false', "'Confirm agent delete button is enabled")
  .click('@deleteAgent')
  .verify.containsText('span.title', 'Agent successfully deleted', "Delete confirmation message is being displayed")
  //bug .verify.containsText('.empty-row', 'No data to display')

},

countAgent: function(browser) {
  return this.getText('.page-count',  function(result){
    //console.log('Value is:', result.value);
    if (result.value == "0 total") {
        browser.expect.elements('datatable-row-wrapper').count.to.equal(parseInt(result.value))
      
        browser.verify.containsText('.sink-info-accent', 'There are no Agent yet.', "Info message of Agent count is correctly displayed")
        browser.verify.containsText('.empty-row', 'No data to display', "View table info message is correctly displayed")
    } else {
        browser.expect.elements('datatable-row-wrapper').count.to.equal(parseInt(result.value))
        browser.verify.containsText('.justify-content-between > div:nth-child(1)', 'You have', "Beginning of info message is correctly displayed")
        browser.verify.containsText('.justify-content-between > div:nth-child(1)', parseInt(result.value), "Number of Agents is correctly displayed")
        // bug need to insert regions count test
        browser.verify.containsText('.justify-content-between > div:nth-child(1)', 'agents deployed', "End of info message is correctly displayed")
    }
  })

}
  }
  
  module.exports = {
    url: '/pages/fleet/agents',
    commands: [loginActions],
    elements: {
      path: 'xng-breadcrumb.orb-breadcrumb',
      pathNew: '.xng-breadcrumb-root',
      agentPath: '.xng-breadcrumb-link',
      view: '.xng-breadcrumb-trail',
      headerNew: 'ngx-agent-add-component.ng-star-inserted > div:nth-child(1) > header:nth-child(1) > h4:nth-child(2)',
      header: 'ngx-agent-list-component.ng-star-inserted > div:nth-child(1) > header:nth-child(1) > h4:nth-child(2)',
      table:'.datatable-body',
      new:'.status-primary',
      filter:'.select-button',
      search:'input.size-medium',
      agentsListed: '.datatable-row-wrapper',
      info: '.sink-info-accent',
      emptyRow: '.empty-row',
      countMessage: '.justify-content-between > div:nth-child(1)',
      count:'.page-count',
      newButton: '.appearance-ghost',
      agentDetails: 'div.step:nth-child(1) > div:nth-child(2) > div:nth-child(1) > strong:nth-child(1)',
      next:'.next-button',
      back: '.appearance-ghost',
      close: '.nb-close',
      closeCredentialsModal: 'ngx-agent-key-component.ng-star-inserted > nb-card:nth-child(1) > nb-card-footer:nth-child(3) > button:nth-child(1)',
      newNameInput: '[formcontrolname="name"]',
      key: 'div.col-5:nth-child(1) > div:nth-child(2) > input:nth-child(1)',
      value: 'div.d-flex:nth-child(3) > div:nth-child(2) > input:nth-child(1)',
      addTag: 'button.status-basic',
      agentCredencialsBody: 'ngx-agent-key-component.ng-star-inserted > nb-card:nth-child(1) > nb-card-body:nth-child(2)',
      agentCredentialsHeader: 'ngx-agent-key-component.ng-star-inserted > nb-card:nth-child(1) > nb-card-header:nth-child(1)',
      agentKey: 'ngx-agent-key-component.ng-star-inserted > nb-card:nth-child(1) > nb-card-body:nth-child(2) > pre:nth-child(3)',
      agentProvisioningCommand: 'ngx-agent-key-component.ng-star-inserted > nb-card:nth-child(1) > nb-card-body:nth-child(2) > pre:nth-child(5) > code:nth-child(2)',
      copyKey: 'ngx-agent-key-component.ng-star-inserted > nb-card:nth-child(1) > nb-card-body:nth-child(2) > pre:nth-child(5) > button:nth-child(1)',
      copyProvisioningCommand: 'ngx-agent-key-component.ng-star-inserted > nb-card:nth-child(1) > nb-card-body:nth-child(2) > pre:nth-child(5) > button:nth-child(1)',
      deleteAgent: '.orb-sink-delete-warning-button',
      deleteAgentModal: 'ngx-agent-delete-component.ng-star-inserted'
      


    }
  }
