"use strict"

var React = require("react");
var connect = require('react-redux').connect;
var Holdem = require("../actions/holdemAction");
var Auth = require("../actions/authAction");

class HoldemMenu extends React.Component {
    render() {
        return(<div> MENU: {this.props.gameID} </div>);
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
    }
    componentDidMount() {
        Auth.wsConnect(this.props.dispatch, this.props.wsConnection);
        Holdem.Initialize(this.props.dispatch, this.props.initialized, this.props.params.gameid);
    }
    componentWillReceiveProps(nextProps) {
        this.props = nextProps;
    }
    /*joinGame(gameID){
        Lobby.Join(
            this.props.dispatch, 
            this.props.wsConnection,
            gameID
        );    
    }*/
    render() {
        if (!this.props.initialized){
            return(<div> loading... </div>);
        }
        var g = this.props.game;
        return (
        <div>
            <HoldemMenu gameName={g.gameName} gameID={g.gameID}/>
            <HoldemTable gameName={g.gameName} gameID={g.gameID}/>
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
