var Async = require("../actions/asyncAction");

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

    var newState = {}
    switch (action.type){
    case Async.FETCH:
        newState = {isFetching:true};
        break;
    case Async.FETCHED:
        newState = {isFetching:false};
        break;
    case Async.GET:
        newState = {data: action.data};
        break;
    default:
        return state;
    }

    nextState = Object.assign({}, state, newState);
    return nextState;
};

