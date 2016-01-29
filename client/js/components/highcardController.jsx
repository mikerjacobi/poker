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
            players.push(this.props.players[i].name);
        }
        return(<div> 
                <div className="ui green label"> {this.props.gameName} </div>
                <div className="ui teal label"> {this.props.gameType}  </div>
                <div className="ui blue label"> Players: {players.join(", ")} </div>
                <button 
                    className="ui black mini button"
                    onClick={this.props.leaveGame}>
                    Leave Game
                </button>
            </div>
        );
    };
};

class HighCardTable extends React.Component {
    render() {
        var card = "";
        if (this.props && this.props.card){
            card = this.props.card.display;
        }
        return(<div> 
                {card}
        </div>);
    };
};


class HighCardController extends React.Component {
    constructor(props){
        super(props);
        this.leaveGame = this.leaveGame.bind(this);
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
    render() {
        if (!this.props.initialized){
            return(<div> loading... </div>);
        }
        return (
            <div>
                <HighCardMenu 
                    {...this.props.gameInfo}
                    leaveGame={this.leaveGame}
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
