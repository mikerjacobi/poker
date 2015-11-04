"use strict"
var React = require("react");

class Counter extends React.Component {
    render() {
        return(
            <div> 
                Count: {this.props.count} <br/>
                <button onClick={this.props.increment}>Increment</button>
                <button onClick={this.props.decrement}>Decrement</button>
                <button onClick={this.props.square}>Square</button>
                <button onClick={this.props.root}>Root</button>
            </div>
        )}
};
exports.Counter = Counter;
