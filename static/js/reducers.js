var Redux = require("redux");
var Actions = require("./actions");
var reactCookie = require("react-cookie");

var count = function(state, action) {
    if (state == undefined){state = 0}

    switch (action.type) {
    case Actions.INCREMENT:
        return state + 1;
    case Actions.DECREMENT:
        return state - 1;
    case Actions.SQUARE:
        return state * state;
    case Actions.ROOT:
        return Math.sqrt(state);
    default:
        return state;
    }
};

var asyncget = function(state, action){
    if (state == undefined){
        state = {
            isFetching:false,
            data: "ur mom"
        };
    }

    switch (action.type){
    case Actions.FETCH:
        nextState = {
            isFetching:true
        };
        return nextState;
    case Actions.GET:
        nextState = {
            isFetching:false,
            data: action.data
        };
        return nextState;
    default:
        return state;
    }
};

var logout = function(state, action){
    if (state == undefined){
        state = {
            loggedIn:false
        };
    };

    var newState = {loggedIn:false};
    switch (action.type){
    case Actions.LOGIN:
        newState.loggedIn = true;
        break;
    case Actions.LOGOUT:
        if (action.success) {
            reactCookie.remove('session');
            action.history.replaceState({ nextPathname: "/"}, '/')
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

var auth = function(state, action){
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
    case Actions.AUTHFETCH:
        newState = {isFetching:true};
        break;
    case Actions.CHANGEUSERNAME:
        newState.username = action.username;
        break;
    case Actions.CHANGEPASSWORD:
        newState.password = action.password;
        break;
    case Actions.CHANGEREPEAT:
        newState.repeat = action.repeat;
        break;
    case Actions.LOGIN:
        newState = {isFetching:false};
        if (action.success) {
            reactCookie.save('session', action.session_id);
            action.history.replaceState({ nextPathname: "/auth"}, '/')
        } else {
            console.log(action.error);
        }
        break;
    case Actions.CREATEACCOUNT:
        newState = {isFetching:false};
        if (!action.success) {
            console.log("failed to login");
        }
        break;
    default:
        return state;
    }

    nextState = Object.assign({}, state, newState);
    return nextState;
}

const rootReducer = Redux.combineReducers({
    count,
    asyncget,
    auth,
    logout
})
exports.rootReducer = rootReducer;

