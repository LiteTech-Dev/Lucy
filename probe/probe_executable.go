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
	"lucy/syntax"
	"lucy/tools"
	"os"
	"path"
	"strings"
)

var getServerExecutable = tools.Memoize(
	func() lucytypes.ServerExecutable {
		var suspectedExecutables []*lucytypes.ServerExecutable
		for _, jarFile := range findJarFiles(getServerWorkPath()) {
			if exec := analyzeServerExecutable(jarFile); exec != nil {
				suspectedExecutables = append(suspectedExecutables, exec)
			}
		}

		switch len(suspectedExecutables) {
		case 0:
			logger.CreateFatal(errors.New("no server executable found"))
			return lucytypes.ServerExecutable{} // unreachable
		case 1:
			return *suspectedExecutables[0]
		default:
			index := output.PromptSelectExecutable(suspectedExecutables)
			return *suspectedExecutables[index]
		}
	},
)

func findJarFiles(dir string) (jarFiles []string) {
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if path.Ext(entry.Name()) == ".jar" {
			jarFiles = append(jarFiles, path.Join(dir, entry.Name()))
		}
	}
	return
}

func analyzeServerExecutable(executableFile string) (serverExecutable *lucytypes.ServerExecutable) {
	serverExecutable = &lucytypes.ServerExecutable{}
	serverExecutable.Path = executableFile
	zipReader, _ := zip.OpenReader(executableFile)
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(zipReader)

	for _, f := range zipReader.File {
		switch f.Name {
		case fabricAttributeFileName:
			serverExecutable.Type = syntax.Fabric
			r, _ := f.Open()
			data, _ := io.ReadAll(r)
			serverExecutable.GameVersion = strings.Split(
				strings.Split(
					string(data), "\n",
				)[1], "=",
			)[1]
			serverExecutable.GameVersion = strings.Split(
				strings.Split(
					string(data), "\n",
				)[0], "=",
			)[1]
			return
		case vanillaAttributeFileName:
			versionDotJson := JarVersionDotJson{}
			serverExecutable.Type = syntax.Minecraft
			r, _ := f.Open()
			data, _ := io.ReadAll(r)
			_ = json.Unmarshal(data, &versionDotJson)
			serverExecutable.GameVersion = versionDotJson.Id
			return
		}
	}

	return nil
}
