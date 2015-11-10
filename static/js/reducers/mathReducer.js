var MathAction = require("../actions/mathAction");
var Auth = require("../actions/authAction");

var getInitialState = function(){
    return {
        count: 0
    };
};

exports.Math = function(state, action) {
    if (state == undefined){
        state = getInitialState();
    };

    switch (action.type) {
    case MathAction.INCREMENT:
        return {count:state.count + 1};
    case MathAction.DECREMENT:
        return {count:state.count - 1};
    case MathAction.SQUARE:
        return {count:state.count * state.count};
    case MathAction.ROOT:
        return {count:Math.sqrt(state.count)};
    case Auth.LOGOUT:
        return getInitialState();
    default:
        return state;
    }
};
