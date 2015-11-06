require("whatwg-fetch")
exports.INCREMENT = 'INCREMENT'
exports.DECREMENT = 'DECREMENT'
exports.SQUARE = 'SQUARE'
exports.ROOT = 'ROOT'
exports.GET = 'GET'
exports.FETCH = 'FETCH'

exports.increment = function(action) {
    action.type = exports.INCREMENT;
    return action
}

exports.decrement = function(action) {
    action.type = exports.DECREMENT;
    return action;
}

exports.square = function(action) {
    action.type = exports.SQUARE;
    return action
}
exports.root = function(action) {
    action.type = exports.ROOT;
    return action
}
exports.getA = function(action) {
    action.type = exports.GET;
    action.timeout = 500;
    action.url = "http://jacobra.com:8004/geta";
    return action
}
exports.getB = function(action) {
    action.type = exports.GET;
    action.timeout = 3000;
    action.url = "http://jacobra.com:8004/getb";
    return action
}
exports.Get = function(dispatch, action){
    dispatch({type:exports.FETCH});
    fetch(action.url)
    .then(function(resp){
        return resp.json();
    }).then(function(json){
        action.data = json.payload.data;
        setTimeout(function(){
            dispatch(action);
        }, action.timeout);
    }).catch(function(err){
        action.data = err;
        dispatch(action);
    })
}
exports.Do = function(dispatch, action){
    dispatch(action);
}

