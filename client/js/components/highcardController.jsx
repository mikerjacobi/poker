"use strict"

var React = require("react");
var connect = require('react-redux').connect;
var HighCard = require("../actions/highcardAction");
var Lobby = require("../actions/lobbyAction");
var Auth = require("../actions/authAction");

class HighCardMenu extends React.Component {
    render() {
        return(<div> 
                <div className="ui label"> {this.props.username} </div>
                <div className="ui green label"> {this.props.gameName} </div>
                <div className="ui teal label"> {this.props.gameType}  </div>
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

        var actionTo = ""
        var players = [];
        var playerCard = "";
        for (var i=0; i<this.props.players.length; i++){
            //set the player who has action
            if (this.props.players[i].gamePlayer.accountID == this.props.actionTo.accountID){
                actionTo = this.props.players[i].gamePlayer.name;
            }
            if (this.props.accountID == this.props.players[i].gamePlayer.accountID){
                playerCard = this.props.players[i].card.display;
            }
            players.push(<div> 
                {this.props.players[i].gamePlayer.name}: {this.props.players[i].gamePlayer.chips} chips
            </div>  );
        }
        return(<div id="game_div"> 
            {players} <br/>
            <div className="ui teal label"> {playerCard} </div> <br/><br/>
            action to: {actionTo}
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
                    username={this.props.username}
                    accountID={this.props.accountID}
                />
                <br/>
                <HighCardTable 
                    {...this.props.gameState}
                    username={this.props.username}
                    accountID={this.props.accountID}
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
        wsConnection: state.Auth.wsConnection,
        username: state.Auth.username,
        accountID: state.Auth.accountID
    };
}

exports.HighCardController = connect(dataMapper)(HighCardController);
