var assert = require("assert");
var timeoutDuration = 15000;
var exec = require('child_process').execSync;
var webdrivercssDir = "poker/"

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
                    try{
                        var msg = "ERROR: "+e.body.value.class+" >>> "+e.body.value.message
                        console.log(msg);
                    } catch(err){
                        //pass
                    }
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
            case "account":
                uri = "/#/account";
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
            .pause(200)
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
