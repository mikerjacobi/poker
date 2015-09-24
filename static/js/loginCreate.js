var LoginCreateForm = React.createClass({
    getInitialState: function() {
        return {username: '', password: '', repeat: ''};
    },
    changeUsername: function(event) {
        this.setState({username: event.target.value});
    },
    changePassword: function(event) {
        this.setState({password: event.target.value});
    },
    changeRepeat: function(event) {
        this.setState({repeat: event.target.value});
    },
    clickLogin: function(event) {
        $.ajax({
            url: this.props.baseurl + "/login",
            method: "POST",
            data:{
                "username":this.state.username,
                "password":this.state.password,
            },
            dataType:"json",
            success: function(data) {
                reactCookie.save('session', data.payload.session_id);
                window.location.replace("/a/");
            }.bind(this),
                error: function(xhr, status, err) {
                console.error(this.props, status, err.toString());
            }.bind(this)
        });
    },
    clickCreateAccount: function(event) {
        $.ajax({
            url: this.props.baseurl + "/create_account",
            method: "POST",
            data:{
                "username":this.state.username,
                "pw1":this.state.password,
                "pw2":this.state.repeat
            },
            dataType:"json",
            success: function(data) {
                console.log("success");
                console.log(data);
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(this.props, status, err.toString());
            }.bind(this)
        });
    },
    render: function() {
        var input_username = this.state.input_username;
        var input_password = this.state.input_password;
        var input_repeat = this.state.input_repeat;
        return(
            <fieldset>
                <div className="col-md-3">
                    <input type="text"
                           placeholder="username"
                           value={input_username}
                           onChange={this.changeUsername}
                           id="username_input"
                           ref="username_input"
                           className="form-control"
                    />
                </div> <br/><br/>
                <div className="col-md-3">
                    <input type="password"
                           placeholder="password"
                           value={input_password}
                           onChange={this.changePassword}
                           id="password_input"
                           ref="password_input"
                           className="form-control"
                    />
                </div>
                <div className="col-md-3">
                    <button className="btn btn-default" onClick={this.clickLogin}>
                        Login
                    </button>
                </div><br/><br/>
                <div className="col-md-3">
                    <input type="password"
                           placeholder="password repeat"
                           value={input_repeat}
                           onChange={this.changeRepeat}
                           id="repeat_input"
                           ref="repeat_input"
                           className="form-control"
                    />
                </div>
                <div className="col-md-3">
                    <button className="btn btn-default" onClick={this.clickCreateAccount}>
                        Create Account 
                    </button>
                </div> 
            </fieldset>
        );
    }
});
React.render(
  <LoginCreateForm baseurl={config.url}/>,
  document.getElementById('login_create_form')
);
