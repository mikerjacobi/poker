var assert = require("assert");

module.exports = function(){
    this.Given(/^(.*) has (.*) account balance$/, function(user, balance, done){
        var accountID = this.fixtures.users[user].accountID;
        this.db.get('accounts').update(
            {accountID:accountID},
            {$set:{balance:parseInt(balance)}})
        done();
    });
    this.When(/^(.*) requests (.*) chips$/, function(user, numChips, done){
        this[user].client
            .setValue('#balance_textfield', numChips)
            .click('#request_chips_button')
            .pause(500)
            .call(done);
    });
    this.Then(/^(.*) should have (.*) chips in their account$/, function(user, numChips, done){
        var accountID = this.fixtures.users[user].accountID;
        this.db.get('accounts').find({accountID:accountID})
        .on('success', function (doc) {
            assert.ok(doc[0].balance == parseInt(numChips), 'predicted: '+numChips+'.  actual: '+doc[0].balance)
            done();
        });
    });
}
