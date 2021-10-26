var commands = {

  new: function() {
    return this.navigate()
    .waitForElementVisible('@newButton', 10000)
    .click('@newButton')
  },

  fillInput: function (selector, data) {
    return this.setValue(selector, data)
  },


  listView: function() {
    return this.verify.containsText('@agentGroupList','Agent Groups List', "Agent Groups view is named 'Agent Groups List'")
    .verify.containsText('@agentGroupAll', 'All Agent Groups', "Agent Groups Header is 'All Agent Groups'")
    .verify.containsText('@agentPath', 'Fleet Management', "Agent Groups is inherited from Fleet Management")
    .verify.visible('@new', "'New Agent Group' button is visible")
    .verify.visible('.flex-column', "Agent Groups table view is visible")
    .verify.visible("@table", "Agent Groups body table view is visible")
    .verify.visible("@filter", "Filter type is visible")
    .verify.visible("@search", "Search by filter is visible")
  },

    agentGroupCreationPage: function() {
        return this.verify.containsText('.header', 'Agent Group Details', "Header is 'Agent Group Details'")
        .verify.containsText('.header', 'This is how you will be able to easily identify your Agent Group', "Help text about name is correctly written")
        .verify.containsText('.header', 'Agent Group Tags', "'Agent Group Tags' is being displayed")
        .verify.containsText('.header', 'Set the tags that will be used to group Agents', "Help text about tags is correctly written")
        .verify.containsText('.header', 'Review & Confirm', "'Review & Confirm' is being displayed")
        .verify.containsText('.step-content', 'Agent Group Name*', "'Agent Group Name*' field is being displayed")
        .verify.containsText('.step-content', 'Agent Group Description',  "'Agent Group Description' field is being displayed")
        .verify.attributeEquals('@next','aria-disabled', 'true', "'Next' button is not enabled")
        .verify.attributeEquals('@back','aria-disabled', 'false', "'Back' button is not enabled")
    },

    agentGroupEditPage: function() {
      return this.verify.containsText('.xng-breadcrumb-trail', 'Edit Agent Group', "Header is 'Edit Agent Group'")
      .verify.containsText('ngx-sink-add-component.ng-star-inserted > div:nth-child(1) > header:nth-child(1)', 'Edit Agent Groups', "Page description is 'Agent Groups'")
      .verify.containsText('@editSinkHeader', 'Agent Group Details', "'Agent Group Details' is being displayed")
      .verify.containsText('@editSinkHeader', 'This is how you will be able to easily identify your Agent Group', "Help text about name is correctly written")
      .verify.containsText('@editSinkHeader', 'Agent Group Tags', "'Agent Group Tags' is being displayed")
      .verify.containsText('@editSinkHeader', 'Set the tags that will be used to group Agents', "Help text about tags is correctly written")
      .verify.containsText('@editSinkHeader', 'Review & Confirm', "'Review & Confirm' is being displayed")
      .verify.attributeEquals('@sinkNext','aria-disabled', 'false', "Next button is enabled")
      .verify.attributeEquals('@back','aria-disabled', 'false', "Cancel button is enabled")
  },

    agentGroupCreation: function(name, description, key, value, verify) {
        return this.setValue('@newNameInput', name)
        .verify.attributeEquals('@next','aria-disabled', 'false', "'Next' button is enabled")
        .setValue('@newDescriptionInput', description)
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
        .verify.containsText('@agentGroupList','New Agent Group', 'Page header is "New Agent Group"')
        .verify.attributeEquals('@back','aria-disabled', 'false', "'Back' button is enabled")
        .verify.attributeEquals('@next','aria-disabled', 'false', "'Next' button is enabled")
        .click('@next')
        .verify.containsText('span.title', verify, "Confirmation message is correctly displayed")


    },

    agentGroupsDelete: function() {
        return this.verify.attributeEquals('button.orb-action-hover:nth-child(3)', 'aria-disabled', 'false', "'Remove' agent group button is enabled")
        .click('button.orb-action-hover:nth-child(3)')
        .verify.attributeEquals('@deleteAgentGroups','aria-disabled', 'true', "'Confirm agent group delete button is not enabled")
        .verify.visible('@agentGroupsDeleteModal', "Delete agent groups modal is visible")
        .verify.containsText('ngx-agent-group-delete-component.ng-star-inserted > nb-card:nth-child(1) > nb-card-header:nth-child(1)', 'Delete Agent Group Confirmation', "Header of help text about delete confirmation is correctly written")
        .verify.containsText('@agentGroupsDeleteModal', 'Are you sure you want to delete this Agent Group? This may cause Datasets which use this Agent Group to become invalid. This action cannot be undone.*', "Help text about delete confirmation is correctly written")
        .verify.containsText('@agentGroupsDeleteModal', '*To confirm, type the Agent Group name exactly as it appears', "End of help text about delete confirmation is correctly written")
        .getAttribute('.input-full-width', 'placeholder',  function(result) {this.setValue('.input-full-width', result.value) })
        .verify.attributeEquals('@deleteAgentGroups','aria-disabled', 'false', "'Confirm agent group delete button is enabled")
        .click('@deleteAgentGroups')
        .verify.containsText('span.title', 'Agent Group successfully deleted', "Delete confirmation message is being displayed")
        //bug .verify.containsText('.empty-row', 'No data to display')

    },

    agentGroupVisualization: function() {
        return this.verify.attributeEquals('button.orb-action-hover:nth-child(1)', 'aria-disabled', 'false', "'Visualization' button is visible")
        .click('button.orb-action-hover:nth-child(1)')
        .verify.elementPresent('.cdk-overlay-backdrop', "'Visualization' modal is visible")
        .verify.containsText('.nb-card-medium > nb-card-header:nth-child(1)', 'Agent Group Details', "'Visualization' header is correctly written")
    },


    demo_Test: function() {
              return this.waitForElementVisible('[class="orb-action-hover detail-button appearance-ghost size-medium status-basic shape-rectangle icon-start icon-end nb-transition"]')
              .findElements('[class="orb-action-hover detail-button appearance-ghost size-medium status-basic shape-rectangle icon-start icon-end nb-transition"]', function(result) {
              var agentGroupsView = result.value
              this.elementIdClick(agentGroupsView[agentGroupsView.length-1].ELEMENT)
     
      
       })
      },

    addTags: function(key, value) {
      return this.setValue('@key', key)
      .setValue('@value', value)
      .verify.attributeEquals('@addTag','aria-disabled', 'false', "'Add tags' button is enabled")
      .click('@addTag')
      .verify.attributeEquals('@next','aria-disabled', 'false', "'Next' button is enabled")

    },

    //bug : need to insert a test for checking if is possible to create two identicals tags

    agentGroupsEdit: function(name, description, key, value, key2, value2, verify) {
      return this.verify.attributeEquals('@next','aria-disabled', 'false', "'Next' button is enabled")
      .verify.attributeEquals('@back','aria-disabled', 'false', "'Back' button is enabled")
      .clearValue('@newNameInput')
      .setValue('@newNameInput', name)
      .clearValue('@newDescriptionInput')
      .setValue('@newDescriptionInput', description)
      .click('@next')
      .click('.eva-close-outline')
      .verify.attributeEquals('@next','aria-disabled', 'true', "'Next' button is not enabled")
      .verify.attributeEquals('button.status-primary:nth-child(1)','aria-disabled', 'false', "'Back' button is enabled")
      .verify.attributeEquals('@addTag','aria-disabled', 'true', "'Add tags' button is not enabled")
      .addTags(key, value)
      .addTags(key2, value2)
      .verify.attributeEquals('button.status-primary:nth-child(1)','aria-disabled', 'false', "'Back' button is enabled")
      .click('@next')
      .verify.containsText('@agentGroupList','Edit Agent Group', 'Page header is "Edit Agent Group"')
      .verify.attributeEquals('@back','aria-disabled', 'false', "'Back' button is enabled")
      .verify.attributeEquals('@next','aria-disabled', 'false', "'Next' button is enabled")
      .click('@next')
      .verify.containsText('span.title', verify, "Confirmation message is correctly displayed")
    },

    
    agentGroupCheck: function(name, description){
      return this.verify.containsText('div.row:nth-child(1) > div:nth-child(1) > p:nth-child(1)', 'Agent Group Name*', "View contain Agent Group Name Field")
      .verify.containsText('div.row:nth-child(1) > div:nth-child(1) > p:nth-child(2)', name, "Name of Agent Group is correctly displayed")
      .verify.containsText('div.row:nth-child(1) > div:nth-child(2) > p:nth-child(1)', 'Agent Group Description', "View contain Agent Group Description Field")
      .verify.containsText('div.row:nth-child(1) > div:nth-child(2) > p:nth-child(2)', description, "Description of Agent Group is correctly displayed")
      .verify.containsText('div.row:nth-child(2) > div:nth-child(1) > p:nth-child(1)', 'Date Created', "View contain Agent Group Date Created Field")
      .verify.visible('div.row:nth-child(2) > div:nth-child(1) > p:nth-child(2)', "Agent Group Date Created is visible")
      .verify.containsText('div.row:nth-child(2) > div:nth-child(2) > p:nth-child(1)', 'Matches Against', "View contain Agent Group Matches Against Field")
      .verify.containsText('div.row:nth-child(2) > div:nth-child(2) > p:nth-child(2)', 'Agent(s)', "Matches of Agent Group is correctly displayed")
      .verify.containsText('div.row:nth-child(3) > div:nth-child(1) > p:nth-child(1)', 'Tags*', "View contain Agent Group Tags Field")
      .click('@close')
      
    },

  

    countAgentGroups: function(browser) {
      return this.getText('.page-count',  function(result){
        //console.log('Value is:', result.value);
        if (result.value == "0 total") {
            browser.expect.elements('datatable-row-wrapper').count.to.equal(parseInt(result.value))
          
            browser.verify.containsText('.sink-info-accent', 'There are no Agent Groups yet.', "Info message of Agent Groups count is correctly displayed")
            browser.verify.containsText('.empty-row', 'No data to display', "View table info message is correctly displayed")
        } else {
            browser.expect.elements('datatable-row-wrapper').count.to.equal(parseInt(result.value))
            browser.verify.containsText('.justify-content-between > div:nth-child(1)', 'You have', "Beginning of info message is correctly displayed")
            browser.verify.containsText('.justify-content-between > div:nth-child(1)', parseInt(result.value), "Number of Agents is correctly displayed")
            browser.verify.containsText('.justify-content-between > div:nth-child(1)', 'Agent Groups.', "End of info message is correctly displayed")
        }
      })

    }

};

