var Redux = require("redux");
var React = require("react");
var Root = require("./components/root").Root;
var render = require("react-dom").render;
var thunkMiddleware = require("redux-thunk");
var loggerMiddleware = require("redux-logger")();
var rootReducer = require("./reducers").rootReducer;
var Provider = require("react-redux").Provider;

const createStoreWithMiddleware = Redux.applyMiddleware(
    thunkMiddleware,
    loggerMiddleware
)(Redux.createStore);

var initialState = {count:55};
var store = createStoreWithMiddleware(rootReducer, initialState);

render(
    <Provider store={store}><Root /></Provider>,
    document.getElementById('root')
);

