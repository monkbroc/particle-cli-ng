package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/monkbroc/particle-cli-ng/gode"
)

var corePluginName = "particle-cli"

func SetupCore() {
	SetupPlugins(corePluginName)
}

func getCorePlugin() *Plugin {
	return GetPlugins()[corePluginName]
}

func runCoreCommand(args []string) {
	plugin := getCorePlugin()
	readLockPlugin(plugin.Name)
	argsJSON, err := json.Marshal(args)
	if err != nil {
		panic(err)
	}
	title, _ := json.Marshal(processTitle())
	cwd, _ := os.Getwd()
	script := fmt.Sprintf(`
	'use strict';
	var moduleName = '%s';
	var moduleVersion = '%s';
	process.title = %s;
	process.argv = %s;
	process.argv.unshift('node');
	var logPath = %s;
	var cwd = %s;
	process.chdir(cwd);
	process.on('uncaughtException', function (err) {
		// ignore EPIPE errors (usually from piping to head)
		if (err.code === "EPIPE") return;
		console.error(' !   Error in ' + moduleName + ':')
		console.error(' !   ' + err.message || err);
		if (err.stack) {
			var fs = require('fs');
			var log = function (line) {
				var d = new Date().toISOString()
				.replace(/T/, ' ')
				.replace(/-/g, '/')
				.replace(/\..+/, '');
				fs.appendFileSync(logPath, d + ' ' + line + '\n');
			}
			log(err.stack);
			console.error(' !   See ' + logPath + ' for more info.');
		}
		process.exit(1);
	});
	require(moduleName);`, plugin.Name, plugin.Version, string(title), argsJSON, strconv.Quote(ErrLogPath), strconv.Quote(cwd))

	// swallow sigint since the plugin will handle it
	swallowSignal(os.Interrupt)

	cmd := gode.RunScript(script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if debugging {
		cmd = gode.DebugScript(script)
	}
	if err := cmd.Run(); err != nil {
		os.Exit(getExitCode(err))
	}
}

func swallowSignal(s os.Signal) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, s)
	go func() {
		<-c
	}()
}

func getExitCode(err error) int {
	switch e := err.(type) {
	case nil:
		return 0
	case *exec.ExitError:
		status, ok := e.Sys().(syscall.WaitStatus)
		if !ok {
			panic(err)
		}
		return status.ExitStatus()
	default:
		panic(err)
	}
}