module.exports = {
    url: '/pages/fleet/groups',
    elements: {
        newButton: '.appearance-ghost',
        newHeading: 'header h4',
        activeBreadcrumb: '.xng-breadcrumb-item:last-child .xng-breadcrumb-trail',
        selectedStep: '.selected span',
        completedStep: '.completed span',
        stepLabel: '.step-label strong',
        stepCaption: '.step-label p',
        detailsLabels: '.nb-form-control-container div:not(.d-flex)',
        newNameInput: '[formcontrolname="name"]',
        newDescriptionInput: '[formcontrolname="description"]',
        tagLabels: '.nb-form-control-container div div div',
        keyInput: '[formcontrolname="key"]',
        valueInput: '[formcontrolname="value"]',
        addTagButton: 'button [icon="plus-outline"]', 
        tagChip: '.mat-chip',
        tagChipDelete: '.mat-chip [icon="close-outline"]',
        agentGroupList: 'xng-breadcrumb.orb-breadcrumb',
        agentGroupAll: 'ngx-agent-group-list-component.ng-star-inserted > div:nth-child(1) > header:nth-child(1) > h4:nth-child(2)',
        key: 'div.col-5:nth-child(1) > div:nth-child(2) > input:nth-child(1)',
        value: 'div.d-flex:nth-child(3) > div:nth-child(2) > input:nth-child(1)',
        addTag: 'button.status-basic',
        next:'.next-button',
        back: '.appearance-ghost',
        deleteAgentGroups: '.orb-sink-delete-warning-button',
        agentGroupsDeleteModal: 'ngx-agent-group-delete-component.ng-star-inserted > nb-card:nth-child(1)',
        close: '.nb-close',
        edit: '.sink-edit-button',
        agentPath: '.xng-breadcrumb-link',
        new:'.status-primary',
        filter:'.select-button',
        table:'.datatable-body',
        search:'input.size-medium'

    },
    commands: [commands]
};
