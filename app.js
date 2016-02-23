// Patched app.js to allow Go shell to call the Node.js Particle CLI

var Interpreter = require('./lib/interpreter.js');

function command(argv) {
  var cli = new Interpreter();
  cli.supressWarmupMessages = true;
  cli.startup();
  cli.handle(argv, true);
}

module.exports = command;
