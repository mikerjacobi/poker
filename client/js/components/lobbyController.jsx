"use strict"

var React = require("react");
var connect = require('react-redux').connect;
var Lobby = require("../actions/lobbyAction");
var Auth = require("../actions/authAction");

class CreateGameForm extends React.Component {
    constructor(props){
        super(props);
        this.changeGameName = this.changeGameName.bind(this);
        this.changeGameType = this.changeGameType.bind(this);
        this.state = {gameType: "", gameName: ""};
    }
    changeGameName(event){
        this.setState({gameName:event.target.value});
    }
    changeGameType(event){
        this.setState({gameType: event.target.value});
    }
    render() {
        var gameName = this.state.gameName; 
        var gameType = this.state.gameType;
            
        return (
            <div> 
                <div className="ui input"><input 
                    type="text"
                    id="gamename_textfield"
                    placeholder="gamename"
                    value={gameName}
                    onChange={this.changeGameName}/> </div>
                <select id="gametype_dropdown" className="ui selection dropdown" onChange={this.changeGameType}>
                    <option value="">Game Type</option>
                    <option value="holdem">Hold Em</option>
                    <option value="highcard">High Card</option>
                </select>
                <button 
                    className="ui primary button"
                    id="create_game_button"
                    onClick={this.props.createGame.bind(this, gameName, gameType)}>
                    Create Game
                </button>
            </div>
        );
    }
};

class JoinGameListing extends React.Component {
    render() {
        var players = [];
        for (var i=0; i<this.props.game.players.length; i++){
            players.push(this.props.game.players[i].name); 
        }
        return(
            <tr onClick={this.props.joinGame.bind(this, this.props.game)} id={this.props.id}>
                <td>{this.props.game.gameName}</td>
                <td>{this.props.game.gameType}</td>
                <td>{players.join(", ")}</td>
            </tr>
        );
    }
};

class GameList extends React.Component {
    render() {
        var gameRows = [];
        var keys = Object.keys(this.props.games);
        for (var i=0; i < keys.length; i++) {
            var key = keys[i]
            var listingID = "game_listing_"+i;
            gameRows.push(
                <JoinGameListing
                    key={this.props.games[key].gameID}
                    id={listingID}
                    game={this.props.games[key]}
                    joinGame={this.props.joinGame}>
                </JoinGameListing>
            );
        }

        return(
            <table id="game_listings" className="ui inverted selectable single line table">
                <thead>
                    <tr>
                        <th>Game Name</th>
                        <th>Type</th>
                        <th>Players</th>
                    </tr>
                </thead>
                <tbody>
                    {gameRows}
                </tbody>
            </table>
        );
    }
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
    createGame(createGameName, createGameType){
        Lobby.Create(
            this.props.dispatch, 
            this.props.wsConnection,
            createGameName,
            createGameType
        );    
    }
    joinGame(game){
        Lobby.Join(
            this.props.dispatch, 
            this.props.wsConnection,
            game
        );    
    }
    render() {
        var data = <div> loading... </div>;
        if (this.props.initialized){
            var games ={games:this.props.games};
            data = <div>
                <CreateGameForm gameName={""} createGame={this.createGame}/>
        
                <br/>

                <GameList 
                    {...games}
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
        wsConnection: state.Auth.wsConnection
    };
}

exports.LobbyController = connect(dataMapper)(LobbyController);
