var Redux = require("redux");
var Math = require("./mathReducer").Math;
var Game = require("./gameReducer").Game;
var Async = require("./asyncReducer").Async;
var Auth = require("./authReducer").Auth;
var Nav = require("./navReducer").Nav;
var LoginCreate = require("./loginCreateReducer").LoginCreate;

const Root = Redux.combineReducers({
    Math,
    Game,
    Async,
    Auth,
    Nav,
    LoginCreate
})
exports.Root = Root;
