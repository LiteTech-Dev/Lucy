package local

import (
	"archive/zip"
	"encoding/json"
	"io"
	"os"
	"path"
	"strings"

	"lucy/logger"
	"lucy/lucytypes"
	"lucy/output"
	"lucy/tools"
)

// TODO: Improve probe logic, plain executable unpacking do not work well
// TODO: Research on forge installation

var getExecutableInfo = tools.Memoize(
	func() *lucytypes.ExecutableInfo {
		var valid []*lucytypes.ExecutableInfo
		workPath := getServerWorkPath()
		jars := findJar(workPath)
		for _, jar := range jars {
			exec := analyzeExecutable(jar)
			if exec == nil {
				continue
			}
			valid = append(valid, exec)
		}

		if len(valid) == 0 {
			logger.Info("no server under current directory")
			return UnknownExecutable
		} else if len(valid) == 1 {
			return valid[0]
		}
		index := output.PromptSelectExecutable(valid)
		return valid[index]
	},
)

func findJar(dir string) (jarFiles []*os.File) {
	jarFiles = []*os.File{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		logger.Info("cannot read current directory, most local-related features will be disabled")
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if path.Ext(entry.Name()) == ".jar" {
			file, err := os.Open(path.Join(dir, entry.Name()))
			if err != nil {
				logger.Warning(err)
				continue
			}
			jarFiles = append(jarFiles, file)
		}
	}
	return
}

var UnknownExecutable = &lucytypes.ExecutableInfo{
	Path:        "",
	GameVersion: "unknown",
	BootCommand: nil,
	Platform:    lucytypes.UnknownPlatform,
}

const (
	fabricSingleIdentifierFile   = "install.properties"
	vanillaIdentifierFile        = "version.json"
	fabricLauncherIdentifierFile = "fabric-server-launch.properties"
	fabricLauncherManifest       = "META-INF/MANIFEST.MF"
)

// analyzeExecutable gives nil if the jar file is invalid. The constant UnknownExecutable
// is not yet used in the codebase, however still reserved for future use.
func analyzeExecutable(file *os.File) (exec *lucytypes.ExecutableInfo) {
	// exec is a nil before an analysis function is called
	// Anything other than exec.Path is set in the analysis function
	stat, err := file.Stat()
	if err != nil {
		return nil
	}
	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return nil
	}

	for _, f := range reader.File {
		switch f.Name {
		case fabricSingleIdentifierFile:
			if exec != nil {
				return nil
			}
			exec = analyzeFabricSingle(f)
		case fabricLauncherIdentifierFile:
			if exec != nil {
				return nil
			}
			for _, ff := range reader.File {
				if ff.Name == fabricLauncherManifest {
					exec = analyzeFabricLauncher(ff)
				}
			}
		case vanillaIdentifierFile:
			if exec != nil {
				return nil
			}
			exec = analyzeVanilla(f)
		}
	}

	if exec == nil {
		return
	}
	// Set the path to the file at the end
	exec.Path = file.Name()
	return
}

func analyzeVanilla(versionJson *zip.File) (exec *lucytypes.ExecutableInfo) {
	exec = &lucytypes.ExecutableInfo{}
	exec.Platform = lucytypes.Minecraft
	reader, _ := versionJson.Open()
	defer tools.CloseReader(reader, logger.Warning)
	data, _ := io.ReadAll(reader)
	obj := VersionDotJson{}
	_ = json.Unmarshal(data, &obj)
	exec.GameVersion = obj.Id
	return
}

// install.properties looks like this:
// fabric-loader-version=0.16.9
// game-version=1.21.4

func analyzeFabricSingle(installProperties *zip.File) (exec *lucytypes.ExecutableInfo) {
	exec = &lucytypes.ExecutableInfo{}
	exec.Platform = lucytypes.Fabric
	r, _ := installProperties.Open()
	defer tools.CloseReader(r, logger.Warning)
	data, _ := io.ReadAll(r)
	s := string(data)

	// Read second line, split by "=" and get the second part
	exec.GameVersion = strings.Split(strings.Split(s, "\n")[1], "=")[1]

	// Read first line, split by "=" and get the second part
	exec.LoaderVersion = strings.Split(strings.Split(s, "\n")[0], "=")[1]

	return
}

// META-INF/MANIFEST.MF looks like this:
// Manifest-Version: 1.0
// Main-Class: net.fabricmc.loader.impl.launch.server.FabricServerLauncher
// Class-Path: libraries/org/ow2/asm/asm/9.7.1/asm-9.7.1.jar libraries/org/
// ow2/asm/asm-analysis/9.7.1/asm-analysis-9.7.1.jar libraries/org/ow2/asm
// /asm-commons/9.7.1/asm-commons-9.7.1.jar libraries/org/ow2/asm/asm-tree
// /9.7.1/asm-tree-9.7.1.jar libraries/org/ow2/asm/asm-util/9.7.1/asm-util
// -9.7.1.jar libraries/net/fabricmc/sponge-mixin/0.15.4+mixin.0.8.7/spong
// e-mixin-0.15.4+mixin.0.8.7.jar libraries/net/fabricmc/intermediary/1.21
// .4/intermediary-1.21.4.jar libraries/net/fabricmc/fabric-loader/0.16.9/
// fabric-loader-0.16.9.jar
// Note that line breaks are "\r\n " and the last line ends with "\r\n"

func analyzeFabricLauncher(
	manifest *zip.File,
) (exec *lucytypes.ExecutableInfo) {
	exec = &lucytypes.ExecutableInfo{}
	exec.Platform = lucytypes.Fabric
	r, _ := manifest.Open()
	defer tools.CloseReader(r, logger.Warning)
	data, _ := io.ReadAll(r)
	s := string(data)
	s = strings.Split(s, "Class-Path: ")[1] // Start reading from Class-Path
	s = strings.ReplaceAll(s, "\r\n ", "")  // Remove line breaks
	s = strings.ReplaceAll(s, "\r\n", "")   // Remove last line breaks
	classPaths := strings.Split(s, " ")
	for _, classPath := range classPaths {
		if strings.Contains(classPath, "libraries/net/fabricmc/intermediary") {
			exec.GameVersion = strings.Split(classPath, "/")[4]
		}
		if strings.Contains(classPath, "libraries/net/fabricmc/fabric-loader") {
			exec.LoaderVersion = strings.Split(classPath, "/")[4]
		}
	}
	return
}
