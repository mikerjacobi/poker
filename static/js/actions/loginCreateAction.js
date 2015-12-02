"use strict"
var reactCookie = require("react-cookie");
var Config = require("../common").Config;
var Nav = require("./navAction")
require("whatwg-fetch");

exports.FETCH = 'FETCH';
exports.FETCHED = 'FETCHED';
exports.CHANGEUSERNAME = 'CHANGEUSERNAME';
exports.CHANGEPASSWORD = 'CHANGEPASSWORD';
exports.CHANGEREPEAT = 'CHANGEREPEAT';
exports.CREATEACCOUNT = 'CREATEACCOUNT';
exports.CLEARFORM = 'CLEARFORM';

exports.ChangeUsername = function(dispatch, username){
    var action = {
        type:exports.CHANGEUSERNAME,
        username:username
    }
    dispatch(action);
}
exports.ChangePassword = function(dispatch, password){
    var action = {
        type:exports.CHANGEPASSWORD,
        password:password
    }
    dispatch(action);
}
exports.ChangeRepeat = function(dispatch, repeat){
    var action = {
        type:exports.CHANGEREPEAT,
        repeat:repeat
    }
    dispatch(action);
}
exports.CreateAccount = function(dispatch, username, password, repeat){
    dispatch({type:exports.FETCH});

    var url = Config.baseURL + "/create_account"
    var action = {type: exports.CREATEACCOUNT};

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

}

exports.ClearForm = function(dispatch, clearUser){
    var action = {
        type: exports.CLEARFORM,
        clearUser: clearUser
    };
    dispatch(action);
}

