package probe

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log"
	"lucy/types"
	"os"
	"path"
	"strings"
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

func analyzeServerExecutable(executableFile string) *types.ServerExecutable {
	serverExecutable := types.ServerExecutable{}
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
			serverExecutable.ModLoaderType = "fabric"
			r, _ := f.Open()
			data, _ := io.ReadAll(r)
			serverExecutable.GameVersion = strings.Split(
				strings.Split(
					string(data), "\n",
				)[1], "=",
			)[1]
			serverExecutable.ModLoaderVersion = strings.Split(
				strings.Split(
					string(data), "\n",
				)[0], "=",
			)[1]
			return &serverExecutable
		case vanillaAttributeFileName:
			versionDotJson := types.VersionDotJson{}
			serverExecutable.ModLoaderType = "vanilla"
			r, _ := f.Open()
			data, _ := io.ReadAll(r)
			_ = json.Unmarshal(data, &versionDotJson)
			serverExecutable.GameVersion = versionDotJson.Id
			serverExecutable.ModLoaderVersion = ""
			return &serverExecutable
		}
	}

	return nil
}
