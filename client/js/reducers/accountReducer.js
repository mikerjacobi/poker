var Account = require("../actions/accountAction");

var getInitialAccountState = function(){
    return {
        account: {
            username:"",
            balance: -1
        },
        isFetching:false
    };
};

exports.Account = function(state, action){
    if (state == undefined){
        state = getInitialAccountState();
    }
    var newState = {};

    switch (action.type){
    case Account.LOAD:
        newState = action;
        break;
    case Account.FETCH:
        newState = {isFetching:true};
        break;
    case Account.FETCHED:
        newState = {isFetching:false};
        break;
    case Account.CREATE:
        newState = {isFetching:false};
        if (!action.success) {
            console.log("failed to login");
        }
        break;
    default:
        return state;
    }

    nextState = Object.assign({}, state, newState);
    return nextState;
}


