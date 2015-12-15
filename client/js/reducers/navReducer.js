var Nav = require("../actions/navAction");
var Config = require("../common").Config;

var getInitialNavState = function(){
    return {
        nextPath:"/",
        history:false
    };
};

exports.Nav = function(state, action){
    if (state == undefined){
        state = getInitialNavState();
    }
    var newState = {};
    
    switch (action.type){
    case Nav.SETHISTORY:
        newState.history = action.history;
    case Nav.GOPATH:
        if (state.history){
            state.history.replaceState({ nextPathname: "/"}, state.nextPath);
        }
        break;
    case Nav.GOTOPATH:
        if (state.history){
            state.history.replaceState({ nextPathname: "/"}, action.path);
        }
        break;
    case Nav.NEXTPATH:
        newState = {nextPath:action.nextPath};
        break;
    default:
        return state;
    }

    nextState = Object.assign({}, state, newState);
    return nextState;
}
