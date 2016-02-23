package main

import (
	"github.com/monkbroc/particle-cli-ng/gode"
)

// SetupNode sets up node and npm in ~/.particle
func SetupNode() {
	gode.SetRootPath(AppDir())
	setup, err := gode.IsSetup()
	PrintError(err, false)
	if !setup {
		setupNode()
	}
}

func setupNode() {
	Err("particle-cli: Adding dependencies...")
	PrintError(gode.Setup(), true)
	Errln(" done")
}

func updateNode() {
	gode.SetRootPath(AppDir())
	needsUpdate, err := gode.NeedsUpdate()
	PrintError(err, true)
	if needsUpdate {
		setupNode()
	}
}
