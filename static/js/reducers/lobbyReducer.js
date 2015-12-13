var Lobby = require("../actions/lobbyAction");
var Auth = require("../actions/authAction");

var getInitialState = function(){
    return {
        games: {},
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
    case Lobby.JOIN:
        var newGames = Object.assign({}, state.games, {});
        newGames[action.game.gameID] = action.game;
        newState.games = newGames;
        break;
    case Lobby.CREATE:
        var newGames = Object.assign({}, state.games, {});
        newGames[action.game.gameID] = action.game;
        newState.games = newGames;
        break;
    case Lobby.LEAVE:
        var newGames = Object.assign({}, state.games, {});
        newGames[action.game.gameID] = action.game;
        newState.games = newGames;
        break;
    case Auth.LOGIN:
        return getInitialState();
    default:
        return state;
    }

    return Object.assign({}, state, newState);
};
