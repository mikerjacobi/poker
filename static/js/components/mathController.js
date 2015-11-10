"use strict"

var React = require("react");
var connect =  require('react-redux').connect;
var Math  = require("../actions/mathAction");

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
    increment(){
        Math.Increment(this.props.dispatch);
    }
    decrement(){
        Math.Decrement(this.props.dispatch);
    }
    square(){
        Math.Square(this.props.dispatch);
    }
    sqrt(){
        Math.Sqrt(this.props.dispatch);
    }
    componentWillReceiveProps(nextProps) {
        this.props = nextProps;
    }
    render() {
        return (
            <Counter 
                count={this.props.count} 
                increment={this.increment}
                decrement={this.decrement}
                square={this.square}
                sqrt={this.sqrt}>
            </Counter>
        )
    }
};

var dataMapper = function(state){
    return {
        count: state.Math.count
    };
}

exports.MathController = connect(dataMapper)(MathController);
