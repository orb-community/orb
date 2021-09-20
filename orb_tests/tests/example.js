module.exports = {
'Create an account' : function(browser) {
	var login = browser.page.login()
	login
		.navigate()
		.orbRegister()
		.end


	},


'Login with correct email and password' : function(browser) {
 	var login = browser.page.login()
	login
	   .navigate()
	   .assert.titleContains('ORB')
	   .assert.visible('#orb-pane-div > p:nth-child(2)')
	   .assert.containsText('#orb-pane-div > p:nth-child(2)','An Open-Source dynamic edge observability platform')
	   .fillLoginForm()
	   .submitLogin()

  	},

'Sink Creation' : function(browser) {
	var login = browser.page.login()
	login
		.assert.visible('li.menu-item:nth-child(4)')
		.click('li.menu-item:nth-child(4)')
		.sinkManagementPage()
		.assert.containsText('.page-count', '0 total')
		.assert.visible('.appearance-ghost')
		.assert.attributeEquals('.appearance-ghost','aria-disabled', 'false')
		.click('.appearance-ghost')
		.sinkCreation('some_name', 'some_description', 'remote_host', 'tester', 'password', 'key', 'value', 'Sink successfully created')
		.sinkManagementPage()
		.assert.containsText('.page-count', '1 total')		
	},

'Duplicate sink creation' : function(browser) {
	var login = browser.page.login()
	login
		.assert.visible('li.menu-item:nth-child(4)')
		.click('li.menu-item:nth-child(4)')
		.sinkManagementPage()
		.assert.visible('.appearance-ghost')
		.assert.attributeEquals('.appearance-ghost','aria-disabled', 'false')
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

},

'Sink Visualization' : function(browser){
	var login = browser.page.login()
	login
	.sinkVisualization()
	.click('.nb-close')
	.sinkVisualization()
	.click('.sink-edit-button')
	.sinkEditPage()
	.click('@cancel')
},


'Sink Visualization and Edit Description' : function(browser){
	var login = browser.page.login()
	login
	.sinkVisualization()
	.click('.nb-close')
	.sinkVisualization()
	.click('.sink-edit-button')
	.sinkEditPage()
	.sinkEditAttribute('@sinkDescription','_new')
	.click('@sinkNext')
	.click('@sinkNext')
	.click('@sinkNext')
	.click('@sinkNext')
	.assert.containsText('span.title', 'Sink successfully updated')
	.sinkCheckEdition('_new')

},

'Sink Visualization and Edit Remote Host' : function(browser){
	var login = browser.page.login()
	login
	.sinkVisualization()
	.click('.nb-close')
	.sinkVisualization()
	.click('.sink-edit-button')
	.sinkEditPage()
	.click('@sinkNext')
	.sinkEditAttribute('@sinkRemoteHost', '_new_rm')
	.click('@sinkNext')
	.click('@sinkNext')
	.click('@sinkNext')
	.assert.containsText('span.title', 'Sink successfully updated')
	.sinkCheckEdition('_new_rm')


},

'Sink Visualization and Edit Username' : function(browser){
	var login = browser.page.login()
	login
	.sinkVisualization()
	.click('.nb-close')
	.sinkVisualization()
	.click('.sink-edit-button')
	.sinkEditPage()
	.click('@sinkNext')
	.sinkEditAttribute('@sinkUsername', '_new_usr')
	.click('@sinkNext')
	.click('@sinkNext')
	.click('@sinkNext')
	.assert.containsText('span.title', 'Sink successfully updated')
	//BUG
	//.sinkCheckEdition('_new_usr')


},

'Sink Visualization and Edit Password' : function(browser){
	var login = browser.page.login()
	login
	.sinkVisualization()
	.click('.nb-close')
	.sinkVisualization()
	.click('.sink-edit-button')
	.sinkEditPage()
	.click('@sinkNext')
	.sinkEditAttribute('@sinkPassword', '_new_pass')
	.click('@sinkNext')
	.click('@sinkNext')
	.click('@sinkNext')
	.assert.containsText('span.title', 'Sink successfully updated')

},

'Sink Visualization and Edit Tags' : function(browser){
	var login = browser.page.login()
	login
	.sinkVisualization()
	.click('.nb-close')
	.sinkVisualization()
	.click('.sink-edit-button')
	.sinkEditPage()
	.click('@sinkNext')
	.click('@sinkNext')
	.click('.ml-1')
	.assert.attributeEquals('@submit','aria-disabled', 'true')
	.sinkEditTags('@key', '@value', 'new_key', 'new_value')
	.click('@sinkNext')
	.click('@sinkNext')
	.assert.containsText('span.title', 'Sink successfully updated')
	.sinkCheckEdition('new_key')
	.sinkCheckEdition('new_value')
},



'Sink Edit Description' : function(browser){
	var login = browser.page.login()
	login
	.sinkEdit()
	.sinkEditAttribute('@sinkDescription','_n3w')
	.click('@sinkNext')
	.click('@sinkNext')
	.click('@sinkNext')
	.click('@sinkNext')
	.assert.containsText('span.title', 'Sink successfully updated')
	.sinkCheckEdition('_n3w')

},

'Sink Edit Remote Host' : function(browser){
	var login = browser.page.login()
	login
	.sinkEdit()
	.click('@sinkNext')
	.sinkEditAttribute('@sinkRemoteHost', '_n3w_rm')
	.click('@sinkNext')
	.click('@sinkNext')
	.click('@sinkNext')
	.assert.containsText('span.title', 'Sink successfully updated')
	.sinkCheckEdition('_n3w_rm')

},

'Sink Edit Username' : function(browser){
	var login = browser.page.login()
	login
	.sinkEdit()
	.click('@sinkNext')
	.sinkEditAttribute('@sinkUsername', '_n3w_usr')
	.click('@sinkNext')
	.click('@sinkNext')
	.click('@sinkNext')
	.assert.containsText('span.title', 'Sink successfully updated')
	// BUG
	//.sinkCheckEdition('_n3w_usr')

},

'Sink Edit Password' : function(browser){
	var login = browser.page.login()
	login
	.sinkEdit()
	.click('@sinkNext')
	.sinkEditAttribute('@sinkPassword', '_n3w_pass')
	.click('@sinkNext')
	.click('@sinkNext')
	.click('@sinkNext')
	.assert.containsText('span.title', 'Sink successfully updated')

},

'Sink Edit Tags' : function(browser){
	var login = browser.page.login()
	login
	.sinkEdit()
	.click('@sinkNext')
	.click('@sinkNext')
	.click('.ml-1')
	.assert.attributeEquals('@submit','aria-disabled', 'true')
	.sinkEditTags('@key', '@value', '_n3w_key', '_n3w_value')
	.click('@sinkNext')
	.click('@sinkNext')
	.assert.containsText('span.title', 'Sink successfully updated')
	.sinkCheckEdition('_n3w_key')
	.sinkCheckEdition('_n3w_value')
},

'Sink Edit Description, Remote Host, Username, Password and Keys' : function(browser){
	var login = browser.page.login()
	login
	.sinkEdit()
	.sinkEditAttribute('@sinkDescription','_n33w')
	.click('@sinkNext')
	.sinkEditAttribute('@sinkRemoteHost', '_n33w_rm')
	.sinkEditAttribute('@sinkUsername', '_n33w_usr')
	.sinkEditAttribute('@sinkPassword', '_n33w_pass')
	.click('@sinkNext')
	.click('.ml-1')
	.assert.attributeEquals('@submit','aria-disabled', 'true')
	.sinkEditTags('@key', '@value', '_n33w_key', '_n33w_value')
	.click('@sinkNext')
	.click('@sinkNext')
	.assert.containsText('span.title', 'Sink successfully updated')
	.sinkCheckEdition('_n33w')
	.sinkCheckEdition('_n33w_rm')
	// BUG
	//.sinkCheckEdition('_n33w_usr')
	.sinkCheckEdition('_n33w_key')
	.sinkCheckEdition('_n33w_value')
},

'Sink Delete' : function(browser){
	var login = browser.page.login()
	login
	.assert.containsText('.page-count', '1 total')
	.sinkDelete()
	.assert.containsText('.page-count', '0 total')
	.sinkManagementPage()
	

}

}
