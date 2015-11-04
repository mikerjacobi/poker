"use strict"
var Redux = require("redux");
var React = require("react");
var Root = require("./components/root").Root;
var About = require("./components/about").About;
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
                <h3>Flux Demo!</h3>
                <Link to="/">Home</Link> -- 
                <Link to="/math">Math</Link> -- 
                <Link to="/about">About</Link>
                {this.props.children}
            </div></Provider>
        )
    }
}

render((
    <Router>
        <Route path="/" component={App}>
        <IndexRoute component={Dashboard} />
            <Route path="math" component={Root} />
            <Route path="about" component={About} />
        </Route> 
    </Router>
    ),document.getElementById('root')
);

