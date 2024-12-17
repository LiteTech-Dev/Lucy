package probe

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"strings"
)

func findJarFiles(dir string) (jarFiles []string) {
	jarFiles = make([]string, 0)
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

func analyzeServerExecutable(executableFile string) (
	gameVersion string,
	modLoaderType string,
	modLoaderVersion string,
) {
	zipReader, _ := zip.OpenReader(executableFile)
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {

		}
	}(zipReader)

	for _, f := range zipReader.File {
		switch f.Name {
		case fabricPropertiesFileName:
			modLoaderType = "fabric"
			executableReader, _ := f.Open()
			fabricPropertiesData, _ := io.ReadAll(executableReader)
			gameVersion = strings.Split(
				strings.Split(
					string(fabricPropertiesData), "\n",
				)[1], "=",
			)[1]
			modLoaderVersion = strings.Split(
				strings.Split(
					string(fabricPropertiesData), "\n",
				)[0], "=",
			)[1]
			return
		}
	}

	return "", "", ""
}
