var React = require('react')
var LogoutForm = React.createClass({
    getInitialState: function() {
        var u = document.getElementById('logout_form').getAttribute("username");
        return {username: u};
    },
    handleClick: function(event) {
        $.ajax({
            url: this.props.baseurl + "/logout",
            method: "POST",
            dataType:"json",
            headers:{"x-session":reactCookie.load("session")},
            success: function(data) {
                reactCookie.remove('session');
                window.location.replace("/");
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(this.props, status, err.toString());
            }.bind(this)
        });
    },
    render: function() {
        return(
            <fieldset>
                {this.state.username} &nbsp;&nbsp;
                <button className="btn btn-default" onClick={this.handleClick}>Logout</button>
            </fieldset>
        );
    }
});

exports.LogoutForm = LogoutForm;
module.exports = exports;

