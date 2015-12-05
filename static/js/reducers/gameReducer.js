var Game = require("../actions/gameAction");
var Auth = require("../actions/authAction");

var getInitialState = function(){
    return {
        games: [],
        initialized: false,
        createGameName:""
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
    case Game.CREATE:
        console.log(action);
        newState.games = state.games.slice(0);
        newState.games.push(action.game);
    case Game.CHANGEGAMENAME:
        newState.createGameName = action.createGameName;
        break;
    case Auth.LOGIN:
        return getInitialState();
    default:
        return state;
    }

    nextState = Object.assign({}, state, newState);
    return nextState;
};
