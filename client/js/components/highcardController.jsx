"use strict"

var React = require("react");
var connect = require('react-redux').connect;
var HighCard = require("../actions/highcardAction");
var Lobby = require("../actions/lobbyAction");
var Auth = require("../actions/authAction");

class HighCardMenu extends React.Component {
    render() {
        var players = [];
        for (var i=0; i<this.props.game.players.length; i++){
            players.push(this.props.game.players[i].name);
        }
        return(<div> 
                <div className="ui green label"> {this.props.game.gameName} </div>
                <div className="ui teal label"> {this.props.game.gameType}  </div>
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
        return(<div> 
            <div className="ui three column stackable padded middle aligned centered color grid">
                <div className="orange column"></div>
                HIGHCARD place holder 
                <div className="violet column"></div>
            </div>
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
        var game = {game:this.props.game};
        return (
            <div>
                <HighCardMenu 
                    {...game}
                    leaveGame={this.leaveGame}
                />
                <br/>
                <HighCardTable 
                    {...game}
                />
            </div>
        );
    }
};

var dataMapper = function(state){
    return {
        initialized: state.HighCard.initialized,
        game: state.HighCard.game,
        wsConnection: state.Auth.wsConnection
    };
}

exports.HighCardController = connect(dataMapper)(HighCardController);
