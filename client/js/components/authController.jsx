"use strict"
var React = require('react')
var connect =  require('react-redux').connect;
var Auth = require("../actions/authAction");
var Account = require("../actions/accountAction");
var Nav = require("../actions/navAction");

class LoginCreateForm extends React.Component{
    constructor(props){
        super(props);
        this.changeUsername = this.changeUsername.bind(this);
        this.changePassword = this.changePassword.bind(this);
        this.changeRepeat = this.changeRepeat.bind(this);
    }
    changeUsername(event){
        this.setState({username:event.target.value});
    }
    changePassword(event){
        this.setState({password:event.target.value});
    }
    changeRepeat(event){
        this.setState({repeat:event.target.value});
    }
    render(){
        var data = <div> loading... </div>;
        if (!this.props.isFetching){
            var username = "";
            var password = "";
            var repeat = "";
            if (this.state != null){
                username = this.state.username; 
                password = this.state.password; 
                repeat = this.state.repeat; 
            }

            data = <div>
                        <div className="ui input"><input 
                            type="text"
                            id="username_textfield"
                            placeholder="username"
                            value={username} 
                            onChange={this.changeUsername}/>
                        </div>
                        <br/>
                        <div className="ui input"><input 
                            type="password"
                            id="password_textfield"
                            placeholder="password"
                            value={password} 
                            onChange={this.changePassword}/>
                        </div>
                        <button 
                            className="ui primary button"
                            id="login_button"
                            onClick={this.props.login.bind(this, username, password)}> 
                            Login 
                        </button>    

                        <br/>
                        <div className="ui input"><input 
                            type="password"
                            id="repeat_textfield"
                            placeholder="password repeat"
                            value={repeat} 
                            onChange={this.changeRepeat}/>
                        </div>
                        <button     
                            className="ui button"
                            id="create_account_button"
                            onClick={this.props.createAccount.bind(this, username, password, repeat)}> 
                            Create Account 
                        </button>
                </div>;
        }
        return data;
    }
}

class AuthController extends React.Component {
    constructor(props){
        super(props);
        this.login = this.login.bind(this);
        this.createAccount = this.createAccount.bind(this);
    }
    componentDidMount() {
        Auth.wsConnect(this.props.dispatch, this.props.wsConnection);
    }
    login(username, password) {
        Auth.Login(
            this.props.dispatch, 
            username, 
            password,
            this.props.wsConnection
        );
    }
    createAccount(username, password, repeat) {
        Account.Create(
            this.props.dispatch, 
            username, 
            password,
            repeat
        );
    }
    componentWillReceiveProps(nextProps) {
        this.props = nextProps;
    }
    render() {
        return(
            <div>
                <LoginCreateForm
                    isFetching={this.props.isFetching}
                    login={this.login}
                    createAccount={this.createAccount}>
                </LoginCreateForm>
            </div>
        )
    }
}

var loginMapper = function(state){
    return {
        isFetching: state.Account.isFetching,
        wsConnection: state.Auth.wsConnection
    };
}
exports.AuthController = connect(loginMapper)(AuthController);


class Logout extends React.Component {
    constructor(props){
        super(props);
        this.clickLogout = this.clickLogout.bind(this);
    }
    clickLogout() {
        Auth.Logout(this.props.dispatch, this.props.wsConnection);
        Nav.GoNextPath(this.props.dispatch);
    }
    componentWillReceiveProps(nextProps) {
        this.props = nextProps;
    }
    render() {
        if (!this.props.loggedIn){
            return false;
        }
        return(
            <label onClick={this.clickLogout}>
                <a href="#/">Logout</a>
            </label>
        );
    }
};
var logoutMapper = function(state){
    return {
        loggedIn:state.Auth.loggedIn,
        wsConnection: state.Auth.wsConnection
    }; 
};
exports.Logout = connect(logoutMapper)(Logout);
