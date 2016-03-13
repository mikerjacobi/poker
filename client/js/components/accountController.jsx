"use strict"

var React = require("react");
var connect = require('react-redux').connect;
var Account = require("../actions/accountAction");
var Auth = require("../actions/authAction");

class ChipController extends React.Component {
    constructor(props){
        super(props);
        this.changeChipRequest = this.changeChipRequest.bind(this);
        this.state = {requestAmount: 0};
    }
    changeChipRequest(event){
        this.setState({requestAmount:event.target.value});
    }
    render() {
        if (this.props.account == null){
            return(<div> loading... </div>);
        }

        return(<div id="chips_controller_div"> 
            username: {this.props.account.username} <br/>
            balance: {this.props.account.balance} <br/>

                <div className="ui input"><input 
                    type="text"
                    id="balance_textfield"
                    placeholder="balance"
                    value={this.state.requestAmount}
                    onChange={this.changeChipRequest}/> </div>
                <button 
                    className="ui button"
                    id="request_chips_button"
                    onClick={this.props.requestChips.bind(this, parseInt(this.state.requestAmount))}>
                    Request Chips
                </button>
        </div>);
    };
};


class AccountController extends React.Component {
    constructor(props){
        super(props);
        this.requestChips = this.requestChips.bind(this);
    }
    componentDidMount() {
        Auth.connect(this.props.dispatch, this.props.wsConnection, Account.Init);
    }
    componentWillReceiveProps(nextProps) {
        this.props = nextProps;
    }
    requestChips(amount){
        Account.RequestChips(this.props.dispatch, this.props.wsConnection, amount);
    }
    render() {
        return (
            <div>
                <ChipController {...this.props.account} requestChips={this.requestChips}/>
            </div>
        );
    }
};

var dataMapper = function(state){
    return {
        account: state.Account,
        wsConnection: state.Auth.wsConnection
    };
}

exports.AccountController = connect(dataMapper)(AccountController);
