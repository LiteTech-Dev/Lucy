package probe

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"gopkg.in/ini.v1"
	"io"
	"lucy/apitypes"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/syntaxtypes"
	"lucy/tools"
	"os"
	"path"
	"sort"
	"sync"
)

const mcdrConfigFileName = "config.yml"

// GetServerInfo is the exposed function for external packages to get serverInfo`.
// As we can assume that the environment do not change while the program is
// running, a sync.Once is used to prevent further calls to this function. Rather,
// the cached serverInfo is used as the return value.
var GetServerInfo = tools.Memoize(buildServerInfo)

// buildServerInfo
// Sequence:
//  1. Check for MCDR
//  2. Locate and unzip the jar file, if multiple valid jar files exist, prompt
//     the user to select one
//  3. From the jar we can detect Minecraft, Forge and(or) Fabric versions
//  4. Then search for related dirs (mods/, config/, plugins/, etc.)
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
			serverInfo.Mcdr = &lucytypes.Mcdr{
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
		modList := getModList()
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
		if exec.Platform == syntaxtypes.Fabric || exec.Platform == syntaxtypes.Forge {
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
		propertiesPath := path.Join(getServerWorkPath(), "server.properties")
		file, err := ini.Load(propertiesPath)
		if err != nil {
			logger.CreateWarning(errors.New("this server is missing a server.properties"))
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

var getModList = tools.Memoize(
	func() (mods []*lucytypes.Package) {
		path := getServerModPath()
		jars := findJar(path)
		for _, jar := range jars {
			mod := analyzeModJar(jar)
			if mod != nil {
				mods = append(mods, mod)
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
			modInfo := &apitypes.FabricModIdentifier{}
			err = json.Unmarshal(data, modInfo)
			if err != nil {
				return nil
			}
			p := &lucytypes.Package{
				Id: syntaxtypes.PackageId{
					Platform: syntaxtypes.Fabric,
					Name:     syntaxtypes.PackageName(modInfo.Id),
					Version:  syntaxtypes.PackageVersion(modInfo.Version),
				},
				Path:      file.Name(),
				Installed: true,
				Info:      nil, // Don't need this for now
				Deps:      nil, // TODO: This is not yet implemented, because the deps field is an expression, we need to parse it
			}
			return p
		}
	}

	return nil
}
