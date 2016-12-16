SystemJS.config({
    meta: { 
        'node_modules/mocha/mocha.js': { format: 'global' }
    },
    map: {
        chai: "node_modules/chai/chai.js",
        mocha: "node_modules/mocha/mocha.js"
    }
});
SystemJS.import("mocha").then(m => {
    mocha.setup("bdd"); // https://mochajs.org/#bdd
    SystemJS.import("build/fetch-tests.js").then(m => {
        mocha.run();
    });
});