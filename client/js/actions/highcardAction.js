var Config = require("../common").Config;
var Auth = require("./authAction");
var Nav = require("./navAction");
var Actions = require("./actions").Actions;

//highcard actions
exports.INIT = 'HIGHCARDINIT'
exports.UPDATE = 'HIGHCARDUPDATE'
exports.REPLAY = 'HIGHCARDREPLAY'
exports.ERROR = 'HIGHCARDERROR'

//register highcard funcs
Actions.Register(exports.UPDATE)
Actions.Register(exports.ERROR)

exports.Initialize = function(dispatch, initialized, gameID){
    if (initialized){return;}
    
    var url = Config.baseURL + "/game/" + gameID
    fetch(url, {headers:{"x-session":reactCookie.load("session") || ""}})
    .then(function(resp){
        if (resp.status != 200){
            throw "status code received: " + resp.status;
        }
        return resp.json();
    }).then(function(json){
        var action = {
            type: exports.INIT,
            gameInfo: json.payload,
            gameState:{},
            initialized: true
        };
        dispatch(action);
    }).catch(function(err){
        console.log("failed to init highcard game ", err);
        dispatch({type: Auth.LOGOUT});
        dispatch({type: Nav.GoToPath, path:"/auth"});
    })
};

exports.Replay = function(ws, gameID){
    var action = {
        type:exports.REPLAY,
        gameInfo: {gameID:gameID}
    };
    ws.jsend(action);
};

