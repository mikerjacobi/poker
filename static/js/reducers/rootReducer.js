var Redux = require("redux");
var Math = require("./mathReducer").Math;
var Async = require("./asyncReducer").Async;
var Auth = require("./authReducer").Auth;
var Nav = require("./navReducer").Nav;
var LoginCreate = require("./loginCreateReducer").LoginCreate;

const Root = Redux.combineReducers({
    Math,
    Async,
    Auth,
    Nav,
    LoginCreate
})
exports.Root = Root;
