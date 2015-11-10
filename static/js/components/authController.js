"use strict"
var React = require('react')
var connect =  require('react-redux').connect;
var Auth = require("../actions/authAction");

class LoginCreateForm extends React.Component{
    render(){
        var data = <div> loading... </div>;
        if (!this.props.isFetching){
            data = <div>
                    <input type="text"
                        placeholder="username"
                        value={this.props.username} 
                        onChange={this.props.changeUsername}/>
                    
                    <br/><br/>
                    <input type="password"
                        placeholder="password"
                        value={this.props.password} 
                        onChange={this.props.changePassword}/>
                    <button onClick={this.props.clickLogin}> Login </button>    

                    <br/><br/>
                    <input type="password"
                            placeholder="password repeat"
                            value={this.props.repeat} 
                            onChange={this.props.changeRepeat}/>
                    <button onClick={this.props.clickCreateAccount}> Create Account </button>
                </div>;
        }
        return data;
    }
}

class AuthController extends React.Component {
    constructor(props){
        super(props);
        this.changeUsername = this.changeUsername.bind(this);
        this.changePassword = this.changePassword.bind(this);
        this.changeRepeat = this.changeRepeat.bind(this);
        this.clickLogin = this.clickLogin.bind(this);
        this.clickCreateAccount = this.clickCreateAccount.bind(this);
    }
    changeUsername(event) {
        var username = event.target.value;
        Auth.ChangeUsername(this.props.dispatch, username);
    }
    changePassword(event) {
        var password = event.target.value;
        Auth.ChangePassword(this.props.dispatch, password);
    }
    changeRepeat(event) {
        var repeat = event.target.value;
        Auth.ChangeRepeat(this.props.dispatch, repeat);
    }
    clickLogin() {
        Auth.Login(
            this.props.dispatch, 
            this.props.username, 
            this.props.password,
            this.props.history
        );
    }
    clickCreateAccount() {
        Auth.CreateAccount(
            this.props.dispatch, 
            this.props.username, 
            this.props.password,
            this.props.repeat
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
                    username={this.props.username}
                    password={this.props.password}
                    repeat={this.props.repeat}
                    changeUsername={this.changeUsername}
                    changePassword={this.changePassword}
                    changeRepeat={this.changeRepeat}
                    clickLogin={this.clickLogin}
                    clickCreateAccount={this.clickCreateAccount}>
                </LoginCreateForm>
            </div>
        )
    }
}

var dataMapper = function(state){
    return {
        isFetching: state.Auth.isFetching,
        username: state.Auth.username,
        password: state.Auth.password,
        repeat: state.Auth.repeat
    };
}
exports.AuthController = connect(dataMapper)(AuthController);


class Logout extends React.Component {
    constructor(props){
        super(props);
        this.clickLogout = this.clickLogout.bind(this);
    }
    clickLogout() {
        Auth.Logout(this.props.dispatch, this.props.history);
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
                -- <a href="#/">Logout</a>
            </label>
        );
    }
};
var logoutMapper = function(state){
    return {
        loggedIn:state.Logout.loggedIn
    }; 
};
exports.Logout = connect(logoutMapper)(Logout);