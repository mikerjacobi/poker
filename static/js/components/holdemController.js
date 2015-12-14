"use strict"

var React = require("react");
var connect = require('react-redux').connect;
var Holdem = require("../actions/holdemAction");
var Lobby = require("../actions/lobbyAction");
var Auth = require("../actions/authAction");

class HoldemMenu extends React.Component {
    render() {
        var players = [];
        for (var i=0; i<this.props.game.players.length; i++){
            players.push(this.props.game.players[i].name);
        }
        return(<div> 
                {this.props.game.gameName} -
                {this.props.game.gameType} - 
                {players.join(", ")} - 
                <button onClick={this.props.leaveGame}>
                    Leave Game
                </button>
            </div>
        );
    };
};

class HoldemTable extends React.Component {
    render() {
        return(<div> game place holder </div>);
    };
};

class HoldemController extends React.Component {
    constructor(props){
        super(props);
        this.leaveGame = this.leaveGame.bind(this);
    }
    componentDidMount() {
        Auth.wsConnect(this.props.dispatch, this.props.wsConnection);
        Holdem.Initialize(this.props.dispatch, this.props.initialized, this.props.params.gameid);
    }
    componentWillReceiveProps(nextProps) {
        this.props = nextProps;
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
                <HoldemMenu 
                    {...game}
                    leaveGame={this.leaveGame}
                />

                <HoldemTable 
                    {...game}
                />
            </div>
        );
    }
};

var dataMapper = function(state){
    return {
        initialized: state.Holdem.initialized,
        game: state.Holdem.game,
        wsConnection: state.Auth.wsConnection
    };
}

exports.HoldemController = connect(dataMapper)(HoldemController);
