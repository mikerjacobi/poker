var Redux = require("redux");
var Actions = require("./actions")

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

var auth = function(state, action){
    if (state == undefined){
        state = {
            username:"",
            password:"",
            repeat:""
        }
    }
    var newState = {};

    switch (action.type){
    case Actions.CHANGEUSERNAME:
        newState.username = action.username;
        break;
    case Actions.CHANGEPASSWORD:
        newState.password = action.password;
        break;
    case Actions.CHANGEREPEAT:
        newState.repeat = action.repeat;
        break;
    case Actions.CLICKLOGIN:
        x = 1;
        break;
    case Actions.CLICKCREATEACCOUNT:
        x = 1;
        break;
    default:
        console.log('default');
        return state;
    }

    nextState = Object.assign({}, state, newState);
    console.log('next', nextState);
    return nextState;
}

const rootReducer = Redux.combineReducers({
    count,
    asyncget,
    auth
})
exports.rootReducer = rootReducer;

