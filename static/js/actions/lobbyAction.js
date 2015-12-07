var Config = require("../common").Config;
var Auth = require("./authAction");
var Nav = require("./navAction");

//lobby actions
exports.INIT = 'LOBBYINIT'
exports.CREATE = 'GAMECREATE'
exports.START = 'GAMESTART'
exports.JOIN = 'GAMEJOIN'
exports.LEAVE = 'GAMELEAVE'

exports.Create = function(dispatch, ws, gameName){
    var action = {
        type:exports.CREATE,
        game: {game_name: gameName}
    };
    ws.jsend(action);
};
exports.Join = function(dispatch, ws, gameID){
    var action = {
        type:exports.JOIN,
        game: {game_id: gameID}
    };
    ws.jsend(action);
};

exports.Initialize = function(dispatch, initialized){
    if (initialized){return;}
    
    var url = Config.baseURL + "/games"
    fetch(url, {headers:{"x-session":reactCookie.load("session") || ""}})
    .then(function(resp){
        if (resp.status != 200){
            throw "status code received: " + resp.status;
        }
        return resp.json();
    }).then(function(json){
        var action = {
            type: exports.INIT,
            games: json.payload
        };
        dispatch(action);
    }).catch(function(err){
        console.log("failed to init lobby ", err);
        dispatch({type: Auth.LOGOUT});
        dispatch({type: Nav.GoToPath, path:"/auth"});
    })
};

