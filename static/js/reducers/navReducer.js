var Nav = require("../actions/navAction");
var Config = require("../common").Config;

var getInitialNavState = function(){
    return {
        nextPath:"/",
    };
};

exports.Nav = function(state, action){
    if (state == undefined){
        state = getInitialNavState();
    }
    var newState = {};
    
    switch (action.type){
    case Nav.GOPATH:
        action.history.replaceState({ nextPathname: "/"}, state.nextPath);
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
