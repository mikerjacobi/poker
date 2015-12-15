"use strict"

var React = require("react");
var connect = require('react-redux').connect;
var Math = require("../actions/mathAction");
var Auth = require("../actions/authAction");

class Counter extends React.Component {
    render() {
        return(
            <div> 
                Count: {this.props.count} <br/>
                <button onClick={this.props.increment}>Increment</button>
                <button onClick={this.props.decrement}>Decrement</button>
                <button onClick={this.props.square}>Square</button>
                <button onClick={this.props.sqrt}>Sqrt</button>
            </div>
        )}
};

class MathController extends React.Component {
    constructor(props){
        super(props);
        this.increment = this.increment.bind(this)
        this.decrement = this.decrement.bind(this)
        this.square = this.square.bind(this)
        this.sqrt = this.sqrt.bind(this)
    }
    componentDidMount() {
        Auth.wsConnect(this.props.dispatch, this.props.wsConnection);
        Math.Initialize(this.props.dispatch, this.props.initialized);
    }
    increment(){
        Math.Increment(this.props.dispatch, this.props.wsConnection);
    }
    decrement(){
        Math.Decrement(this.props.dispatch, this.props.wsConnection);
    }
    square(){
        Math.Square(this.props.dispatch, this.props.wsConnection);
    }
    sqrt(){
        Math.Sqrt(this.props.dispatch, this.props.wsConnection);
    }
    componentWillReceiveProps(nextProps) {
        this.props = nextProps;
    }
    render() {
        var data = <div> loading... </div>;
        if (this.props.initialized){
            data = <Counter 
                    count={this.props.count} 
                    increment={this.increment}
                    decrement={this.decrement}
                    square={this.square}
                    sqrt={this.sqrt}>
                </Counter>;
        }
        return data;
    }
};

var dataMapper = function(state){
    return {
        count: state.Math.count,
        initialized: state.Math.initialized,
        wsConnection: state.Auth.wsConnection
    };
}

exports.MathController = connect(dataMapper)(MathController);
