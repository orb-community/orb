var commands = {

    agentGroupsPage: function() {
        return this.assert.containsText('@agentGroupList','Agent Groups List')
        .assert.containsText('@agentGroupAll', 'All Agents Groups')
        .assert.visible('.flex-column')
        .assert.visible('.status-primary')

    },

    agentGroupCreationPage: function() {
        return this.assert.containsText('.header', 'Agent Group Details')
        .assert.containsText('.header', 'This is how you will be able to easily identify your Agent Group')
        .assert.containsText('.header', 'Agent Group Tags')
        .assert.containsText('.header', 'Set the tags that will be used to group Agents')
        .assert.containsText('.header', 'Review & Confirm')
        .assert.containsText('.step-content', 'Agent Group Name*')
        .assert.containsText('.step-content', 'Agent Group Description')
        .assert.attributeEquals('@next','aria-disabled', 'true')
        .assert.attributeEquals('@back','aria-disabled', 'false')
    },

    agentGroupCreation: function(name, description, key, value, assert) {
        return this.setValue('@agentGroupName', name)
        .assert.attributeEquals('@next','aria-disabled', 'false')
        .setValue('@agentGroupDescription', description)
        .click('@next')
        .assert.attributeEquals('@next','aria-disabled', 'true')
        .assert.attributeEquals('button.status-primary:nth-child(1)','aria-disabled', 'false')
        .assert.attributeEquals('@addTag','aria-disabled', 'true')
        .setValue('@key', key)
        .setValue('@value', value)
        .click('@addTag')
        .assert.attributeEquals('@next','aria-disabled', 'false')
        .assert.attributeEquals('button.status-primary:nth-child(1)','aria-disabled', 'false')
        .click('@next')
        .assert.containsText('@agentGroupList','New Agent Group')
        .assert.attributeEquals('@back','aria-disabled', 'false')
        .assert.attributeEquals('@next','aria-disabled', 'false')
        .click('@next')
        .assert.containsText('span.title', assert)


    },

    agentGroupsDelete: function() {
        return this.assert.attributeEquals('button.orb-action-hover:nth-child(3)', 'aria-disabled', 'false')
        .click('button.orb-action-hover:nth-child(3)')
        .assert.attributeEquals('@deleteAgentGroups','aria-disabled', 'true')
        .assert.visible('@agentGroupsDeleteModal')
        .assert.containsText('ngx-agent-group-delete-component.ng-star-inserted > nb-card:nth-child(1) > nb-card-header:nth-child(1)', 'Delete Agent Group Confirmation')
        .assert.containsText('@agentGroupsDeleteModal', 'Are you sure you want to delete this Agent Group? This may cause Datasets which use this Agent Group to become invalid. This action cannot be undone.*')
        .assert.containsText('@agentGroupsDeleteModal', '*To confirm, type the Agent Group name exactly as it appears')
        .getAttribute('.input-full-width', 'placeholder',  function(result) {this.setValue('.input-full-width', result.value) })
        .assert.attributeEquals('@deleteAgentGroups','aria-disabled', 'false')
        .click('@deleteAgentGroups')
        .assert.containsText('span.title', 'Agent Group successfully deleted')
        // .assert.containsText('.empty-row', 'No data to display')

    },

    agentGroupVisualization: function() {
        return this.assert.not.elementPresent('.cdk-overlay-backdrop')
        .assert.attributeEquals('button.orb-action-hover:nth-child(1)', 'aria-disabled', 'false')
        .click('button.orb-action-hover:nth-child(1)')
        .assert.elementPresent('.cdk-overlay-backdrop')
        .assert.containsText('.nb-card-medium > nb-card-header:nth-child(1)', 'Agent Group Details')
    },

    agentGroupEditPage: function() {
        return this.assert.containsText('.xng-breadcrumb-trail', 'Edit Agent Group')
        .assert.containsText('ngx-sink-add-component.ng-star-inserted > div:nth-child(1) > header:nth-child(1)', 'Update Sink')
        .assert.containsText('@editSinkHeader', 'Agent Group Details')
        .assert.containsText('@editSinkHeader', 'This is how you will be able to easily identify your Agent Group')
        .assert.containsText('@editSinkHeader', 'Agent Group Tags')
        .assert.containsText('@editSinkHeader', 'Set the tags that will be used to group Agents')
        .assert.containsText('@editSinkHeader', 'Review & Confirm')
        .assert.attributeEquals('@sinkNext','aria-disabled', 'false')
    },

};

module.exports = {
    url: 'http://localhost:4200',
    elements: {
        agentGroupList: 'xng-breadcrumb.orb-breadcrumb',
        agentGroupAll: 'ngx-agent-group-list-component.ng-star-inserted > div:nth-child(1) > header:nth-child(1) > h4:nth-child(2)',
        agentGroupCreation: '.appearance-ghost',
        agentGroupName: 'input.ng-invalid',
        agentGroupDescription: '.ng-pristine',
        key: 'div.col-5:nth-child(1) > div:nth-child(2) > input:nth-child(1)',
        value: 'div.d-flex:nth-child(3) > div:nth-child(2) > input:nth-child(1)',
        addTag: 'button.status-basic',
        next:'.next-button',
        back: '.appearance-ghost',
        deleteAgentGroups: '.orb-sink-delete-warning-button',
        agentGroupsDeleteModal: 'ngx-agent-group-delete-component.ng-star-inserted > nb-card:nth-child(1)',
        close: '.nb-close',
        edit: '.sink-edit-button'

    },
    commands: [commands]
};
