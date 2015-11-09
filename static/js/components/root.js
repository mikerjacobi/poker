"use strict"

var React = require("react");
var connect =  require('react-redux').connect;
var Actions  = require("../actions/actions");
var Counter = require("./counter.js").Counter


class Root extends React.Component {
    constructor(props){
        super(props);
        this.increment = this.increment.bind(this)
        this.decrement = this.decrement.bind(this)
        this.square = this.square.bind(this)
        this.root = this.root.bind(this)
    }
    increment(){
        Actions.Do(this.props.dispatch, Actions.increment({}));
    }
    decrement(){
        Actions.Do(this.props.dispatch, Actions.decrement({}));
    }
    square(){
        Actions.Do(this.props.dispatch, Actions.square({}));
    }
    root(){
        Actions.Do(this.props.dispatch, Actions.root({}));
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
                root={this.root}>
            </Counter>
        )
    }
};

var dataMapper = function(state){
    return {
        count: state.count
    };
}

exports.Root = connect(dataMapper)(Root);
