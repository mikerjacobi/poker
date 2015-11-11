var reactCookie = require("react-cookie");
var Auth = require("./actions/authAction");

exports.Config = {
    baseURL:"http://jacobra.com:8004"
};

exports.RequireAuth = function(store) {
    return function(nextState, replaceState){
        var session = reactCookie.load("session") || "";
        if (session == "" ){
            Auth.SetNextPath(store.dispatch, nextState.location.pathname);
            replaceState({ nextPathname: nextState.location.pathname }, '/auth');
        }
    };
}

exports.GetInitialState = function(){
    return {};
};
