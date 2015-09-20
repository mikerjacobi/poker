var LogoutForm = React.createClass({
  getInitialState: function() {
    return {};
  },
  handleClick: function(event) {
    $.ajax({
        url: this.props.posturl,
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
                <div className="col-md-3">
                    <button className="btn btn-default" onClick={this.handleClick}>
                        Logout
                    </button>
                </div>
            </fieldset>
        );
    }
});

React.render(
    <LogoutForm posturl={config.url + "/a/logout"}/>,
    document.getElementById('logout_form')
);
