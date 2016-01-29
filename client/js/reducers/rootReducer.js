var Redux = require("redux");
var Math = require("./mathReducer").Math;
var Lobby = require("./lobbyReducer").Lobby;
var Holdem = require("./holdemReducer").Holdem;
var HighCard = require("./highcardReducer").HighCard;
var Async = require("./asyncReducer").Async;
var Auth = require("./authReducer").Auth;
var Nav = require("./navReducer").Nav;
var Account = require("./accountReducer").Account;

exports.Root = Redux.combineReducers({
    Math:Math,
    Lobby:Lobby,
    Holdem:Holdem,
    HighCard:HighCard,
    Async:Async,
    Auth:Auth,
    Nav:Nav,
    Account:Account
})
