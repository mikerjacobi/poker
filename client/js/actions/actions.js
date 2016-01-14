var actions = {};

defaultCallback = function(dispatch, msg){
    dispatch(msg);
}

exports.Actions = {
    Register: function(action, callback){
        if (callback === undefined){
            callback = defaultCallback;
        }
        actions[action] = callback; 
    },
    Call: function(action){
        if (actions[action] == undefined){
            throw action + " is not a registered action.";
        }
        return actions[action];
    },
    List: function(){
        console.log(actions);
    }
};

