exports.INCREMENT = 'INCREMENT'
exports.DECREMENT = 'DECREMENT'
exports.SQUARE = 'SQUARE'
exports.ROOT = 'ROOT'

exports.increment = function(action) {
    action.type = exports.INCREMENT;
    return action
}

exports.decrement = function(action) {
    action.type = exports.DECREMENT;
    return action;
}

exports.square = function(action) {
    action.type = exports.SQUARE;
    return action
}
exports.root = function(action) {
    action.type = exports.ROOT;
    return action
}

exports.Do = function(dispatch, action){
    dispatch(action);
}

