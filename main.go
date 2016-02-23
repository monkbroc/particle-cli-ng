package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/monkbroc/particle-cli-ng/Godeps/_workspace/src/github.com/stvp/rollbar"
)

// Version is the version of the CLI.
// This is set by a build flag in the `Rakefile`.
// If it is set to `dev` it will not autoupdate.
var Version = "dev"

// Channel is the git branch the code was compiled on.
// This is set by a build flag in the `Rakefile` based on the git branch.
var Channel = "?"

var cli = &Cli{}

// BuiltinPlugins are the core plugins that will be autoinstalled
var BuiltinPlugins = []string{
	"particle-cli",
}

func init() {
	rollbar.Platform = "client"
	rollbar.Token = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	rollbar.Environment = Channel
	rollbar.ErrorWriter = nil
}

func main() {
	defer handlePanic()
	runtime.GOMAXPROCS(1) // more procs causes runtime: failed to create new OS thread on Ubuntu
	ShowDebugInfo()
	Update(Channel, "block")
	SetupNode()
	err := cli.Run(os.Args)
	SetupBuiltinPlugins()
	TriggerBackgroundUpdate()
	if err != nil {
		PrintError(err, false)
		os.Exit(2)
	}
}

func handlePanic() {
	if rec := recover(); rec != nil {
		err, ok := rec.(error)
		if !ok {
			err = errors.New(rec.(string))
		}
		Errln("ERROR:", err)
		if Channel == "?" {
			debug.PrintStack()
		} else {
			rollbar.Error(rollbar.ERR, err, rollbarFields()...)
			rollbar.Wait()
		}
		Exit(1)
	}
}

func rollbarFields() []*rollbar.Field {
	var cmd string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	return []*rollbar.Field{
		{"Version", Version},
		{"GOOS", runtime.GOOS},
		{"GOARCH", runtime.GOARCH},
		{"command", cmd},
	}
}

// ShowDebugInfo prints debugging information if HEROKU_DEBUG=1
func ShowDebugInfo() {
	if !isDebugging() {
		return
	}
	info := []string{version(), binPath}
	if len(os.Args) > 1 {
		info = append(info, fmt.Sprintf("cmd: %s", os.Args[1]))
	}
	proxy := getProxy()
	if proxy != nil {
		info = append(info, fmt.Sprintf("proxy: %s", proxy))
	}
	Debugln(strings.Join(info, " "))
}

func getProxy() *url.URL {
	req, err := http.NewRequest("GET", "https://api.heroku.com", nil)
	PrintError(err, false)
	proxy, err := http.ProxyFromEnvironment(req)
	PrintError(err, false)
	return proxy
}
