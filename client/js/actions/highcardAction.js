var Config = require("../common").Config;
var Auth = require("./authAction");
var Nav = require("./navAction");
var Actions = require("./actions").Actions;

//highcard actions
exports.INIT = '/highcard/init'
exports.UPDATE = '/highcard/update'
exports.PLAY = '/highcard/play'
exports.CHECK = '/highcard/check'
exports.BET = '/highcard/bet'

//register highcard funcs
Actions.Register(exports.UPDATE)

exports.Initialize = function(dispatch, initialized, gameID){
    if (initialized){return;}
    
    var url = Config.baseURL + "/game/" + gameID
    fetch(url, {headers:{"x-session":reactCookie.load("session") || ""}})
    .then(function(resp){
        if (resp.status != 200){
            throw "status code received: " + resp.status;
        }
        console.log(resp)
        return resp.json();
    }).then(function(json){
        var action = {
            type: exports.INIT,
            gameInfo: json.payload,
            gameState:{players:[]},
            initialized: true
        };
        dispatch(action);
    }).catch(function(err){
        console.log("failed to init highcard game ", err);
        dispatch({type: Auth.LOGOUT});
        dispatch({type: Nav.GoToPath, path:"/auth"});
    })
};

exports.Play = function(ws, gameID){
    var action = {
        type:exports.PLAY,
        gameID:gameID
    };
    ws.jsend(action);
};

exports.Check = function(ws, gameID){
    var action = {
        type:exports.CHECK,
        gameID:gameID
    };
    ws.jsend(action);
};

exports.Bet = function(ws, gameID, amount){
    var action = {
        type:exports.BET,
        gameID:gameID,
        amount:amount
    };
    ws.jsend(action);
};

