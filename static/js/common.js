var reactCookie = require("react-cookie");

exports.Config = {
    baseURL:"http://jacobra.com:8004"
};

exports.RequireAuth = function(nextState, replaceState){
    var session = reactCookie.load("session") || "";
    if (session == "" ){
        replaceState({ nextPathname: nextState.location.pathname }, '/auth');
    }
}

exports.GetInitialState = function(){
    return {};
};
