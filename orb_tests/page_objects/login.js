var commands = {

    orbRegister: function() {
        return this.assert.visible('@register')
		.click('@register')
		//why??
		.click('@register')
		.waitForElementVisible('@fullNameInput')
        .setValue('@fullNameInput', 'tester')
		.waitForElementVisible('@username')
        .setValue('@username', 'tester@email.com')
		.waitForElementVisible('@pwd')
        .setValue('@pwd', '12345678')
		.waitForElementVisible('@confirmPassword')
        .setValue('@confirmPassword', '12345678')
        .waitForElementVisible('@submit')
        .assert.attributeEquals('@submit','aria-disabled', 'false')
        .click('@submit')

    },


    fillLoginForm: function() {
        return this.waitForElementVisible('@username')
        .setValue('@username', 'tester@email.com')
        .waitForElementVisible('@pwd')
        .setValue('@pwd', '12345678')
    },
    submitLogin: function() {
        return this.waitForElementVisible('@submit')
        .assert.attributeEquals('@submit','aria-disabled', 'false')
        .click('@submit')
    },

    sinkManagementPage: function() {
        return this.waitForElementVisible('@allSinksPage')
        .assert.containsText('@allSinksPage', 'Sink Management')
		.assert.visible('ngx-sink-list-component.ng-star-inserted')

    },
    
    
    sinkCreation: function(name_label, description, remote_host, username, password, key, value, assert) {
        return this.assert.attributeEquals('@sinkNext','aria-disabled', 'true')
        .waitForElementVisible('@sinkNameLabel')
        .setValue('@sinkNameLabel', name_label)
        .assert.attributeEquals('@sinkNext','aria-disabled', 'false')
        .waitForElementVisible('@sinkDescription')
        .setValue('@sinkDescription', description)
        .assert.attributeEquals('@sinkNext','aria-disabled', 'false')
        .click('@sinkNext')
        .assert.attributeEquals('@sinkNext','aria-disabled', 'true')
        .waitForElementVisible('@sinkRemoteHost')
        .setValue('@sinkRemoteHost', remote_host)
        .assert.attributeEquals('@sinkNext','aria-disabled', 'true')
        .waitForElementVisible('@sinkUsername')
        .setValue('@sinkUsername', username)
        .assert.attributeEquals('@sinkNext','aria-disabled', 'true')
        .waitForElementVisible('@sinkPassword')
        .setValue('@sinkPassword', password)
        .assert.attributeEquals('@sinkNext','aria-disabled', 'false')
        .click('@sinkNext')
        .assert.attributeEquals('@sinkNext','aria-disabled', 'true')
        .waitForElementVisible('@key')
        .setValue('@key', key)
        .assert.attributeEquals('@sinkNext','aria-disabled', 'true')
        .waitForElementVisible('@value')
        .setValue('@value', value)
        .assert.attributeEquals('@sinkNext','aria-disabled', 'true')
        .waitForElementVisible('@addTag')
        .click('@addTag')
        .click('@sinkNext')
        .click('@sinkNext')
        .assert.containsText('span.title', assert)
    },

    sinkVisualization: function() {
        return this.assert.not.elementPresent('.cdk-overlay-backdrop')
        .assert.attributeEquals('button.orb-action-hover:nth-child(1)', 'aria-disabled', 'false')
        .click('button.orb-action-hover:nth-child(1)')
        .assert.elementPresent('.cdk-overlay-backdrop')
        .assert.containsText('.nb-card-medium > nb-card-header:nth-child(1)', 'Sink Details')
    },

    sinkEditPage: function() {
        return this.assert.containsText('.xng-breadcrumb-trail', 'Edit Sink')
        .assert.containsText('ngx-sink-add-component.ng-star-inserted > div:nth-child(1) > header:nth-child(1)', 'Update Sink')
        .assert.containsText('div.step:nth-child(1) > div:nth-child(2) > div:nth-child(1)', 'Sink Details')
        .assert.containsText('@editSinkHeader', 'Provide a name and description for the Sink')
        .assert.containsText('@editSinkHeader', 'Sink Destination')
        .assert.containsText('@editSinkHeader', 'Configure your sink settings')
        .assert.containsText('@editSinkHeader', 'Sink Tags')
        .assert.containsText('@editSinkHeader', 'Enter tags for this sink')
        .assert.containsText('@editSinkForm', 'Name Label')
        .assert.containsText('@editSinkForm', 'Sink Description')
        .assert.containsText('@editSinkForm', 'Sink Type')
        .assert.attributeEquals('@sinkNext','aria-disabled', 'false')

    },

    sinkEditDescription: function(description) {
        return this.waitForElementVisible('@sinkDescription')
        .setValue('@sinkDescription', description)
        .assert.attributeEquals('@sinkNext','aria-disabled', 'false')

    },
    
    sinkDelete: function() {
        return this.assert.attributeEquals('button.orb-action-hover:nth-child(3)', 'aria-disabled', 'false')
        .click('button.orb-action-hover:nth-child(3)')
        .assert.attributeEquals('@deleteSink','aria-disabled', 'true')
        .getAttribute('.input-full-width', 'placeholder',  function(result) {this.setValue('.input-full-width', result.value) })
        .assert.attributeEquals('@deleteSink','aria-disabled', 'false')
        .click('@deleteSink')
        .assert.containsText('span.title', 'Sink successfully deleted')
        .assert.containsText('.empty-row', 'No data to display')

    },

};

module.exports = {
    url: 'http://localhost:4200',
    elements: {
        username: 'input[id=input-email]',
        pwd: 'input[id=input-password]',
        submit: '.appearance-filled',
        loginBody: 'div.login_wrapper',
        allSinksPage : '.xng-breadcrumb-trail',
        sinkNameLabel:'input[data-orb-qa-id=name]',
        sinkDescription: 'input[data-orb-qa-id=description]',
        sinkNext:  'button[data-orb-qa-id=next]',
        sinkRemoteHost: 'input[data-orb-qa-id=remote_host]',
        sinkUsername: 'input[data-orb-qa-id=username]',
        sinkPassword: 'input[data-orb-qa-id=password]',
        key: 'input[data-orb-qa-id=key]',
        value: 'input[data-orb-qa-id=value]',
        addTag: 'button[data-orb-qa-id=addTag]',
        spanTitle: 'span.title',
        register: '.text-link',
        fullNameInput:'input[id=input-name]',
        confirmPassword: 'input[id=input-re-password]',
        deleteSink: '.orb-sink-delete-warning-button',
        cancel: 'button[data-orb-qa-id=cancel]',
        back: 'button[data-orb-qa-id=back]',
        previous: 'button[data-orb-qa-id=previous]',
        editSinkHeader: '.header',
        editSinkForm: 'form.ng-pristine'
    },
    commands: [commands]
};
