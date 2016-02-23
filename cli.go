package main

import (
	"os"
	"strings"
)

// Cli handles parsing and dispatching of commands
type Cli struct {
}

// Run parses command line arguments and runs the associated command or help.
// Also does lookups for app name and/or auth token if the command needs it.
func (cli *Cli) Run(args []string) (err error) {
	// update is a special command that runs in the Go shell
	if args[0] == "update" {
		UpdateCommand(args)
		return nil
	}

	// Run all other commands in the Node.js plugin
	for _, plugin := range GetPlugins() {
		runInPlugin(plugin, args)
		return nil
	}
	return nil
}

func processTitle() string {
	return "particle " + strings.Join(os.Args[1:], " ")
}
