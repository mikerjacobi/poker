var Config = require("../common").Config;

//math actions
exports.INCREMENT = 'INCREMENT'
exports.DECREMENT = 'DECREMENT'
exports.SQUARE = 'SQUARE'
exports.SQRT = 'SQRT'
exports.INIT = 'MATHINIT'

exports.Increment = function(dispatch, ws) {
    var action = {type:exports.INCREMENT};
    ws.send(JSON.stringify(action));
}
exports.Decrement = function(dispatch, ws){
    var action = {type:exports.DECREMENT};
    ws.send(JSON.stringify(action));
}
exports.Square = function(dispatch, ws){
    var action = {type:exports.SQUARE};
    ws.send(JSON.stringify(action));
}
exports.Sqrt = function(dispatch, ws){
    var action = {type:exports.SQRT};
    ws.send(JSON.stringify(action));
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
        action.data = err;
        dispatch(action);
    })
};
