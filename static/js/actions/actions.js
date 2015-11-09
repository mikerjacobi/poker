var reactCookie = require("react-cookie");
require("whatwg-fetch")

//math actions
exports.INCREMENT = 'INCREMENT'
exports.DECREMENT = 'DECREMENT'
exports.SQUARE = 'SQUARE'
exports.ROOT = 'ROOT'

//async get actions
exports.GET = 'GET'
exports.FETCH = 'FETCH'

//auth actions
exports.AUTHFETCH = 'AUTHFETCH'
exports.CHANGEUSERNAME = 'CHANGEUSERNAME'
exports.CHANGEPASSWORD = 'CHANGEPASSWORD'
exports.CHANGEREPEAT = 'CHANGEREPEAT'
exports.LOGIN = 'LOGIN'
exports.LOGOUT = 'LOGOUT'
exports.CREATEACCOUNT = 'CREATEACCOUNT'

exports.increment = function(action) {
    action.type = exports.INCREMENT;
    return action;
}
exports.decrement = function(action) {
    action.type = exports.DECREMENT;
    return action;
}
exports.square = function(action) {
    action.type = exports.SQUARE;
    return action
}
exports.root = function(action) {
    action.type = exports.ROOT;
    return action
}

exports.getA = function(action) {
    action.type = exports.GET;
    action.timeout = 500;
    action.url = "http://jacobra.com:8004/geta";
    return action
}
exports.getB = function(action) {
    action.type = exports.GET;
    action.timeout = 3000;
    action.url = "http://jacobra.com:8004/getb";
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

exports.changeUsername = function(action){
    action.type = exports.CHANGEUSERNAME;
    return action;
}
exports.changePassword = function(action){
    action.type = exports.CHANGEPASSWORD;
    return action;
}
exports.changeRepeat = function(action){
    action.type = exports.CHANGEREPEAT;
    return action;
}
exports.CreateAccount = function(dispatch, username, password, repeat){
    dispatch({type:exports.AUTHFETCH});

    url = "http://jacobra.com:8004/create_account"
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

    var url = "http://jacobra.com:8004/login"
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
    var url = "http://jacobra.com:8004/logout"
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

exports.Do = function(dispatch, action){
    dispatch(action);
}

