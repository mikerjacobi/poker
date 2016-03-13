db = db.getSiblingDB("echo")
load("/fixtures/data.js");
for (var i=0; i<Object.keys(fixtures.users).length;i++){
    var u = fixtures.users["user"+i];
    db.accounts.update({accountID:u.accountID}, u, {upsert:true})
}

db.games.remove({})
