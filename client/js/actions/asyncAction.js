"use strict"
var reactCookie = require("react-cookie");
var Config = require("../common").Config;
require("whatwg-fetch");

//async actions
exports.GET = '/async/get'
exports.FETCH = '/async/fetch'
exports.FETCHED = '/async/fetched'

var Get = function(dispatch, url, timeout){
    var action = {type: exports.GET};

    fetch(url, {headers:{"x-session":reactCookie.load("session") || ""}})
    .then(function(resp){
        return resp.json();
    }).then(function(json){
        setTimeout(function(){
            action.data = json.payload.data;
            dispatch(action);
            dispatch({type:exports.FETCHED});
        }, timeout);
    }).catch(function(err){
        action.data = err;
        dispatch(action);
        dispatch({type:exports.FETCHED});
    })
};

exports.GetA = function(dispatch) {
    dispatch({type:exports.FETCH});
    var url = Config.baseURL + "/geta"
    var timeout = 500;
    Get(dispatch, url, timeout);
};

exports.GetB = function(dispatch) {
    dispatch({type:exports.FETCH});
    var url = Config.baseURL + "/getb"
    var timeout = 3000;
    Get(dispatch, url, timeout);
};
