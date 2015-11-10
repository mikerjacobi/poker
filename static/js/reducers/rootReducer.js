var Redux = require("redux");
var Math = require("./mathReducer").Math;
var Async = require("./asyncReducer").Async;
var Logout = require("./authReducer").Logout;
var Auth = require("./authReducer").Auth;

const Root = Redux.combineReducers({
    Math,
    Async,
    Auth,
    Logout
})
exports.Root = Root;
