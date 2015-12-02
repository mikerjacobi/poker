"use strict"
var reactCookie = require("react-cookie");
var Config = require("../common").Config;
var Nav = require("./navAction")
var LoginCreate = require("./loginCreateAction")
require("whatwg-fetch");

//auth actions
exports.LOGIN = 'LOGIN';
exports.LOGOUT = 'LOGOUT';
exports.WSCONNECT = 'WSCONNECT';
exports.WSDISCONNECT = 'WSDISCONNECT';

exports.wsConnect = function(dispatch, currentWSConnection){
    var loggedIn = (reactCookie.load("session") || "") != "";

    //only attempt to wsconnect if we are logged in and dont have a ws conn
    if (!loggedIn || currentWSConnection){
        return false;
    }

    var action = {type:exports.WSCONNECT};
    try{
        var wsConnection = new WebSocket(Config.wsURL); 
    } catch(err){
        return false;
    }
    action.wsConnection = wsConnection;

    wsConnection.onopen = function () {
        wsConnection.send('ping'); 
    };

    wsConnection.onerror = function (error) {
        console.log('ws error: ' + error);
    };

    wsConnection.onmessage = function(event) {
        try{
            var msg = JSON.parse(event.data);
        } catch(err){
            console.log(event.data);
            console.log("JSON parse error: ", err);
            return;
        }
        dispatch(msg);
    };
    dispatch(action);
    return true;
}

exports.wsDisconnect = function(dispatch, wsConnection){
    if (wsConnection == false){
        return;
    }
    wsConnection.close();
    dispatch({type:exports.WSDISCONNECT});
}


exports.Login = function(dispatch, username, password, wsConn, history){
    dispatch({type:LoginCreate.FETCH});

    var url = Config.baseURL + "/login"
    var action = {type: exports.LOGIN};

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
        var err = "";
        if (resp.status == 200){
            LoginCreate.ClearForm(dispatch, true);
            return resp.json();
        } else if (resp.status == 400){
            err = "bad login input";
        } else if (resp.status == 401){
            err = "bad login creds";
        } else {
            err = "login unknown";
        }
        LoginCreate.ClearForm(dispatch, false);
        throw err;
    }).then(function(json){
        action.session_id = json.payload.session_id;
        dispatch(action);

        //if we create a new ws, gonextpath
        if (exports.wsConnect(dispatch, wsConn)){
            Nav.GoNextPath(dispatch, history);
        } else {
            exports.Logout(dispatch, false)
        }
    }).catch(function(err){
        console.log(err);
    })
    dispatch({type:LoginCreate.FETCHED});
}

exports.Logout = function(dispatch, wsConn){
    Nav.SetNextPath(dispatch, "/auth");

    var url = Config.baseURL + "/logout"
    fetch(url, {
        method:"post",
        headers: {
            "x-session":reactCookie.load("session") || "",
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
    })
    .then(function(resp){
        var err = "";
        if (resp.status == 200){
            //close web socket connection
            exports.wsDisconnect(dispatch, wsConn);
            dispatch({type: exports.LOGOUT});
            return;
        } else if (resp.status == 401){
            err = "bad logout creds";
        } else {
            err = "logout unknown";
        }
        throw err;
    }).catch(function(err){
        //close web socket connection
        exports.wsDisconnect(dispatch, wsConn);
        dispatch({type: exports.LOGOUT});
        console.log(err);
    })
}
