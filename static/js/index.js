"use strict"
var Redux = require("redux");
var React = require("react");
var Root = require("./components/root").Root;
var AsyncGet = require("./components/asyncget").AsyncGet;
var AuthController = require("./components/auth-controller").AuthController;
var Dashboard = require("./components/dashboard").Dashboard;
var render = require("react-dom").render;
var thunkMiddleware = require("redux-thunk");
var loggerMiddleware = require("redux-logger")();
var rootReducer = require("./reducers").rootReducer;
var Provider = require("react-redux").Provider;
var Router = require('react-router').Router;
var Route = require('react-router').Route;
var Link = require('react-router').Link;
var IndexRoute = require('react-router').IndexRoute;

const createStoreWithMiddleware = Redux.applyMiddleware(
    thunkMiddleware,
    loggerMiddleware
)(Redux.createStore);

var initialState = {count:7};
var store = createStoreWithMiddleware(rootReducer, initialState);

class App extends React.Component {
    render() {
        return (
            <Provider store={store}><div>
                <h4>Flux Demo!</h4>
                <Link to="/">Home</Link> -- 
                <Link to="/math">Math</Link> -- 
                <Link to="/asyncget">AsyncGet</Link> -- 
                <Link to="/auth">Auth</Link>

                <br/><br/>

                {this.props.children}
            </div></Provider>
        )
    }
}

var enterMath = function(nextState, replaceState){
    console.log(nextState);
    replaceState({ nextPathname: nextState.location.pathname }, '/auth')
    
}

render((
    <Router>
        <Route path="/" component={App}>
        <IndexRoute component={Dashboard} />
            <Route path="math" component={Root} onEnter={enterMath} />
            <Route path="asyncget" component={AsyncGet} />
            <Route path="auth" component={AuthController} />
        </Route> 
    </Router>
    ),document.getElementById('root')
);

