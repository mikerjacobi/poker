module.exports = function () {
    this.Given(/^(.*) navigates to login$/, function(user, done){
        this.user = this.fixtures.users[user];
        this.url = this.appHost + "/#/auth";
        //this.visit(this.url, done);
        done();
    });
    this.Given(/^(.*) is screenshot/, 10000, function(ssName, done){
        this.client.init().sync()
            .url(this.url)
            .webdrivercss(ssName, [{
                name: 'element',
                elem: '#root',
                screenWidth: [640]
            }], function(err,res) {
                if (err != undefined){
                  //console.log("err: "+err);
                } 
               //assert.ok(res.element[0].isWithinMisMatchTolerance);
            }).sync()
            //.end()
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

  this.When(/^the index page is loaded$/, function (done) {
    var url = this.appHost + "/";
    this.client
        .init()
        .sync()
        .url(url)
        .getTitle().then(function(title) {
            console.log('title: ' + title);
        })
        .webdrivercss('indexPage', [{
            name: 'element',
            elem: '#root',
            screenWidth: [640]
        }], function(err,res) {
            if (err != undefined){
              console.log("err: "+err);
            } else {
              console.log("success");
            }
           //assert.ok(res.element[0].isWithinMisMatchTolerance);
        })
        .sync()
        //.end()
        .call(done);
    });
};
