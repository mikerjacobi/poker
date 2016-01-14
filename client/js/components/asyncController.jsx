"use strict"
var React = require("react");
var connect =  require('react-redux').connect;
var Async = require("../actions/asyncAction")

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

class AsyncController extends React.Component {
    constructor(props){
        super(props);
        this.getA = this.getA.bind(this);
        this.getB = this.getB.bind(this);
    }
    getA(){
        Async.GetA(this.props.dispatch);
    }
    getB(){
        Async.GetB(this.props.dispatch);
    }
    componentWillReceiveProps(nextProps) {
        this.props = nextProps;
    }
    render() {
        var data = <div> loading... </div>;
        if (!this.props.isFetching){
            data = <div> 
                <GetWidget
                    isFetching={this.props.isFetching}
                    data={this.props.data}
                    getA={this.getA}
                    getB={this.getB}>
                </GetWidget>
            </div> ;
        }
        return data;
    }
};
var dataMapper = function(state){
    return {
        data: state.Async.data,
        isFetching: state.Async.isFetching
    };
}
exports.AsyncController = connect(dataMapper)(AsyncController);
