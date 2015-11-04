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
}

const rootReducer = Redux.combineReducers({
    count
})
exports.rootReducer = rootReducer;

