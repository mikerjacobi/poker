var Config = require("../common").Config;
var Auth = require("./authAction");
var Nav = require("./navAction");

//game actions
exports.INIT = 'GAMEINIT'
exports.CREATE = 'GAMECREATE'
exports.START = 'GAMESTART'
exports.JOIN = 'GAMEJOIN'
exports.LEAVE = 'GAMELEAVE'

//game form actions
exports.CHANGEGAMENAME = 'CHANGEGAMENAME';

exports.Create = function(dispatch, ws, gameName){
    //dispatch({type:GameForm.FETCH});
    var action = {
        type:exports.CREATE,
        game: {game_name: gameName}
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
        console.log("failed to init games ", err);
        dispatch({type: Auth.LOGOUT});
        dispatch({type: Nav.GoToPath, path:"/auth"});
    })
};

exports.ChangeGameName = function(dispatch, gameName){
    var action = {
        type:exports.CHANGEGAMENAME,
        createGameName:gameName
    }
    dispatch(action);
}
