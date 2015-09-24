var CreateGameForm = React.createClass({
    getInitialState: function() {
        return {gamename:''};
    },
    changeGamename: function(event) {
        this.setState({gamename: event.target.value});
    },
    handleClick: function(event) {
        $.ajax({
            url: config.url + "/a/game",
            method: "POST",
            data: {'gamename':this.state.gamename },
            dataType:"json",
            headers:{"x-session":reactCookie.load("session")},
            success: function(data) {
                console.log("game created");
                console.log(data);
                this.props.createGame(data.payload);
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(this.props, status, err.toString());
            }.bind(this)
        });
    },
    render: function() {
        var gamename = this.state.gamename;
        return(
            <fieldset>
                <div className="col-md-3">
                    <input type="text"
                           placeholder="game name"
                           value={gamename}
                           onChange={this.changeGamename}
                           id="gamename_input"
                           ref="gamename_input"
                           className="form-control"
                    />
                    <button className="btn btn-default" onClick={this.handleClick}>
                        Create Game 
                    </button>
                </div>
            </fieldset>
        );
    }
});

var Game = React.createClass({
    render: function() {
        return (
            <div className="game">
                {this.props.gameName} 
                <button className="btn btn-default" onClick={this.props.clickFunc}>Join</button>
                {this.props.players} 
            </div>
        );
    }
});

var GameList = React.createClass({
    getInitialState: function() {
        return {data: []};
    },
    clickJoinGame: function(i) {
        this.props.joinGame(i);
    },
    render: function() {
        var gameNodes = this.props.games.map(function (game, i) {
            return (
                <Game clickFunc={this.clickJoinGame.bind(this, i)} key={i} gameID={game.game_id} gameName={game.name} players={game.players}></Game>
            );
        }, this);
        return (
            <div className="gameList">
                {gameNodes}
            </div>
        );
    }
});

var GameForms = React.createClass({
    getInitialState: function() {
        return {games: []};
    },
    createGame: function(newGame){
        var games = this.state.games
        games.push(newGame);
        this.setState(games);
    },
    joinGame: function(i){
        $.ajax({
            url: config.url + "/a/game/" + this.state.games[i].game_id + "/join",
            method: "POST",
            headers:{"x-session":reactCookie.load("session")},
            dataType:"json",
            success: function(data) {
                $.ajax({
                    url: config.url + "/a/game/" + this.state.games[i].game_id,
                    method: "GET",
                    headers:{"x-session":reactCookie.load("session")},
                    success: function(data) {
                        var games = this.state.games
                        games[i] = data.payload;
                        this.setState(games);
                    }.bind(this),
                    error: function(xhr, status, err) {
                        console.error(this.props, status, err.toString());
                        return [];
                    }.bind(this)
                });
            }.bind(this),
                error: function(xhr, status, err) {
                console.error(this.state, status, err.toString());
            }.bind(this)
        });
    },
    componentDidMount: function() {
        $.ajax({
            url: config.url + "/a/games",
            method: "GET",
            headers:{"x-session":reactCookie.load("session")},
            success: function(data) {
                this.setState({"games": data["payload"]});
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(this.props, status, err.toString());
                return [];
            }.bind(this)
        });
    },
    render: function(){
        return (
            <div>
                <CreateGameForm createGame={this.createGame}/>
                <GameList joinGame={this.joinGame} games={this.state.games}/>
            </div>
        ); 
        
    }

})

React.render(
    <GameForms />,
    document.getElementById('game_forms')
);
