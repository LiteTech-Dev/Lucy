/*
Copyright 2024 4rcadia

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package local provides functionality to gather and manage server information
// for a Minecraft server. It includes methods to retrieve server configuration,
// mod list, executable information, and other relevant details. The package
// utilizes memoization to avoid redundant calculations and resolve any data
// dependencies issues. Therefore, all probe functions are 100% concurrent-safe.
//
// The main exposed function is GetServerInfo, which returns a comprehensive
// ServerInfo struct containing all the gathered information. To avoid side
// effects, the ServerInfo struct is returned as a copy, rather than reference.
package local

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"sort"
	"sync"

	"gopkg.in/ini.v1"

	"lucy/datatypes"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/tools"
)

const mcdrConfigFileName = "config.yml"

// GetServerInfo is the exposed function for external packages to get serverInfo.
// As we can assume that the environment does not change while the program is
// running, a sync.Once is used to prevent further calls to this function. Rather,
// the cached serverInfo is used as the return value.
var GetServerInfo = tools.Memoize(buildServerInfo)

// buildServerInfo builds the server information by performing several checks
// and gathering data from various sources. It uses goroutines to perform these
// tasks concurrently and a sync.Mutex to ensure thread-safe updates to the
// serverInfo struct.
func buildServerInfo() lucytypes.ServerInfo {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var serverInfo lucytypes.ServerInfo

	// MCDR Stage
	wg.Add(1)
	go func() {
		defer wg.Done()
		mcdrConfig := getMcdrConfig()
		if mcdrConfig != nil {
			mu.Lock()
			serverInfo.Mcdr = &lucytypes.McdrInstallation{
				PluginPaths: mcdrConfig.PluginDirectories,
			}
			mu.Unlock()
		}
	}()

	// Server Work Path
	wg.Add(1)
	go func() {
		defer wg.Done()
		workPath := getServerWorkPath()
		mu.Lock()
		serverInfo.WorkPath = workPath
		mu.Unlock()
	}()

	// Executable Stage
	wg.Add(1)
	go func() {
		defer wg.Done()
		executable := getExecutableInfo()
		mu.Lock()
		serverInfo.Executable = executable
		mu.Unlock()
	}()

	// Mod Path
	wg.Add(1)
	go func() {
		defer wg.Done()
		modPath := getServerModPath()
		mu.Lock()
		serverInfo.ModPath = modPath
		mu.Unlock()
	}()

	// Mod List
	wg.Add(1)
	go func() {
		defer wg.Done()
		modList := getMods()
		mu.Lock()
		serverInfo.Mods = modList
		mu.Unlock()
	}()

	// Save Path
	wg.Add(1)
	go func() {
		defer wg.Done()
		savePath := getSavePath()
		mu.Lock()
		serverInfo.SavePath = savePath
		mu.Unlock()
	}()

	// Check for Lucy installation
	wg.Add(1)
	go func() {
		defer wg.Done()
		hasLucy := checkHasLucy()
		mu.Lock()
		serverInfo.HasLucy = hasLucy
		mu.Unlock()
	}()

	// Check if the server is running
	wg.Add(1)
	go func() {
		defer wg.Done()
		activity := checkServerFileLock()
		mu.Lock()
		serverInfo.Activity = activity
		mu.Unlock()
	}()

	// Server Mod Path
	wg.Add(1)
	go func() {
		defer wg.Done()
		modPath := getServerModPath()
		mu.Lock()
		serverInfo.ModPath = modPath
		mu.Unlock()
	}()

	wg.Wait()
	return serverInfo
}

// Some functions that gets a single piece of information. They are not exported,
// as GetServerInfo() applies a memoization mechanism. Every time a serverInfo
// is needed, just call GetServerInfo() without the concern of redundant calculation.

var getServerModPath = tools.Memoize(
	func() string {
		exec := getExecutableInfo()
		if exec.Platform == lucytypes.Fabric || exec.Platform == lucytypes.Forge {
			return path.Join(getServerWorkPath(), "mods")
		}
		return ""
	},
)

var getServerWorkPath = tools.Memoize(
	func() string {
		if mcdrConfig := getMcdrConfig(); mcdrConfig != nil {
			return mcdrConfig.WorkingDirectory
		}
		return "."
	},
)

var getServerDotProperties = tools.Memoize(
	func() MinecraftServerDotProperties {
		exec := getExecutableInfo()
		propertiesPath := path.Join(getServerWorkPath(), "server.properties")
		file, err := ini.Load(propertiesPath)
		if err != nil {
			if exec != UnknownExecutable {
				logger.Warning(errors.New("this server is missing a server.properties"))
			}
			return nil
		}

		properties := make(map[string]string)
		for _, section := range file.Sections() {
			for _, key := range section.Keys() {
				properties[key.Name()] = key.String()
			}
		}

		return properties
	},
)

var getSavePath = tools.Memoize(
	func() string {
		serverProperties := getServerDotProperties()
		if serverProperties == nil {
			return ""
		}
		levelName := serverProperties["level-name"]
		return path.Join(getServerWorkPath(), levelName)
	},
)

var checkHasLucy = tools.Memoize(
	func() bool {
		_, err := os.Stat(".lucy")
		return err == nil
	},
)

var getMods = tools.Memoize(
	func() (mods []lucytypes.Package) {
		path := getServerModPath()
		jars := findJar(path)
		for _, jar := range jars {
			mod := analyzeModJar(jar)
			if mod != nil {
				mods = append(mods, *mod)
			}
		}
		sort.Slice(
			mods,
			func(i, j int) bool { return mods[i].Id.Name < mods[j].Id.Name },
		)
		return mods
	},
)

const fabricModIdentifierFile = "fabric.mod.json"

// const forgeModIdentifierFile =
// TODO: forgeModIdentifierFile

func analyzeModJar(file *os.File) *lucytypes.Package {
	stat, err := file.Stat()
	if err != nil {
		return nil
	}
	r, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return nil
	}

	for _, f := range r.File {
		if f.Name == fabricModIdentifierFile {
			rr, err := f.Open()
			data, err := io.ReadAll(rr)
			if err != nil {
				return nil
			}
			modInfo := &datatypes.FabricModIdentifier{}
			err = json.Unmarshal(data, modInfo)
			if err != nil {
				return nil
			}
			p := &lucytypes.Package{
				Id: lucytypes.PackageId{
					Platform: lucytypes.Fabric,
					Name:     lucytypes.PackageName(modInfo.Id),
					Version:  lucytypes.PackageVersion(modInfo.Version),
				},
				Local: &lucytypes.PackageInstallation{
					Path: file.Name(),
				}, Information: nil, // Don't need this for now
				Dependencies:   nil, // TODO: This is not yet implemented, because the deps field is an expression, we need to parse it
			}
			return p
		}
	}

	return nil
}
