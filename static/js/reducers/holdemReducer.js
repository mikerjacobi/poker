var Holdem = require("../actions/holdemAction");
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
            gameID: action.game.game_id,
            gameName: action.game.game_name,
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
