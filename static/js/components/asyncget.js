"use strict"
var React = require("react");
var connect =  require('react-redux').connect;
var Actions = require("../actions")

class GetWidget extends React.Component {
    render(){
        var data = <div>loading...</div>;
        if (!this.props.isFetching){
            data = <div> Data: {this.props.data}  </div>;
        }
        return(
            <div>   
                { data } 
                <button onClick={this.props.getA}>Get A</button>
                <button onClick={this.props.getB}>Get B</button>
            </div>
        );
    }
};

class AsyncGet extends React.Component {
    constructor(props){
        super(props);
        this.getA = this.getA.bind(this);
        this.getB = this.getB.bind(this);
    }
    getA(){
        Actions.Get(this.props.dispatch, Actions.getA({}));
    }
    getB(){
        Actions.Get(this.props.dispatch, Actions.getB({}));
    }
    componentWillReceiveProps(nextProps) {
        this.props = nextProps;
    }
    render() {
        return(
            <div> 
                <GetWidget
                    isFetching={this.props.isFetching}
                    data={this.props.data}
                    getA={this.getA}
                    getB={this.getB}>
                </GetWidget>
            </div> 
        )
    }
};
var dataMapper = function(state){
    return {
        data: state.asyncget.data,
        isFetching: state.asyncget.isFetching
    };
}
exports.AsyncGet = connect(dataMapper)(AsyncGet);
