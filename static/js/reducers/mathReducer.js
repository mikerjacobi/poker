var MathActions = require("../actions/mathAction");
var Auth = require("../actions/authAction");

var getInitialState = function(){
    return {
        count: 0,
        initialized: false
    };
};

exports.Math = function(state, action) {
    if (state == undefined){
        state = getInitialState();
    };
    
    var newState = {}
    switch (action.type) {
    case MathActions.INIT:
        newState.count = action.count;
        newState.initialized = true;
        break;
    case MathActions.INCREMENT:
        newState.count = state.count + 1;
        break;
    case MathActions.DECREMENT:
        newState.count = state.count - 1;
        break;
    case MathActions.SQUARE:
        newState.count = state.count * state.count;
        break;
    case MathActions.SQRT:
        newState.count = Math.sqrt(state.count);
        break;
    case Auth.LOGIN:
        return getInitialState();
    default:
        return state;
    }

    nextState = Object.assign({}, state, newState);
    return nextState;
};
