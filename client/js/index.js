"use strict"
var Redux = require("redux");
var React = require("react");
var render = require("react-dom").render;
var thunkMiddleware = require("redux-thunk");
var loggerMiddleware = require("redux-logger")();
var Provider = require("react-redux").Provider;
var router = require('react-router');

var Root = require("./reducers/rootReducer").Root;
var Common = require("./common");
var GetInitialState = require("./common").GetInitialState;

//smart components
var MathController = require("./components/mathController").MathController;
var LobbyController = require("./components/lobbyController").LobbyController;
var HoldemController = require("./components/holdemController").HoldemController;
var AsyncController = require("./components/asyncController").AsyncController;
var AuthController = require("./components/authController").AuthController;
var IndexController = require("./components/indexController").IndexController;
var Logout = require("./components/authController").Logout;

//actions
var Nav = require("./actions/navAction");

const createStoreWithMiddleware = Redux.applyMiddleware(
    thunkMiddleware,
    loggerMiddleware
)(Redux.createStore);

var initialState = GetInitialState();
var store = createStoreWithMiddleware(Root, initialState);

class App extends React.Component {
    componentDidMount(){
        Nav.SetHistory(store.dispatch, this.props.history);
    }
    render() {
        return (
            <Provider store={store}>
                <div>
                    <div className="ui fixed inverted menu">
                        <router.Link className="item" onClick={Common.SetPath(store)} to="/">Index</router.Link>  
                        <router.Link className="item" onClick={Common.SetPath(store)} to="/math">Math</router.Link>  
                        <router.Link className="item" onClick={Common.SetPath(store)} to="/lobby">Lobby</router.Link>  
                        <router.Link className="item" onClick={Common.SetPath(store)} to="/auth">Auth</router.Link>  
                        <div className="item"><Logout /></div>
                    </div>
                         
                    <br/>

                    <div>{this.props.children}</div>
                </div>
            </Provider>
        )
    }
}

render((
    <router.Router>
        <router.Route path="/" component={App}>
            <router.IndexRoute component={IndexController} />
            <router.Route path="math" component={MathController} onEnter={Common.RequireAuth(store)} />
            <router.Route path="lobby">
                <router.IndexRoute component={LobbyController}  onEnter={Common.RequireAuth(store)}/>
                <router.Route path="/holdem/:gameid" component={HoldemController} onEnter={Common.RequireAuth(store)} />
            </router.Route>
            <router.Route path="async" component={AsyncController} />
            <router.Route path="auth" component={AuthController} />
        </router.Route> 
    </router.Router>
    ),document.getElementById('root')
);
