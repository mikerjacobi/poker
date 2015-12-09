var Config = require("../common").Config;
var Auth = require("./authAction");
var Nav = require("./navAction");
var Actions = require("./actions").Actions;

//math actions
exports.INCREMENT = 'INCREMENT'
exports.DECREMENT = 'DECREMENT'
exports.SQUARE = 'SQUARE'
exports.SQRT = 'SQRT'
exports.INIT = 'MATHINIT'

Actions.Register(exports.INCREMENT);
Actions.Register(exports.DECREMENT);
Actions.Register(exports.SQUARE);
Actions.Register(exports.SQRT);

exports.Increment = function(dispatch, ws) {
    var action = {type:exports.INCREMENT};
    ws.jsend(action);
}
exports.Decrement = function(dispatch, ws){
    var action = {type:exports.DECREMENT};
    ws.jsend(action);
}
exports.Square = function(dispatch, ws){
    var action = {type:exports.SQUARE};
    ws.jsend(action);
}
exports.Sqrt = function(dispatch, ws){
    var action = {type:exports.SQRT};
    ws.jsend(action);
}
exports.Initialize = function(dispatch, initialized){
    if (initialized){return;}
    
    var action = {type: exports.INIT};
    var url = Config.baseURL + "/math"
    fetch(url, {headers:{"x-session":reactCookie.load("session") || ""}})
    .then(function(resp){
        return resp.json();
    }).then(function(json){
        action.count = json.payload.count;
        dispatch(action);
    }).catch(function(err){
        console.log("failed to init math ", err);
        dispatch({type: Auth.LOGOUT});
        dispatch({type: Nav.GoToPath, path:"/auth"});
    })
};
