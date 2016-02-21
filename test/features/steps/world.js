var WebDriverIO = require('webdriverio');
var WebDriverCSS = require('webdrivercss');
var exec = require('child_process').execSync;
var zombie = require('zombie');
var fixtures = require("../../../server/fixtures/data");
var mongo = require("mongodb").MongoClient;
var webdrivercssDir = "poker"

function getDockerIP(container){
  return String(exec('docker inspect '+container+' | grep IPA | tail -n1 | awk \'{print $2\'} | cut -d\'"\' -f2')).replace("\n","");
}

function World() {

    this.fixtures = fixtures
    this.appHost = "http://dev:8004";
    var seleniumHub = getDockerIP("server_hub_1");
    var pdiffHost = "http://" + getDockerIP("server_pdiff_1") + ":9000";
    var driverOptions = { 
        //logLevel:"verbose",
        host:seleniumHub,
        desiredCapabilities: { browserName: 'chrome'}
    };
    var cssOptions = {
        screenshotRoot: webdrivercssDir,
        misMatchTolerance:0.001,
        api: pdiffHost + '/api/repositories/'
    };

    //this.browser = WebDriverIO.multiremote({
    //    chrome1: driverOptions,
    //    chrome2: driverOptions
    //});
    this.clients = {
        cli1: WebDriverIO.remote(driverOptions),
        cli2: WebDriverIO.remote(driverOptions)
    }
    WebDriverCSS.init(this.clients.cli1, cssOptions);
    WebDriverCSS.init(this.clients.cli2, cssOptions);
    //WebDriverCSS.init(this.browser, cssOptions);
       
    //database init 
    var table = "echo" 
    var dbconn = getDockerIP("server_mongo_1") + "/" + table;
    this.db = require('monk')(dbconn);
}

module.exports = function() {
    this.World = World;
};

