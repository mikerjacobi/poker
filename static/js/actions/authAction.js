"use strict"
var reactCookie = require("react-cookie");
var Config = require("../common").Config;
var Actions = require("./actions").Actions;
var Nav = require("./navAction")
var Async = require("./asyncAction")
var Account = require("./accountAction")
require("whatwg-fetch");

//auth actions
exports.LOGIN = 'LOGIN';
exports.LOGOUT = 'LOGOUT';
exports.WSCONNECT = 'WSCONNECT';
exports.WSDISCONNECT = 'WSDISCONNECT';
exports.WSERROR = 'WSERROR';

Actions.Register(exports.WSCONNECT);
Actions.Register(exports.WSDISCONNECT);
Actions.Register(exports.WSERROR, function(dispatch, msg){
    dispatch({type:Async.FETCHED});
    dispatch({type:Account.FETCHED});
    console.log("server returned ws error: ", msg);
});

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

    wsConnection.jsend = function(obj){
        try{
            var json = JSON.stringify(obj);
        } catch(err){
            console.log("failed to jsonmarshal: "+obj);
            return;
        }
        wsConnection.send(json);
    }

    wsConnection.onopen = function () {
        wsConnection.jsend({type:exports.WSCONNECT}); 
    };

    wsConnection.onclose = function(){
        console.log("ws closed")
    };

    wsConnection.onerror = function (event) {
        Actions.Call(exports.WSERROR)(dispatch, event.data);
    };

    wsConnection.onmessage = function(event) {
        try{
            var msg = JSON.parse(event.data);
        } catch(err){
            console.log("JSON parse error: ", err);
            console.log("failing event.data: ", event.data);
            return;
        }

        try{
            Actions.Call(msg.type)(dispatch, msg);
        } catch(err){
            console.log("failed to Actions.Call: " + err);
            console.log("failing message: ", msg);
            return;
        }
    };
    dispatch(action);
    return true;
}

exports.wsDisconnect = function(dispatch, wsConnection){
    if (wsConnection == false){
        return;
    }
    wsConnection.jsend({type:exports.WSDISCONNECT});
    wsConnection.close();
    dispatch({type:exports.WSDISCONNECT});
}

exports.Login = function(dispatch, username, password, wsConn, history){
    dispatch({type:Account.FETCH});

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
        if (resp.status != 200){
           throw "received bad status code from login: ", resp.status;
        } 
        return resp.json();
    }).then(function(json){
        action.sessionID = json.payload.sessionID;
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
    dispatch({type:Account.FETCHED});
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
