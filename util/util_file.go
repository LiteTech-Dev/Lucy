package util

import (
	"fmt"
	"lucy/probe"
	"os"
	"path"
)

func InstallMod(file *os.File) {
	filename := path.Base(file.Name())
	fmt.Println("Installing mod", filename)
	serverInfo := probe.GetServerInfo()
	os.Rename(
		file.Name(),
		path.Join(serverInfo.ModPath, filename),
	)
}
