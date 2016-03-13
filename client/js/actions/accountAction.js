"use strict"
var reactCookie = require("react-cookie");
var Config = require("../common").Config;
require("whatwg-fetch");
var Actions = require("./actions").Actions;

//account actions
exports.FETCH = '/account/fetch';
exports.FETCHED = '/account/fetched';
exports.CREATE = '/account/create';
exports.LOAD = '/account/load';
exports.REQUESTCHIPS = '/account/chips/request'

//account registrations
Actions.Register(exports.LOAD);

exports.Init = function(dispatch, ws){
    var action = {type:exports.LOAD};
    ws.jsend(action);
};

exports.RequestChips = function(dispatch, ws, amount){
    var action = {type:exports.REQUESTCHIPS, amount};
    ws.jsend(action);
};

exports.Create = function(dispatch, username, password, repeat){
    dispatch({type:exports.FETCH});

    var url = Config.baseURL + "/account"
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

