var Auth = require("../actions/authAction");
var reactCookie = require("react-cookie");
var Config = require("../common").Config;

var getInitialAuthState = function(){
    var loggedIn = (reactCookie.load("session") || "") != "";
    return {
        isFetching:false,
        username:"",
        accountID:"",
        password:"",
        repeat:"",
        loggedIn:loggedIn,
        wsConnection: false,
    };
};

exports.Auth = function(state, action){
    if (state == undefined){
        state = getInitialAuthState();
    }
    var newState = {};

    switch (action.type){
    case Auth.LOGIN:
        reactCookie.save('session', action.sessionID);
        newState.loggedIn = true;
        newState.username = action.username;
        newState.accountID = action.accountID;
        break;
    case Auth.LOGOUT:
        reactCookie.remove('session');
        return getInitialAuthState();
    case Auth.WSCONNECT:
        newState = {
            wsConnection:action.wsConnection,
        };
        break;
    case Auth.WSDISCONNECT:
        newState = {wsConnection:false};
        break;
    case Auth.WSINFO:
        newState = {
            username: action.username,
            accountID: action.accountID
        };
        break;
    default:
        return state;
    }

    nextState = Object.assign({}, state, newState);
    return nextState;
}

