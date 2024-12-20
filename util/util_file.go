package util

import (
	"fmt"
	"io"
	"lucy/probe"
	"os"
	"path"
)

func InstallLucy() {
	os.Mkdir(LucyPath, 0755)
	os.Mkdir(LucyDownloadDir, 0755)
	os.Mkdir(LucyCacheDir, 0755)
	// 	TODO: create empty config
}

func InstallMod(file *os.File) {
	preventExternalFile(file.Name())
	filename := path.Base(file.Name())
	fmt.Println("Installing mod", filename)
	serverInfo := probe.GetServerInfo()
	os.Rename(
		file.Name(),
		path.Join(serverInfo.ModPath, filename),
	)
}

func CopyToCache(file *os.File) {
	preventExternalFile(file.Name())
	filename := path.Base(file.Name())
	cacheFile, _ := os.Create(path.Join(LucyCacheDir, filename))
	_, _ = io.Copy(cacheFile, file)
}

// preventExternalFile panics when the passed in string do not contain .lucy in
// all its parent directories. This is used to prevent any unexpected access to
// externally manages files.
func preventExternalFile(file string) {
	for ; file == path.Dir(file); file = path.Dir(file) {
		if path.Base(file) == LucyPath {
			return
		}
	}
	panic("incorrect file")
}
