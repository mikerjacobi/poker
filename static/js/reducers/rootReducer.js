var Redux = require("redux");
var Math = require("./mathReducer").Math;
var Lobby = require("./lobbyReducer").Lobby;
var Async = require("./asyncReducer").Async;
var Auth = require("./authReducer").Auth;
var Nav = require("./navReducer").Nav;
var Account = require("./accountReducer").Account;

const Root = Redux.combineReducers({
    Math,
    Lobby,
    Async,
    Auth,
    Nav,
    Account
})
exports.Root = Root;
