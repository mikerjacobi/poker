"use strict"

var React = require("react");
var connect = require('react-redux').connect;
var HighCard = require("../actions/highcardAction");
var Lobby = require("../actions/lobbyAction");
var Auth = require("../actions/authAction");

class HighCardMenu extends React.Component {
    render() {
        var players = [];
        for (var i=0; i<this.props.players.length; i++){
            players.push(this.props.players[i].name + ":" + this.props.players[i].chips);
        }
        return(<div> 
                <div className="ui green label"> {this.props.gameName} </div>
                <div className="ui teal label"> {this.props.gameType}  </div>
                <div className="ui blue label"> Players: {players.join(", ")} </div>
                 <br/><br/>
                <button 
                    className="ui tiny button"
                    id="play_game_button"
                    onClick={this.props.play}>
                    Play
                </button>
                <button 
                    className="ui primary tiny button"
                    onClick={this.props.leaveGame}>
                    Leave Game
                </button>
            </div>
        );
    };
};

class HighCardTable extends React.Component {
    render() {
        if (this.props && this.props.error){
            return(<div id="game_div"> {this.props.error} </div>)
        }
        if (!this.props && !this.props.players){
            return(<div id="game_div"> no player object </div>)
        }

        var players = [];
        for (var i=0; i<this.props.players.length; i++){
            players.push(<div> 
                name:{this.props.players[i].game_player.name}  --
                chips: {this.props.players[i].game_player.chips} --
                card: {this.props.players[i].card.display}
            </div>  );
        }
        return(<div id="game_div"> 
            {players}
        </div>);
    };
};


class HighCardController extends React.Component {
    constructor(props){
        super(props);
        this.leaveGame = this.leaveGame.bind(this);
        this.play = this.play.bind(this);
    }
    componentDidMount() {
        Auth.wsConnect(this.props.dispatch, this.props.wsConnection);
    }
    componentWillReceiveProps(nextProps) {
        this.props = nextProps;
        if (!this.props.initialized){
            HighCard.Initialize(this.props.dispatch, this.props.initialized, this.props.params.gameid);
            return(<div> loading... </div>);
        }
    }
    leaveGame(){
        Lobby.Leave(
            this.props.dispatch, 
            this.props.wsConnection,
            this.props.params.gameid
        );    
    }
    play(){
        HighCard.Play(
            this.props.wsConnection, 
            this.props.gameInfo.gameID
        )
    }
    render() {
        if (!this.props.initialized){
            return(<div> loading... </div>);
        }
        return (
            <div>
                <HighCardMenu 
                    {...this.props.gameInfo}
                    leaveGame={this.leaveGame}
                    play={this.play}
                />
                <br/>
                <HighCardTable 
                    {...this.props.gameState}
                />
            </div>
        );
    }
};

var dataMapper = function(state){
    return {
        initialized: state.HighCard.initialized,
        gameInfo: state.HighCard.gameInfo,
        gameState: state.HighCard.gameState,
        wsConnection: state.Auth.wsConnection
    };
}

exports.HighCardController = connect(dataMapper)(HighCardController);
