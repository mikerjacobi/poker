var assert = require("assert");

module.exports = function () {
    this.Given(/^(.*) navigates to (.*)/, function(user, page, done){
        this.user = this.fixtures.users[user];
        this.url = this.appHost + "/#/" + page 
        //this.visit(this.url, done);
        done();
    });
    this.Given(/^(.*) is screenshot/, 10000, function(ssName, done){
        var wdcssRes;
        this.client
            .init()
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
            .end()
            .call(function(){
                assert.ok(
                    webcssRes.element[0].isWithinMisMatchTolerance,
                    "pdiff not within tolerance");
            })
            .call(done);
    });
    this.When(/^user logs in$/, function(done){
        console.log(this.user);
        done();
    });
    this.Given(/^user has a session cookie$/, function(done){
        console.log("session cookie");
        done();
    });
};
