"use strict"

var React = require("react");
var connect = require('react-redux').connect;
var Game = require("../actions/gameAction");
var Auth = require("../actions/authAction");

class CreateGameForm extends React.Component {
    render() {
        return(
            <div> 
                <button onClick={this.props.createGame}>Create Game</button>
            </div>
        )}
};

class JoinGameListing extends React.Component {
    render() {
        return(
            <div>
                {this.props.gameName} -- 
                <button 
                    value={this.props.gameID}
                    onClick={this.props.joinGame}>
                    Join Game
                </button>
            </div>
        )};
};


class GameList extends React.Component {
    render() {
        var games = [];
        for (var i=0; i < this.props.games.length; i++) {
            games.push(
                <JoinGameListing
                    key={this.props.games[i].gameID}
                    gameID={this.props.games[i].gameID}
                    gameName={this.props.games[i].gameName}
                    joinGame={this.props.joinGame}>
                </JoinGameListing>
            );
        }
        return(<div> {games} </div>);
    };
};

class GameController extends React.Component {
    constructor(props){
        super(props);
        this.createGame = this.createGame.bind(this);
        this.joinGame = this.joinGame.bind(this);
    }
    componentDidMount() {
        Auth.wsConnect(this.props.dispatch, this.props.wsConnection);
        Game.Initialize(this.props.dispatch, this.props.initialized);
    }
    componentWillReceiveProps(nextProps) {
        this.props = nextProps;
    }
    createGame(){
        console.log("create game clicked");
    }
    joinGame(event){
        var gameID = event.target.value;
        console.log("join game clicked: ", gameID);
    }
    render() {
        var data = <div> loading... </div>;
        if (this.props.initialized){
            data = <div>
                <CreateGameForm 
                    //count={this.props.count} 
                    createGame={this.createGame}>
                </CreateGameForm>
        
                <br/>

                <GameList 
                    games={this.props.games}
                    joinGame={this.joinGame}>
                </GameList>
            </div>;
        }
        return data;
    }
};

var dataMapper = function(state){
    return {
        initialized: state.Game.initialized,
        games: state.Game.games,
        wsConnection: state.Auth.wsConnection
    };
}

exports.GameController = connect(dataMapper)(GameController);
