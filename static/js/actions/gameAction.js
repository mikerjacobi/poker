var Config = require("../common").Config;
var Auth = require("./authAction");

//game actions
exports.INIT = 'GAMEINIT'
exports.CREATE = 'GAMECREATE'
exports.START = 'GAMESTART'
exports.JOIN = 'GAMEJOIN'
exports.LEAVE = 'GAMELEAVE'

exports.Initialize = function(dispatch, initialized){
    if (initialized){return;}
    
    var url = Config.baseURL + "/games"
    fetch(url, {headers:{"x-session":reactCookie.load("session") || ""}})
    .then(function(resp){
        return resp.json();
    }).then(function(json){
        var action = {
            type: exports.INIT,
            //games: json.payload.games
            games:[
                {gameID:"abc", gameName:"game one"}, 
                {gameID:"def", gameName:"game two"}
            ]
        };
        dispatch(action);
    }).catch(function(err){
        console.log("failed to init games ", err);
    })
};
