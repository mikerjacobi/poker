var Async = require("../actions/asyncAction");
var Auth = require("../actions/authAction");

var getInitialState = function(){
    return {
        isFetching:false,
        data: "ur mom"
    };
};

exports.Async = function(state, action){
    if (state == undefined){
        state = getInitialState();
    }

    switch (action.type){
    case Async.FETCH:
        nextState = {
            isFetching:true
        };
        return nextState;
    case Async.GET:
        nextState = {
            isFetching:false,
            data: action.data
        };
        return nextState;
    case Auth.LOGOUT:
        return getInitialState();
    default:
        return state;
    }
};

