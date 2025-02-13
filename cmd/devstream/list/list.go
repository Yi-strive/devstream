package list

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// list is the version of DevStream.
// Assign the value when building with the -X parameter. Example:
// -X github.com/devstream-io/devstream/cmd/devstream/list.PluginsName=${PLUGINS_NAME}
// See the Makefile for more info.

var PluginsName string

// List all plugins name
func List(pluginFilter string) {
	r, _ := regexp.Compile(pluginFilter)
	for _, pluginName := range PluginsNameSlice() {
		if r.Match([]byte(pluginName)) {
			fmt.Println(pluginName)
		}
	}
}

// PluginsNameSlice Gets plugins name in slice
func PluginsNameSlice() []string {
	listPluginsName := strings.Fields(PluginsName)
	sort.Strings(listPluginsName)
	return listPluginsName
}

// PluginNamesMap Gets plugins name in map
func PluginNamesMap() map[string]struct{} {
	mp := make(map[string]struct{})

	listPluginsName := strings.Fields(PluginsName)

	for _, pluginName := range listPluginsName {
		mp[pluginName] = struct{}{}
	}

	return mp
}
