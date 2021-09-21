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

    },
    commands: [commands]
};
