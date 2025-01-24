package probe

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"log"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/output"
	"lucy/syntaxtypes"
	"os"
	"path"
	"strings"
)

// TODO: Improve probe logic, plain executable unpacking do not work well
// TODO: Research on forge installation

func getExecutableInfo() lucytypes.ExecutableInfo {
	var valid []lucytypes.ExecutableInfo
	workPath := getServerWorkPath()
	jars := findJar(workPath)
	for _, jar := range jars {
		exec := analyzeExecutable(jar)
		if exec != nil {
			valid = append(valid, *exec)
		}
	}

	switch len(valid) {
	case 0:
		logger.CreateFatal(errors.New("no server executable found"))
	case 1:
		return valid[0]
	default:
		index := output.PromptSelectExecutable(valid)
		return valid[index]
	}

	return lucytypes.ExecutableInfo{}
}

func findJar(dir string) (jarFiles []*os.File) {
	jarFiles = []*os.File{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if path.Ext(entry.Name()) == ".jar" {
			file, err := os.Open(path.Join(dir, entry.Name()))
			if err != nil {
				logger.CreateWarning(err)
				continue
			}
			jarFiles = append(jarFiles, file)
		}
	}

	return
}

var unknownExecutable = lucytypes.ExecutableInfo{
	Path:        "",
	GameVersion: "unknown",
	BootCommand: nil,
	Type:        "unknown",
}

const fabricSingleIdentifierFile = "install.properties"
const vanillaIdentifierFile = "version.json"
const fabricLauncherIdentifierFile = "fabric-server-launch.properties"
const fabricLauncherManifest = "META-INF/MANIFEST.MF"

func analyzeExecutable(file *os.File) (exec *lucytypes.ExecutableInfo) {
	// exec is a nil before an analysis function is called
	// Anything other than exec.Path is set in the analysis function
	stat, _ := file.Stat()
	reader, _ := zip.NewReader(file, stat.Size())

	for _, f := range reader.File {
		switch f.Name {
		case fabricSingleIdentifierFile:
			if exec != nil {
				*exec = unknownExecutable
				exec.Path = file.Name()
				return
			}
			exec = analyzeFabricSingle(f)
		case fabricLauncherIdentifierFile:
			if exec != nil {
				*exec = unknownExecutable
				exec.Path = file.Name()
				return
			}
			for _, ff := range reader.File {
				if ff.Name == fabricLauncherManifest {
					exec = analyzeFabricLauncher(ff)
				}
			}
		case vanillaIdentifierFile:
			if exec != nil {
				*exec = unknownExecutable
				exec.Path = file.Name()
				return
			}
			exec = analyzeVanilla(f)
		}
	}

	if exec == nil {
		exec = &unknownExecutable
		return
	}
	// Set the path to the file at the end
	exec.Path = file.Name()

	return
}

func analyzeVanilla(versionJson *zip.File) (exec *lucytypes.ExecutableInfo) {
	exec = &lucytypes.ExecutableInfo{}
	exec.Type = syntaxtypes.Minecraft
	reader, _ := versionJson.Open()
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
	exec.Type = syntaxtypes.Fabric
	r, _ := installProperties.Open()
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
	exec.Type = syntaxtypes.Fabric
	r, _ := manifest.Open()
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
