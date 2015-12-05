"use strict"

exports.NEXTPATH = 'NAVNEXTPATH';
exports.GOPATH = 'NAVGOPATH';
exports.GOTOPATH = 'NAVGOTOPATH';
exports.SETHISTORY = 'NAVSETHISTORY';

exports.SetHistory = function(dispatch, history){
    var action = {
        type: exports.SETHISTORY,
        history: history
    };
    dispatch(action);
}

exports.GoToPath = function(dispatch, path){
    var action = {
        type: exports.GOTOPATH,
        path: path
    };
    dispatch(action);
}

exports.GoNextPath = function(dispatch){
    var action = {
        type: exports.GOPATH
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
