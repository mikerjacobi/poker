var ReactDOM = require('react-dom')
var Router = require('react-router').Router
var Route = require('react-router').Route
var Link = require('react-router').Link
var Config = require("./config.js").Config

var GameList = require("./dashboard.js").GameList
var Game = require("./dashboard.js").Game
var CreateGameForm = require("./dashboard.js").CreateGameForm
var GameForms = require("./dashboard.js").GameForms
var LoginCreateForm = require("./loginCreate.js").LoginCreateForm
var LogoutForm = require("./logout.js").LogoutForm



ReactDOM.render(
    <GameForms baseurl={Config.url}/>,
    document.getElementById('game_forms')
);

ReactDOM.render(
  <LoginCreateForm baseurl={Config.url}/>,
  document.getElementById('login_create_form')
);

ReactDOM.render(
    <LogoutForm baseurl={Config.url}/>,
    document.getElementById('logout_form')
);
