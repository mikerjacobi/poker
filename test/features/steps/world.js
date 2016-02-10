var WebDriverIO = require('webdriverio');
var WebDriverCSS = require('webdrivercss');
var exec = require('child_process').execSync;
var zombie = require('zombie');
var fixtures = require("../../../server/fixtures/data");

function getDockerIP(container){
  return String(exec('docker inspect '+container+' | grep IPA | tail -n1 | awk \'{print $2\'} | cut -d\'"\' -f2')).replace("\n","");
}

function World() {
  this.browser = new zombie(); // this.browser will be available in step definitions

  this.visit = function (url, callback) {
    this.browser.visit(url, callback);
  };

  this.fixtures = fixtures
  this.appHost = "http://"+getDockerIP("server_echo_1");
  var seleniumHub = getDockerIP("server_selenium_hub_1");
  var pdiffHost = "http://" + getDockerIP("server_pdiff_server_1") + ":9000";

  var options = { 
      host:seleniumHub,
      //waitforTimeout:10000,
      desiredCapabilities: { browserName: 'chrome'} 
  };
  this.client = WebDriverIO.remote(options);
  WebDriverCSS.init(this.client, {
      screenshotRoot: 'poker',
      api: pdiffHost + '/api/repositories/'
  });
}

module.exports = function() {
  this.World = World;
};


