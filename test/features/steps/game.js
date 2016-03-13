var assert = require("assert");

module.exports = function(){
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
    this.When(/^(.*) plays game$/, function(user, done){
        this[user].client
            .click("#play_game_button")
            .pause(500)
            .call(done);
    });
}
