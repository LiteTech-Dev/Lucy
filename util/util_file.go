package util

import (
	"io"
	"os"
	"path"
)

func InstallLucy() {
	os.Mkdir(ProgramPath, 0755)
	os.Mkdir(DownloadPath, 0755)
	os.Mkdir(CachePath, 0755)
	// 	TODO: create empty config
}

func MoveFile(src *os.File, dest string) (err error) {
	err = os.Rename(src.Name(), dest)
	return
}

func CopyToCache(f *os.File) {
	filename := path.Base(f.Name())
	cacheFile, _ := os.Create(path.Join(CachePath, filename))
	_, _ = io.Copy(cacheFile, f)
}
