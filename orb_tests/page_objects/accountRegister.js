var registerCommands = {
    orbRegister: function() {
        return this.verify.visible('@register', "'Register' option is being displayed")
		.click('@register')
		//bug remove this duplicated
		.click('@register')
		.waitForElementVisible('@fullNameInput', "'Full name' field is visible")
        .verify.containsText("[for='input-name']", "Full name:", "'Full name:' description is correctly written")
        .setValue('@fullNameInput', 'tester')
		.waitForElementVisible('@username', "'Email address' field is visible")
        .verify.containsText("[for='input-email']", "Email address:", "'Email Address' description is correctly written")
        .setValue('@username', 'tester@email.com')
		.waitForElementVisible('@pwd', "'Password field is visible")
        .verify.containsText("[for='input-password']", "Password:", "'Password' description is correctly written")
        .setValue('@pwd', '12345678')
		.waitForElementVisible('@confirmPassword', "'Repeat password' field is visible")
        .verify.containsText("[for='input-re-password']", "Repeat password:", "'Repeat password' description is correctly written")
        .setValue('@confirmPassword', '12345678')
        .waitForElementVisible('@submit', "'Register' button is visible")
        .verify.containsText('@submit', "REGISTER", "'REGISTER' text is correctly written")
        .verify.attributeEquals('@submit','aria-disabled', 'false', "'Register' button is clickable")
        .pause(2000)
        .click('@submit')
        .pause(2000)

    },
}
module.exports = {
    url: '/auth/register',
    commands: [registerCommands],
    elements: {
        register: '.text-link',
        fullNameInput:'input[id=input-name]',
        username: 'input[id=input-email]',
        pwd: 'input[id=input-password]',
        confirmPassword: 'input[id=input-re-password]',
        submit: '.appearance-filled'

    }}
