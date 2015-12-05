"use strict"

var React = require("react");
var connect = require('react-redux').connect;
var Game = require("../actions/gameAction");
var Auth = require("../actions/authAction");

class CreateGameForm extends React.Component {
    render() {
        return (
            <div> 
                <input type="text"
                    placeholder="gamename"
                    value={this.props.createGameName} 
                    onChange={this.props.changeGameName}/>
                <button onClick={this.props.createGame}>Create Game</button>
            </div>);
    }
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
            console.log(this.props.games[i])
            games.push(
                <JoinGameListing
                    key={this.props.games[i].game_id}
                    gameID={this.props.games[i].game_id}
                    gameName={this.props.games[i].game_name}
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
        this.changeGameName = this.changeGameName.bind(this);
        this.joinGame = this.joinGame.bind(this);
    }
    componentDidMount() {
        Auth.wsConnect(this.props.dispatch, this.props.wsConnection);
        Game.Initialize(this.props.dispatch, this.props.initialized);
    }
    componentWillReceiveProps(nextProps) {
        this.props = nextProps;
    }
    changeGameName(event){
        var gameName = event.target.value;
        Game.ChangeGameName(this.props.dispatch, gameName);
    }
    createGame(event){
        Game.Create(
            this.props.dispatch, 
            this.props.wsConnection,
            this.props.createGameName
        );    
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
                    createGameName={this.props.createGameName}
                    //isFetching={this.props.gameFormFetching}
                    createGame={this.createGame}
                    changeGameName={this.changeGameName}>
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
        //gameFormFetching: state.GameForm.isFetching,
        createGameName: state.Game.createGameName,
        wsConnection: state.Auth.wsConnection
    };
}

exports.GameController = connect(dataMapper)(GameController);
