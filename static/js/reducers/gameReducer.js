var Game = require("../actions/gameAction");
var Auth = require("../actions/authAction");

var getInitialState = function(){
    return {
        games: [],
        initialized: false
    };
};

exports.Game = function(state, action) {
    if (state == undefined){
        state = getInitialState();
    };
    
    var newState = {}
    switch (action.type) {
    case Game.INIT:
        newState.initialized = true;
        newState.games = action.games;
        break;
    case Auth.LOGIN:
        return getInitialState();
    default:
        return state;
    }

    nextState = Object.assign({}, state, newState);
    return nextState;
};
