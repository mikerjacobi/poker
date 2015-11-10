"use strict"
var reactCookie = require("react-cookie");
var Config = require("../common").Config;
require("whatwg-fetch");

//async get actions
exports.GET = 'GET'
exports.FETCH = 'FETCH'

exports.GetA = function() {
    var action = {
        type:exports.GET,
        timeout:500,
        url:Config.baseURL + "/geta"
    };
    return action
}
exports.GetB = function() {
    var action = {
        type:exports.GET,
        timeout:3000,
        url:Config.baseURL + "/getb"
    };
    return action
}

exports.Get = function(dispatch, action){
    dispatch({type:exports.FETCH});
    fetch(action.url, {headers:{"x-session":reactCookie.load("session") || ""}})
    .then(function(resp){
        return resp.json();
    }).then(function(json){
        action.data = json.payload.data;
        setTimeout(function(){
            dispatch(action);
        }, action.timeout);
    }).catch(function(err){
        action.data = err;
        dispatch(action);
    })
}


