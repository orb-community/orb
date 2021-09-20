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

// 'Sink Creation' : function(browser) {
// 	var login = browser.page.login()
// 	login
// 		.assert.visible('li.menu-item:nth-child(4)')
// 		.click('li.menu-item:nth-child(4)')
// 		.sinkManagementPage()
// 		.assert.visible('.appearance-ghost')
// 		.assert.attributeEquals('.appearance-ghost','aria-disabled', 'false')
// 		.click('.appearance-ghost')
// 		.sinkCreation('some_name', 'some_description', 'remote_host', 'tester', 'password', 'key', 'value', 'Sink successfully created')
// 		.sinkManagementPage()		
// 	},

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


'Sink Visualization and Edit' : function(browser){
	var login = browser.page.login()
	login
	.sinkVisualization()
	.click('.nb-close')
	.sinkVisualization()
	.click('.sink-edit-button')
	.sinkEditPage()
	.sinkEditDescription('new_desc')
	.click('@sinkNext')
	.click('@sinkNext')
	.click('@sinkNext')
	.click('@sinkNext')
	.assert.containsText('span.title', 'Sink successfully updated')

},

'Sink Edit' : function(browser){
	var login = browser.page.login()
	login
	// .sinkEditDescription()
},

'Sink Delete' : function(browser){
	var login = browser.page.login()
	login
	// .sinkDelete()
	// .sinkManagementPage()

}

}
