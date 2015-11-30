var reactCookie = require("react-cookie");
var Nav = require("./actions/navAction");

exports.Config = {
    baseURL:"http://jacobra.com:8004",
    wsURL:"ws://jacobra.com:8004/ws"
};

exports.RequireAuth = function(store) {
    return function(nextState, replaceState){
        var session = reactCookie.load("session") || "";
        if (session == "" ){
            replaceState({ nextPathname: nextState.location.pathname }, '/auth');
        }
    };
}
exports.SetPath = function(store) {
    return function(event){
        Nav.SetNextPath(store.dispatch, event.target.hash.replace("#",""));
    };
}

exports.GetInitialState = function(){
    return {};
};
    
