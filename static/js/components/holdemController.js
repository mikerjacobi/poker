"use strict"

var React = require("react");
var connect = require('react-redux').connect;
var Holdem = require("../actions/holdemAction");
var Lobby = require("../actions/lobbyAction");
var Auth = require("../actions/authAction");

class HoldemMenu extends React.Component {
    render() {
        return(<div> 
                <button onClick={this.props.leaveGame}>
                    Leave Game
                </button>
            </div>
        );
    };
};

class HoldemTable extends React.Component {
    render() {
        return(<div> TABLE: {this.props.gameName} </div>);
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
        var g = this.props.game;
        return (
        <div>
            <HoldemMenu 
                gameName={g.gameName} 
                gameID={g.gameID}
                leaveGame={this.leaveGame}
            />

            <HoldemTable 
                gameName={g.gameName} 
                gameID={g.gameID}
            />
        </div>);
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
