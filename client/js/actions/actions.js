var actions = {};

defaultCallback = function(dispatch, msg){
    dispatch(msg);
}

exports.Actions = {
    Register(action, callback){
        if (callback === undefined){
            callback = defaultCallback;
        }
        actions[action] = callback; 
    },
    Call(action){
        if (actions[action] == undefined){
            throw action + " is not a registered action.";
        }
        return actions[action];
    },
    List(){
        console.log(actions);
    }
}


