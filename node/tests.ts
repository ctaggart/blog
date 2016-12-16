import Mocha = require("mocha");
import fetch = require("node-fetch");
(<any>global).fetch = fetch;
var runner = new Mocha({ui: "bdd"});
runner.addFile("node/build/fetch-tests.js");
runner.run();