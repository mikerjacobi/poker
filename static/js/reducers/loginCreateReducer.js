var LoginCreate = require("../actions/loginCreateAction");

var getInitialLoginCreateState = function(){
    var loggedIn = (reactCookie.load("session") || "") != "";
    return {
        isFetching:false,
        username:"",
        password:"",
        repeat:"",
        loggedIn:loggedIn,
        wsConnection: false
    };
};

exports.LoginCreate = function(state, action){
    if (state == undefined){
        state = getInitialLoginCreateState();
    }
    var newState = {};

    switch (action.type){
    case LoginCreate.FETCH:
        newState = {isFetching:true};
        break;
    case LoginCreate.FETCHED:
        newState = {isFetching:false};
        break;
    case LoginCreate.CHANGEUSERNAME:
        newState.username = action.username;
        break;
    case LoginCreate.CHANGEPASSWORD:
        newState.password = action.password;
        break;
    case LoginCreate.CHANGEREPEAT:
        newState.repeat = action.repeat;
        break;
    case LoginCreate.CREATEACCOUNT:
        newState = {isFetching:false};
        if (!action.success) {
            console.log("failed to login");
        }
        break;
    case LoginCreate.CLEARFORM:
        newState = {
            password:"",
            repeat:""
        }
        if (action.clearUser){
            newState.username = "";
        }
        break;
    default:
        return state;
    }

    nextState = Object.assign({}, state, newState);
    return nextState;
}


