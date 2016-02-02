module.exports = function() {
  this.World = World;
};

module.exports = function () {
  this.Given(/^hello world$/, function (callback) {
    //this.visit('https://github.com/cucumber/cucumber-js', callback);
    console.log(">>> hello world");
    callback();
  });
  this.Given(/^webdriverio$/, function (callback) {
    //this.visit('https://github.com/cucumber/cucumber-js', callback);
    console.log(">>> webdriverio");



    var webdriverio = require('webdriverio');
    var options = { desiredCapabilities: { browserName: 'phantomjs' } };
    var client = webdriverio.remote(options);
     
    client
        .init()
        .url('http://jacobra.com:8004')
        .getTitle().then(function(title) {
            console.log('title: ' + title);
            // outputs: "Title is: WebdriverIO (Software) at DuckDuckGo"
            callback();
        })
        .end();
  });


  this.Then(/^webdrivercss$/, function (done) {
    console.log(">>> webdrivercss");
    var webdriverio = require('webdriverio');
    var options = { desiredCapabilities: { browserName: 'phantomjs' } };
    var client = webdriverio.remote(options);
    require('webdrivercss').init(client, {
        screenshotRoot: 'shots',
        failedComparisonsRoot: 'diffs',
        misMatchTolerance: 0.0,
        //screenWidth: [320,480,640,1024]
        screenWidth: [640]
    });
     
    client
        .init()
        .url('http://jacobra.com:8004/')
        .webdrivercss('test', [{
            name: 'element',
            elem: '#root',
            //screenWidth: [320,640,960]
            screenWidth: [640]
        }], function(err,res) {
           //assert.ok(res.element[0].isWithinMisMatchTolerance);
           done();
        })
        .end()
/*
        //.end()
        //.call(callback);
*/
    });

/*
  this.When(/^I go to the README file$/, function (callback) {
    callback.pending();
  });

  this.Then(/^I should see "(.*)" as the page title$/, function (title, callback) {
    // matching groups are passed as parameters to the step definition

    var pageTitle = this.browser.text('title');
    if (title === pageTitle) {
      callback();
    } else {
      callback(new Error("Expected to be on page with title " + title));
    }
  });
*/
};
