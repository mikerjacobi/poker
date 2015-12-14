var Config = require("../common").Config;
var Auth = require("./authAction");
var Nav = require("./navAction");
var Actions = require("./actions").Actions;

//holdem actions
exports.INIT = 'HOLDEMINIT'

exports.Initialize = function(dispatch, initialized, gameID){
    if (initialized){return;}
    
    var url = Config.baseURL + "/game/" + gameID
    fetch(url, {headers:{"x-session":reactCookie.load("session") || ""}})
    .then(function(resp){
        if (resp.status != 200){
            throw "status code received: " + resp.status;
        }
        return resp.json();
    }).then(function(json){
        var action = {
            type: exports.INIT,
            game: json.payload
        };
        dispatch(action);
    }).catch(function(err){
        console.log("failed to init holdem game ", err);
        dispatch({type: Auth.LOGOUT});
        dispatch({type: Nav.GoToPath, path:"/auth"});
    })
};

