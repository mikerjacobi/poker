
//math actions
exports.INCREMENT = 'INCREMENT'
exports.DECREMENT = 'DECREMENT'
exports.SQUARE = 'SQUARE'
exports.SQRT = 'SQRT'

exports.Increment = function(dispatch) {
    var action = {type:exports.INCREMENT};
    dispatch(action);
}
exports.Decrement = function(dispatch) {
    var action = {type:exports.DECREMENT};
    dispatch(action);
}
exports.Square = function(dispatch) {
    var action = {type:exports.SQUARE};
    dispatch(action);
}
exports.Sqrt = function(dispatch) {
    var action = {type:exports.SQRT};
    dispatch(action);
}
