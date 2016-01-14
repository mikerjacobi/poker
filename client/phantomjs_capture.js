var page = require('webpage').create();
var timeout = 1;
var url = 'http://jacobra.com:8004/dpxdt/demo.html';
var url = 'http://jacobra.com:8004';

page.viewportSize = {
    width: 800,
    height: 600
};
page.open(url, function() {
    setTimeout(function () {
        page.render("jacobra.png");
        phantom.exit();
    }, timeout);  
});
page.onError = function(msg, trace) {
    console.error(msg);
};
