var assert = require("assert");
var timeoutDuration = 150000;

module.exports = function () {
    this.setDefaultTimeout(timeoutDuration)
    this.Before(function(scenario){
        this.client
            .init()
            .timeouts("page load",timeoutDuration)
            .timeouts("implicit",timeoutDuration)
            .timeouts("script",timeoutDuration)
            .on('error', function(e) {
                var msg = "ERROR: "+e.body.value.class+" >>> "+e.body.value.message
                console.log(msg);
                //assert.ok(false, msg);
            })
    });
    this.After(function(scenario){
        this.client.end();
    });
    this.Given(/^we wait (.*) seconds$/, function(seconds, done){
        setTimeout(function(){done();}, seconds * 1000);
    });

    this.Given(/^user navigates to index$/, function(done){
        this.url = this.appHost + "/";
        done();
    });

    this.Given(/^(.*) logs in$/, function(user, done){
        this.user = this.fixtures.users[user];
        this.url = this.appHost + "/#/auth";
        this.client
            .url(this.url)
            .setValue('#username_textfield', this.user.username)
            .setValue('#password_textfield', "111")
            .click('#login_button')
            .call(done);
    });
    this.Given(/^(.*) is screenshot$/, function(ssName, done){
        var wdcssRes;
        this.client
            .sync()
            .url(this.url)
            .webdrivercss(ssName, [{
                name: 'element',
                elem: '#root',
                screenWidth: [640]
            }], function(err,res) {
                assert.ifError(err);
                webcssRes = res;
            }).sync()
            .call(function(){
                assert.ok(
                    webcssRes.element[0].isWithinMisMatchTolerance,
                    "pdiff not within tolerance");
            })
            .call(done);
    });
    this.Then(/^user has a session cookie$/, function(done){
        this.client
            .url(this.url)
            .getCookie("session").then(function(cookie){
                assert.ok(cookie.value.length == 36, "fail session cookie: "+JSON.stringify(cookie))
            })
            .call(done);
    });
};
