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

package util

import (
	"io"
	"os"
	"path"
)

func InstallLucy() {
	os.Mkdir(ProgramPath, 0o755)
	os.Mkdir(DownloadPath, 0o755)
	os.Mkdir(CachePath, 0o755)
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
