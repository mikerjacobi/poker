//import { createStore } from 'redux';
var Redux = require("redux");

var counter = function(state, action) {
    switch (action.type) {
    case 'INCREMENT':
        return state + 1;
    case 'DECREMENT':
        return state - 1;
    default:
        return state;
    }
};

var store = Redux.createStore(counter, 0);
store.subscribe(function(){
    console.log(store.getState())
});

store.dispatch({ type: 'INCREMENT' });
store.dispatch({ type: 'INCREMENT' });
store.dispatch({ type: 'DECREMENT' });
store.dispatch({ type: 'INCREMENT' });
