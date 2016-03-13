var Config = require("../common").Config;
var Auth = require("./authAction");
var Nav = require("./navAction");
var Actions = require("./actions").Actions;

//lobby actions
exports.INIT = '/lobby/init'
exports.CREATE = '/game/create'
exports.JOIN = '/game/join'
exports.JOINALERT = '/game/join/alert'
exports.LEAVE = '/game/leave'
exports.LEAVEALERT = '/game/leave/alert'

//lobby registers
Actions.Register(exports.CREATE);
Actions.Register(exports.JOIN, function(dispatch, msg){
    var gameRoute = "/" + msg.game.gameType + "/" + msg.game.gameID;
    Nav.GoToPath(dispatch, gameRoute);
    dispatch(msg);
});
Actions.Register(exports.JOINALERT)
Actions.Register(exports.LEAVE, function(dispatch, msg){
    Nav.GoToPath(dispatch, "/lobby");
});
Actions.Register(exports.LEAVEALERT)

exports.Create = function(dispatch, ws, gameName, gameType){
    var action = {
        type:exports.CREATE,
        game: {gameName: gameName, gameType: gameType}
    };
    ws.jsend(action);
};
exports.Join = function(dispatch, ws, game){
    var action = {
        type:exports.JOIN,
        game
    };
    ws.jsend(action);
};
exports.Leave = function(dispatch, ws, gameID){
    var action = {
        type:exports.LEAVE,
        game: {gameID: gameID}
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
        var games = {};
        for (var i=0; i<json.payload.length; i++){
            var g = json.payload[i];
            games[g.gameID] = g;
        }
        var action = {
            type: exports.INIT,
            games: games
        };
        dispatch(action);
    }).catch(function(err){
        console.log("failed to init lobby ", err);
        dispatch({type: Auth.LOGOUT});
        dispatch({type: Nav.GoToPath, path:"/auth"});
    })
};

