var WebDriverIO = require('webdriverio');
var WebDriverCSS = require('webdrivercss');

module.exports = function() {
  this.World = World;
};

module.exports = function () {
  var seleniumHub = "172.17.0.3";
  var appHost = "http://172.17.0.6"
  var pdiffHost = "http://dev:9000";
  var options = { 
      host:seleniumHub,
      //waitforTimeout:10000,
      desiredCapabilities: { browserName: 'chrome'} 
  };
  var client = WebDriverIO.remote(options);
  WebDriverCSS.init(client, {
      screenshotRoot: 'poker',
      api: pdiffHost + '/api/repositories/'
  });

  this.Then(/^webdrivercss is demoed$/, function (done) {
    var url = appHost + "/";
    client
        .init()
        .sync()
        .url(url)
        .getTitle().then(function(title) {
            console.log('title: ' + title);
        })
        .webdrivercss('indexPage', [{
            name: 'element',
            elem: '#root',
            screenWidth: [640]
        }], function(err,res) {
            if (err != undefined){
              console.log("err: "+err);
            } else {
              console.log("success");
            }
           //assert.ok(res.element[0].isWithinMisMatchTolerance);
        })
        .sync()
        .end()
        .call(done);
    });
};
