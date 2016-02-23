package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/monkbroc/particle-cli-ng/Godeps/_workspace/src/github.com/dickeyxxx/golock"
	"github.com/monkbroc/particle-cli-ng/gode"
)

// Plugin represents a javascript plugin
type Plugin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ParsePlugin requires the plugin's node module
// to get the commands and metadata
func ParsePlugin(name string) (*Plugin, error) {
	script := `
	var plugin = {};
	var pjson  = require('` + name + `/package.json');

	plugin.name    = pjson.name;
	plugin.version = pjson.version;

	console.log(JSON.stringify(plugin))`
	cmd := gode.RunScript(script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Error reading plugin: %s\n%s\n%s", name, err, string(output))
	}
	var plugin Plugin
	json.Unmarshal([]byte(output), &plugin)
	return &plugin, nil
}

// GetPlugins goes through all the node plugins and returns them in Go stucts
func GetPlugins() map[string]*Plugin {
	plugins := FetchPluginCache()
	for name, plugin := range plugins {
		if plugin == nil || !pluginExists(name) {
			delete(plugins, name)
		}
	}
	return plugins
}

// PluginNames lists all the plugin names
func PluginNames() []string {
	plugins := FetchPluginCache()
	names := make([]string, 0, len(plugins))
	for _, plugin := range plugins {
		if plugin != nil {
			names = append(names, plugin.Name)
		}
	}
	return names
}

// PluginNamesNotSymlinked returns all the plugins that are not symlinked
func PluginNamesNotSymlinked() []string {
	a := PluginNames()
	b := make([]string, 0, len(a))
	for _, plugin := range a {
		if !isPluginSymlinked(plugin) {
			b = append(b, plugin)
		}
	}
	return b
}

func isPluginSymlinked(plugin string) bool {
	path := filepath.Join(AppDir(), "node_modules", plugin)
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeSymlink != 0
}

// ensure all the Javascript plugin are installed
func SetupPlugins(pluginNames ...string) {
	newPluginNames := difference(pluginNames, PluginNames())
	if len(newPluginNames) == 0 {
		return
	}
	Err("particle-cli-ng: Installing plugins...")
	if err := installPlugins(newPluginNames...); err != nil {
		// retry once
		PrintError(gode.RemovePackages(newPluginNames...), true)
		PrintError(gode.ClearCache(), true)
		Err("\rparticle-cli-ng: Installing plugins (retrying)...")
		ExitIfError(installPlugins(newPluginNames...), true)
	}
	Errln(" done")
}

func difference(a, b []string) []string {
	res := make([]string, 0, len(a))
	for _, aa := range a {
		if !contains(b, aa) {
			res = append(res, aa)
		}
	}
	return res
}

func contains(arr []string, s string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}

func installPlugins(names ...string) error {
	for _, name := range names {
		lockPlugin(name)
	}
	defer func() {
		for _, name := range names {
			unlockPlugin(name)
		}
	}()
	err := gode.InstallPackages(names...)
	if err != nil {
		return err
	}
	plugins := make([]*Plugin, 0, len(names))
	for _, name := range names {
		plugin, err := ParsePlugin(name)
		if err != nil {
			return err
		}
		plugins = append(plugins, plugin)
	}
	AddPluginsToCache(plugins...)
	return nil
}

func pluginExists(plugin string) bool {
	exists, _ := fileExists(pluginPath(plugin))
	return exists
}

// directory location of plugin
func pluginPath(plugin string) string {
	return filepath.Join(AppDir(), "node_modules", plugin)
}

// lock a plugin for reading
func readLockPlugin(name string) {
	lockfile := updateLockPath + "." + name
	if exists, _ := fileExists(lockfile); exists {
		lockPlugin(name)
		unlockPlugin(name)
	}
}

// lock a plugin for writing
func lockPlugin(name string) {
	LogIfError(golock.Lock(updateLockPath + "." + name))
}

// unlock a plugin
func unlockPlugin(name string) {
	LogIfError(golock.Unlock(updateLockPath + "." + name))
}
