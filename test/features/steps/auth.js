var assert = require("assert");

module.exports = function(){
    this.Given(/^(.*) logs in with (.*)$/, function(user, client, done){
        this[user] = this.fixtures.users[user];
        this[user].client = this.clients[client];
        this[user].client
            .url(this.appHost + "/#/auth")
            .setValue('#username_textfield', this[user].username)
            .setValue('#password_textfield', "111")
            .click('#login_button')
            .pause(500)
            .call(done);
    });
    this.Then(/^(.*) has a session cookie$/, function(user, done){
        this[user].client
            .getCookie("session").then(function(cookie){
                assert.ok(cookie.value.length == 36, "fail session cookie: "+JSON.stringify(cookie))
            })
            .call(done);
    });
}
