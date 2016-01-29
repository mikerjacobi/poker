var HighCard = require("../actions/highcardAction");
var Lobby = require("../actions/lobbyAction");
var Auth = require("../actions/authAction");

var getInitialState = function(){
    return {
        initialized: false,
        game: {
            gameID:"",
            gameName:"",
            gameType:"",
            players:[]
        }
    };
};

exports.HighCard = function(state, action) {
    if (state == undefined){
        state = getInitialState();
    };
    
    var newState = {}
    switch (action.type) {
    case HighCard.INIT:
        newState.initialized = true;
        newState.game = {
            gameID: action.game.gameID,
            gameName: action.game.gameName,
            gameType: action.game.gameType,
            players: action.game.players
        };
        break;
    case Lobby.JOINALERT:
        newState.game = {
            gameID: action.game.gameID,
            gameName: action.game.gameName,
            gameType: action.game.gameType,
            players: action.game.players
        };
        break;
    case Lobby.JOIN:
        return getInitialState();
    case Lobby.LEAVE:
        return getInitialState();
    case Lobby.LEAVEALERT:
        newState.game = {
            gameID: action.game.gameID,
            gameName: action.game.gameName,
            gameType: action.game.gameType,
            players: action.game.players
        };
        break;
    case Auth.LOGIN:
        return getInitialState();
    default:
        return state;
    }

    nextState = Object.assign({}, state, newState);
    return nextState;
};
