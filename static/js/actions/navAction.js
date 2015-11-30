"use strict"

exports.NEXTPATH = 'NEXTPATH';
exports.GOPATH = 'GOPATH';

exports.GoNextPath = function(dispatch, history){
    var action = {
        type: exports.GOPATH,
        history: history
    };
    dispatch(action);
}

exports.SetNextPath = function(dispatch, nextPath){
    var action = {
        type: exports.NEXTPATH,
        nextPath: nextPath
    };
    dispatch(action);
}
