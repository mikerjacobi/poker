var LogoutForm = React.createClass({
    getInitialState: function() {
        var u = document.getElementById('logout_form').getAttribute("username");
        return {username: u};
    },
    handleClick: function(event) {
        $.ajax({
            url: config.url + "/a/logout",
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
React.render(
    <LogoutForm/>,
    document.getElementById('logout_form')
);
