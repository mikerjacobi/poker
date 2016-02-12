var assert = require("assert");
var timeoutDuration = 15000;
var exec = require('child_process').execSync;
var webdrivercssDir = "poker-tests/"

function DiffScreenShots(ss1, ss2){
    var cmd = "compare -metric MAE " + webdrivercssDir + ss1 + 
        ".game_div.800px.baseline.png " + webdrivercssDir + ss2 + 
        ".game_div.800px.baseline.png null: 2>&1"
    var raw = String(exec(cmd)).replace("\n","");
    var numDiffPixels = parseFloat(raw.split(" ")[0]);
    return numDiffPixels;
}

function DeleteScreenShot(ss){
    //delete game screenshot file
    var cmd = "rm " + webdrivercssDir + ss + ".game_div.800px.baseline.png";
    exec(cmd)
    var cmd = "rm " + webdrivercssDir + ss + ".800px.png"
    exec(cmd)
}

module.exports = function () {
    this.setDefaultTimeout(timeoutDuration)
    this.Before(function(scenario){
        for (var i=1; i<=Object.keys(this.clients).length; i++){
            this.clients["cli"+i]
                .init()
                .timeouts("page load",timeoutDuration)
                .timeouts("implicit",timeoutDuration)
                .timeouts("script",timeoutDuration)
                .on('error', function(e) {
                    var msg = "ERROR: "+e.body.value.class+" >>> "+e.body.value.message
                    console.log(msg);
                })
        }
    });
    this.After(function(scenario){
        this.db.close();
        for (var i=1; i<=Object.keys(this.clients).length; i++){
            this.clients["cli"+i].end();
        }
    });
    this.Given(/^(.*) navigates to (.*)$/, function(user, route, done){
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

        if (this[user] == undefined){
            this[user] = {client: this.clients.cli1};
        }

        this[user].client
            .url(this.appHost + uri)
            .call(done);
    });

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
    this.When(/^(.*) is screenshot$/, function(ssName, done){
        var wdcssRes;
        this.clients.cli1
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
    this.Then(/^(.*) has a session cookie$/, function(user, done){
        this[user].client
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
    this.When(/^(.*) creates (.*) game$/, function(user, gameType, done){
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

        this[user].client
            .setValue('#gamename_textfield', "test_"+gameType)
            .click('//*[@id="gametype_dropdown"]/option[' + typeNum + ']')
            .click('#create_game_button')
            .pause(500)
            .call(done);
    });
    this.When(/^(.*) joins game$/, function(user, done){
        this[user].client
            .url(this.appHost + "/#/lobby")
            .click("#game_listing_0")
            .pause(500)
            .call(done);
    });
    this.Given(/^there is a (.*) game$/, function(gameType, done){
        this.gameType = gameType;
        var games = this.db.get('games');
        games.insert({
            gameID : "9cf67eab-8a93-4005-9a5b-b9cf678a6cb9", 
            gameName : "test_"+gameType, 
            state : "open", 
            players : [ ], 
            gameType : gameType
        }, function (err, doc) {
            assert.ifError(err);
            setTimeout(function(){
                done();
            }, 1000)
        });
    });
    this.When(/^(.*) replays game$/, function(user, done){
        this[user].client
            .click("#replay_game_button")
            .pause(500)
            .call(done);
    });
    this.When(/^(.*) screenshots game$/, function(user, done){
        var currScreenShot = user + "." + this.gameType
        this.gameScreenShots = (this.gameScreenShots || []).concat([currScreenShot])

        this[user].client
            .webdrivercss(currScreenShot, [{
                name: 'game_div',
                elem: '#game_div',
                screenWidth: [800]
            }], function(err,res) {
                assert.ifError(err);
            })
            .call(done)
    })
    this.Then(/^game screenshots should be equal/, function(done){
        assert.ok(this.gameScreenShots.length >= 2, "not enough game screenshots to diff")

        var compareShot = this.gameScreenShots[0]
        //iterate over all game screenshots and diff them with the first screenshot
        for (var i=1; i<this.gameScreenShots.length; i++){
            var numDiffPixels = DiffScreenShots(compareShot, this.gameScreenShots[i]);
            var failMsg = this.gameScreenShots[i] + " and " + compareShot + " should be equal, but are not.  result: " + numDiffPixels
            assert.ok(numDiffPixels < 15, failMsg)
            DeleteScreenShot(this.gameScreenShots[i])
        }
        DeleteScreenShot(compareShot);
        done();
    })
};
