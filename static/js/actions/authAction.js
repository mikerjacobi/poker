"use strict"
var reactCookie = require("react-cookie");
var Config = require("../common").Config;
require("whatwg-fetch");

//auth actions
exports.AUTHFETCH = 'AUTHFETCH'
exports.CHANGEUSERNAME = 'CHANGEUSERNAME'
exports.CHANGEPASSWORD = 'CHANGEPASSWORD'
exports.CHANGEREPEAT = 'CHANGEREPEAT'
exports.LOGIN = 'LOGIN'
exports.LOGOUT = 'LOGOUT'
exports.CREATEACCOUNT = 'CREATEACCOUNT'

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
    dispatch({type:exports.AUTHFETCH});

    url = Config.baseURL + "/create_account"
    action = {type: exports.CREATEACCOUNT};

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

exports.Login = function(dispatch, username, password, history){
    dispatch({type:exports.AUTHFETCH});

    var url = Config.baseURL + "/login"
    var action = {type: exports.LOGIN, history:history, success:false};

    var data = JSON.stringify({
        "username":username,
        "password":password      
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
        if (resp.status == 200){
            action.success = true;
        } else if (resp.status == 400){
            action.error = "bad login input";
        } else if (resp.status == 401){
            action.error = "bad login creds";
        } else {
            action.error = "login unknown";
        }
        return resp.json();
    }).then(function(json){
        if (action.success == true){
            action.session_id = json.payload.session_id;
        }
        dispatch(action);
    }).catch(function(err){
        action.error = err;
        dispatch(action);
    })
}

exports.Logout = function(dispatch, history){
    var url = Config.baseURL + "/logout"
    var action = {type: exports.LOGOUT, history:history, success:false};

    fetch(url,{
        method:"post",
        headers: {
            "x-session":reactCookie.load("session") || "",
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
    })
    .then(function(resp){
        if (resp.status == 200){
            action.success = true;
        } else if (resp.status == 401){
            action.error = "bad logout creds";
        } else {
            action.error = "logout unknown";
        }
        dispatch(action);
    }).catch(function(err){
        action.error = err;
        dispatch(action);
    })
}
