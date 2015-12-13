var Holdem = require("../actions/holdemAction");
var Lobby = require("../actions/lobbyAction");
var Auth = require("../actions/authAction");

var getInitialState = function(){
    return {
        initialized: false,
        game: {
            gameID:"",
            gameName:"",
            players:[]
        }
    };
};

exports.Holdem = function(state, action) {
    if (state == undefined){
        state = getInitialState();
    };
    
    var newState = {}
    switch (action.type) {
    case Holdem.INIT:
        newState.initialized = true;
        newState.game = {
            gameID: action.game.gameID,
            gameName: action.game.gameName,
            players: action.game.players
        };
        break;
    case Lobby.JOIN:
        return getInitialState();
    case Lobby.LEAVE:
        return getInitialState();
    case Auth.LOGIN:
        return getInitialState();
    default:
        return state;
    }

    nextState = Object.assign({}, state, newState);
    return nextState;
};
