"use strict"
var Redux = require("redux");
var React = require("react");
var render = require("react-dom").render;
var thunkMiddleware = require("redux-thunk");
var loggerMiddleware = require("redux-logger")();
var Provider = require("react-redux").Provider;
var router = require('react-router');

var Root = require("./reducers/rootReducer").Root;
var RequireAuth = require("./common").RequireAuth;
var GetInitialState = require("./common").GetInitialState;

//smart components
var MathController = require("./components/mathController").MathController;
var AsyncController = require("./components/asyncController").AsyncController;
var AuthController = require("./components/authController").AuthController;
var IndexController = require("./components/indexController").IndexController;
var Logout = require("./components/authController").Logout;

const createStoreWithMiddleware = Redux.applyMiddleware(
    thunkMiddleware,
    loggerMiddleware
)(Redux.createStore);

var initialState = GetInitialState();
var store = createStoreWithMiddleware(Root, initialState);

class App extends React.Component {
    render() {
        return (
            <Provider store={store}>
                <div>
                    <h4>Flux Demo!</h4>
                    <router.Link to="/">Index</router.Link> -- 
                    <router.Link to="/math">Math</router.Link> -- 
                    <router.Link to="/async">Async</router.Link> -- 
                    <router.Link to="/auth">Auth</router.Link>  
                    <Logout history={this.props.history}/>
                     
                    <br/><br/>
                    {this.props.children}
                </div>
            </Provider>
        )
    }
}

render((
    <router.Router>
        <router.Route path="/" component={App}>
            <router.IndexRoute component={IndexController} />
            <router.Route path="math" component={MathController} onEnter={RequireAuth(store)} />
            <router.Route path="async" component={AsyncController} />
            <router.Route path="auth" component={AuthController} />
        </router.Route> 
    </router.Router>
    ),document.getElementById('root')
);

