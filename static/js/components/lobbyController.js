"use strict"

var React = require("react");
var connect = require('react-redux').connect;
var Lobby = require("../actions/lobbyAction");
var Auth = require("../actions/authAction");

class CreateGameForm extends React.Component {
    constructor(props){
        super(props);
        this.changeGameName = this.changeGameName.bind(this);
    }
    changeGameName(event){
        this.setState({gameName:event.target.value});
    }
    render() {
        var gameName = "";
        if (this.state != null){
            gameName = this.state.gameName; 
        }
            
        return (
            <div> 
                <button onClick={this.props.createGame.bind(this, gameName)}>
                    Create Game
                </button>
                <input type="text"
                    placeholder="gamename"
                    value={gameName}
                    onChange={this.changeGameName}/>
            </div>);
    }
};

class JoinGameListing extends React.Component {
    render() {
        return(
            <div>
                <button 
                    onClick={this.props.joinGame.bind(this, this.props.game.game_id)}>
                    Join Game: {this.props.game.game_name}
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
                    key={this.props.games[i].game_id}
                    game={this.props.games[i]}
                    joinGame={this.props.joinGame}>
                </JoinGameListing>
            );
        }
        return(<div> {games} </div>);
    };
};

class LobbyController extends React.Component {
    constructor(props){
        super(props);
        this.createGame = this.createGame.bind(this);
        this.joinGame = this.joinGame.bind(this);
    }
    componentDidMount() {
        Auth.wsConnect(this.props.dispatch, this.props.wsConnection);
        Lobby.Initialize(this.props.dispatch, this.props.initialized);
    }
    componentWillReceiveProps(nextProps) {
        this.props = nextProps;
    }
    createGame(createGameName){
        Lobby.Create(
            this.props.dispatch, 
            this.props.wsConnection,
            createGameName
        );    
    }
    joinGame(gameID){
        Lobby.Join(
            this.props.dispatch, 
            this.props.wsConnection,
            gameID
        );    
    }
    render() {
        var data = <div> loading... </div>;
        if (this.props.initialized){
            data = <div>
                <CreateGameForm gameName={""} createGame={this.createGame}/>
        
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
        initialized: state.Lobby.initialized,
        games: state.Lobby.games,
        //gameFormFetching: state.GameForm.isFetching,
        createGameName: state.Lobby.createGameName,
        wsConnection: state.Auth.wsConnection
    };
}

exports.LobbyController = connect(dataMapper)(LobbyController);
