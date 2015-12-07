"use strict"
var reactCookie = require("react-cookie");
var Config = require("../common").Config;
require("whatwg-fetch");

exports.FETCH = 'FETCH';
exports.FETCHED = 'FETCHED';
exports.CREATE = 'CREATEACCOUNT';

exports.Create = function(dispatch, username, password, repeat){
    dispatch({type:exports.FETCH});

    var url = Config.baseURL + "/create_account"
    var action = {type: exports.CREATE};

    var data = JSON.stringify({
        "username":username,
        "password":password,
        "repeat":repeat     
    });

    fetch(url,{
        method:"post",
        body:data,
        headers: {
            "x-session":reactCookie.load("session") || "",
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
    })
    .then(function(resp){
        return resp.json();
    }).then(function(json){
        action.success = true;
        action.resp = json;
        dispatch(action);
    }).catch(function(err){
        action.success = false;
        action.resp = err;
        dispatch(action);
    })

};

