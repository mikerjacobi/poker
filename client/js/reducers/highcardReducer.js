var HighCard = require("../actions/highcardAction");
var Lobby = require("../actions/lobbyAction");
var Auth = require("../actions/authAction");

var getInitialState = function(){
    return {
        initialized: false,
        gameInfo: {
            gameID:"",
            gameName:"",
            gameType:"",
            players:[],
        },
        gameState: {
            players: {},
            error: ""
        }
    };
};

exports.HighCard = function(state, action) {
    if (state == undefined){
        state = getInitialState();
    };
    var newState = state
    newState.error =  "";
    switch (action.type) {
    case HighCard.INIT:
        if (!state.initialized){
          newState = action;
        }
        break;
    case HighCard.UPDATE:
        console.log(action)
        newState = action;
        newState.initialized = true;
        break;
    case HighCard.ERROR:
        newState = action;
        break;
    case Lobby.JOINALERT:
        newState = action;
        break;
    case Lobby.JOIN:
        if (!state.initialized){
          return getInitialState();
        }
    case Lobby.LEAVE:
        return getInitialState();
    case Lobby.LEAVEALERT:
        newState.gameInfo = action.game;
        break;
    case Auth.LOGIN:
        return getInitialState();
    default:
        return state;
    }

    nextState = Object.assign({}, state, newState);
    return nextState;
};
