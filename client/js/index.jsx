"use strict"

if (!('bind' in Function.prototype)) {
    Function.prototype.bind = function() {
        var funcObj = this;
        var extraArgs = Array.prototype.slice.call(arguments);
        var thisObj = extraArgs.shift();
        return function() {
            return funcObj.apply(thisObj, extraArgs.concat(Array.prototype.slice.call(arguments)));
        };
    };
}

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
var AccountController = require("./components/accountController").AccountController;
var LobbyController = require("./components/lobbyController").LobbyController;
var HoldemController = require("./components/holdemController").HoldemController;
var HighCardController = require("./components/highcardController").HighCardController;
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
                        <router.Link className="item" onClick={Common.SetPath(store)} to="/account">Account</router.Link>  
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
                <router.Route path="/highcard/:gameid" component={HighCardController} onEnter={Common.RequireAuth(store)} />
            </router.Route>
            <router.Route path="async" component={AsyncController} />
            <router.Route path="account" component={AccountController} onEnter={Common.RequireAuth(store)} />
            <router.Route path="auth" component={AuthController} />
        </router.Route> 
    </router.Router>
    ),document.getElementById('root')
);

