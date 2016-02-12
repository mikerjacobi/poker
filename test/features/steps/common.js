var assert = require("assert");
var timeoutDuration = 15000;

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
            })
    });
    this.After(function(scenario){
        this.db.close();
        this.client.end();
    });
    this.Given(/^we wait (.*) seconds$/, function(seconds, done){
        this.client
            .then(
                setTimeout(function(){
                    done();
                }, seconds * 1000)
            );
    });

    this.Given(/^user navigates to (.*)$/, function(route, done){
        var uri;
        switch (route){
            case "lobby":
                uri = "/#/lobby";
                break;
            case "login":
                uri = "/#/auth?";
                break;
            default:
                uri = "/";
        }

        this.client
            .url(this.appHost + uri)
            .call(done);
    });

    this.Given(/^(.*) logs in$/, function(user, done){
        this.user = this.fixtures.users[user];
        this.client
            .url(this.appHost + "/#/auth")
            .setValue('#username_textfield', this.user.username)
            .setValue('#password_textfield', "111")
            .click('#login_button')
            .call(done);
    });
    this.When(/^(.*) is screenshot$/, function(ssName, done){
        var wdcssRes;
        this.client
            .sync()
            .webdrivercss(ssName, [{
                name: 'element',
                elem: '#root',
                screenWidth: [800]
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
            .getCookie("session").then(function(cookie){
                assert.ok(cookie.value.length == 36, "fail session cookie: "+JSON.stringify(cookie))
            })
            .call(done);
    });
    this.Given(/^there are no games$/, function(done){
        var games = this.db.get('games');
        games.remove({});
        done();
    });
    this.When(/^user creates (.*) game$/, function(gameType, done){
        var typeNum;
        switch (gameType){
            case "holdem":
                typeNum = 2;
                break;
            case "highcard":
                typeNum = 3;
                break;
            default:
                typeNum = 1;
        }

        this.client
            .setValue('#gamename_textfield', "test_"+gameType)
            .click('//*[@id="gametype_dropdown"]/option[' + typeNum + ']')
            .click('#create_game_button')
            .call(done);
    });
    this.When(/^user joins game$/, function(done){
        this.client
            .url(this.appHost + "/#/lobby")
            .click("#game_listing_0")
            .call(done);
    });
    this.Given(/^there is a (.*) game$/, function(gameType, done){
        var games = this.db.get('games');
        games.insert({
            gameID : "9cf67eab-8a93-4005-9a5b-b9cf678a6cb9", 
            gameName : "test_"+gameType, 
            state : "open", 
            players : [ ], 
            gameType : gameType
        }, function (err, doc) {
            assert.ifError(err);
            done();
        });
    });
};
