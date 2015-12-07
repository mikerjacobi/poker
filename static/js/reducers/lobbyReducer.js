var Lobby = require("../actions/lobbyAction");
var Auth = require("../actions/authAction");

var getInitialState = function(){
    return {
        games: [],
        initialized: false
    };
};

exports.Lobby = function(state, action) {
    if (state == undefined){
        state = getInitialState();
    };
    
    var newState = {}
    switch (action.type) {
    case Lobby.INIT:
        newState.initialized = true;
        newState.games = action.games;
        break;
    case Lobby.CREATE:
        newState.games = state.games.slice(0);
        newState.games.push(action.game);
        break;
    case Auth.LOGIN:
        return getInitialState();
    default:
        return state;
    }

    nextState = Object.assign({}, state, newState);
    return nextState;
};
