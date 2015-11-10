var Auth = require("../actions/authAction");
var reactCookie = require("react-cookie");

var getInitialState = function(){
    var loggedIn = (reactCookie.load("session") || "") != "";
    return {
        loggedIn:loggedIn
    };
};
exports.Logout = function(state, action){
    if (state == undefined){
        state = getInitialState();
    };

    var newState = {loggedIn:false};
    switch (action.type){
    case Auth.LOGIN:
        newState.loggedIn = true;
        break;
    case Auth.LOGOUT:
        if (action.success) {
            reactCookie.remove('session');
            action.history.replaceState({ nextPathname: "/"}, '/');
            return getInitialState();
        } else {
            console.log(action.error);
        }
        break;
    default:
        return state;
    }
    nextState = Object.assign({}, state, newState);
    return nextState;
}

exports.Auth = function(state, action){
    if (state == undefined){
        state = {
            isFetching:false,
            username:"",
            password:"",
            repeat:""
        }
    }
    var newState = {};

    switch (action.type){
    case Auth.AUTHFETCH:
        newState = {isFetching:true};
        break;
    case Auth.CHANGEUSERNAME:
        newState.username = action.username;
        break;
    case Auth.CHANGEPASSWORD:
        newState.password = action.password;
        break;
    case Auth.CHANGEREPEAT:
        newState.repeat = action.repeat;
        break;
    case Auth.LOGIN:
        newState = {isFetching:false};
        if (action.success) {
            reactCookie.save('session', action.session_id);
            action.history.replaceState({ nextPathname: "/auth"}, '/')
        } else {
            console.log(action.error);
        }
        break;
    case Auth.CREATEACCOUNT:
        newState = {isFetching:false};
        if (!action.success) {
            console.log("failed to login");
        }
        break;
    case Auth.LOGOUT:
        return getInitialState();
    default:
        return state;
    }

    nextState = Object.assign({}, state, newState);
    return nextState;
}

