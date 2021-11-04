var commands = {
    sinkManagementPage: function() {
        return this.waitForElementVisible('@allSinksPage')
        .verify.containsText('@allSinksPage', 'Sink Management', "Page name is 'Sink Management'")
		.verify.visible('ngx-sink-list-component.ng-star-inserted', "Sink list is visible")

    },

    sinkCreation: function(name_label, description, remote_host, username, password, key, value, verify) {
        return this.verify.attributeEquals('@sinkNext','aria-disabled', 'true', "'Next' button is not enabled")
        .waitForElementVisible('@sinkNameLabel', "Sink Label named field is being displayed")
        .setValue('@sinkNameLabel', name_label)
        .verify.attributeEquals('@sinkNext','aria-disabled', 'false', "'Next' button is enabled")
        .waitForElementVisible('@sinkDescription', "Sink Description field is being displayed")
        .setValue('@sinkDescription', description)
        .verify.attributeEquals('@sinkNext','aria-disabled', 'false', "'Next' button is enabled")
        .click('@sinkNext')
        .verify.attributeEquals('@sinkNext','aria-disabled', 'true', "'Next' button is not enabled")
        .waitForElementVisible('@sinkRemoteHost', "Sink Remote Host field is being displayed")
        .setValue('@sinkRemoteHost', remote_host)
        .verify.attributeEquals('@sinkNext','aria-disabled', 'true', "'Next' button is not enabled")
        .waitForElementVisible('@sinkUsername', "Sink username field is being displayed")
        .setValue('@sinkUsername', username)
        .verify.attributeEquals('@sinkNext','aria-disabled', 'true', "'Next' button is not enabled")
        .waitForElementVisible('@sinkPassword', "Sink password field is being displayed")
        .setValue('@sinkPassword', password)
        .verify.attributeEquals('@sinkNext','aria-disabled', 'false', "'Next' button is enabled")
        .click('@sinkNext')
        .verify.attributeEquals('@sinkNext','aria-disabled', 'true', "'Next' button is not enabled")
        .waitForElementVisible('@key', "Sink key field is being displayed")
        .setValue('@key', key)
        .verify.attributeEquals('@sinkNext','aria-disabled', 'true', "'Next' button is not enabled")
        .waitForElementVisible('@value', "Sink value field is being displayed")
        .setValue('@value', value)
        .verify.attributeEquals('@sinkNext','aria-disabled', 'true', "'Next' button is not enabled")
        .waitForElementVisible('@addTag', "Add tags button is being displayed")
        .click('@addTag')
        .click('@sinkNext')
        .click('@sinkNext')
        .verify.containsText('span.title', verify, "Confirmation message is being correctly displayed")
    },

    sinkVisualization: function() {
        return this.verify.not.elementPresent('.cdk-overlay-backdrop', "Sink visualization modal is not being displayed")
        .verify.attributeEquals('button.orb-action-hover:nth-child(1)', 'aria-disabled', 'false', "Visualization button is enabled")
        .click('button.orb-action-hover:nth-child(1)')
        .verify.elementPresent('.cdk-overlay-backdrop', "Sink visualization modal is being displayed")
        .verify.containsText('.nb-card-medium > nb-card-header:nth-child(1)', 'Sink Details', "Header of visualization page is correctly written")
    },

    sinkEdit: function() {
        return this.verify.not.elementPresent('.cdk-overlay-backdrop', "Sink edit page is not being displayed")
        .verify.attributeEquals('button.orb-action-hover:nth-child(2)', 'aria-disabled', 'false', "Edit sink button is enabled")
        .click('button.orb-action-hover:nth-child(2)')
        .sinkEditPage()

    }, 

    sinkEditPage: function() {
        return this.verify.containsText('.xng-breadcrumb-trail', 'Edit Sink', "Page is named 'Edit Sink'")
        .verify.containsText('ngx-sink-add-component.ng-star-inserted > div:nth-child(1) > header:nth-child(1)', 'Edit Sink', "Option is named 'Edit Sink'")
        .verify.containsText('@editSinkHeader', 'Sink Details', "Sink Details is being displayed")
        .verify.containsText('@editSinkHeader', 'Provide a name and description for the Sink', "Help message for name is correctly written")
        .verify.containsText('@editSinkHeader', 'Sink Destination', "Sink Destination is being correctly displayed")
        .verify.containsText('@editSinkHeader', 'Configure your Sink settings', "Help message for settings is correctly written")
        .verify.containsText('@editSinkHeader', 'Sink Tags', "Sink Tags is being displayed")
        .verify.containsText('@editSinkHeader', 'Enter tags for this Sink', "Help message for tags is correctly written")
        .verify.containsText('@editSinkForm', 'Name Label', "Name Label is being displayed")
        .verify.containsText('@editSinkForm', 'Sink Description', "Sink Description is being displayed")
        .verify.containsText('@editSinkForm', 'Sink Type', "Sink Type is being displayed")
        .verify.attributeEquals('@sinkNext','aria-disabled', 'false', "'Next' button is enabled")

    },

    sinkEditAttribute: function(attribute, value) {
        return this.verify.attributeEquals('@sinkNext','aria-disabled', 'false', "'Next' button is enabled")
        .waitForElementVisible(attribute, "Element edition is enabled")
        .setValue(attribute, value)
        .verify.attributeEquals('@sinkNext','aria-disabled', 'false', "'Next' button is enabled")

    },

    sinkEditTags: function(key, value, key_value, value_value) {
        return this.waitForElementVisible(key, "Key is being displayed")
        .waitForElementVisible(value, "Value is being displayed")
        .setValue(key, key_value)
        .setValue(value, value_value)
        .waitForElementVisible('@addTag', "'Add Tags' button is visible")
        .click('@addTag')
        .verify.attributeEquals('@sinkNext','aria-disabled', 'false', "'Next' button is enabled")

    },


    sinkDelete: function() {
        return this.verify.attributeEquals('button.orb-action-hover:nth-child(3)', 'aria-disabled', 'false', "'Delete' button is enabled")
        .click('button.orb-action-hover:nth-child(3)')
        .verify.attributeEquals('@deleteSink','aria-disabled', 'true', "'Confirm Delete' button is not enabled")
        .verify.visible('@sinkDeleteModal', "Delete modal is visible")
        .verify.containsText('ngx-sink-delete-component.ng-star-inserted > nb-card:nth-child(1) > nb-card-header:nth-child(1)', 'Delete Sink Confirmation', "Header of delete modal is correctly written")
        .verify.containsText('@sinkDeleteModal', 'Are you sure you want to delete this Sink? This may cause Datasets which use this Sink to become invalid. This action cannot be undone.', "Help message for delete is correctly written")
        .verify.containsText('@sinkDeleteModal', 'To confirm, type your Sink name exactly as it appears', "Confirm message is correctly wirtten")
        .getAttribute('.input-full-width', 'placeholder',  function(result) {this.setValue('.input-full-width', result.value) })
        .verify.attributeEquals('@deleteSink','aria-disabled', 'false', "'Confirm Delete' button is enabled")
        .click('@deleteSink')
        .verify.containsText('span.title', 'Sink successfully deleted', "Confirmation message is being correctly displayed")
        //bug insert count
        .verify.containsText('.empty-row', 'No data to display', "List sink is correctly reload")

    },

    sinkCheckEdition: function(value) {
        return this.sinkVisualization()	
        .verify.containsText('ngx-sink-details-component.ng-star-inserted',value, "Element is correctly edited")
        .click('.nb-close')
        .sinkManagementPage()
    },

    countSinks: function(browser) {
        return this.getText('.page-count',  function(result){
          //console.log('Value is:', result.value);
          if (result.value == "0 total") {
              browser.expect.elements('datatable-row-wrapper').count.to.equal(parseInt(result.value))
            
              browser.verify.containsText('.sink-info-accent', 'There are no Sinks yet.', "Info message of Sinks count is correctly displayed")
              browser.verify.containsText('.empty-row', 'No data to display', "View table info message is correctly displayed")
          } else {
              browser.expect.elements('datatable-row-wrapper').count.to.equal(parseInt(result.value))
              browser.verify.containsText('.justify-content-between > div:nth-child(1)', 'sinks total', "Beginning of info message is correctly displayed")
              browser.verify.containsText('.justify-content-between > div:nth-child(1)', parseInt(result.value), "Number of Sinks is correctly displayed")
              browser.verify.containsText('.justify-content-between > div:nth-child(1)', 'have errors.', "End of info message is correctly displayed")
          }
        })}

};

module.exports = {
    url: '/pages/sinks',
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
        editSinkForm: 'form.ng-pristine',
        sinkDeleteModal: 'ngx-sink-delete-component.ng-star-inserted > nb-card:nth-child(1)'
    },
    commands: [commands]
};
