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

package modrinth

import (
	"lucy/datatypes"
	"lucy/lucytypes"
)

func GetFile(id lucytypes.PackageId) (url string, filename string, err error) {
	version, err := getVersion(id)
	if err != nil {
		return "", "", err
	}
	primary := primaryFile(version.Files)
	return primary.Url, primary.Filename, nil
}

func getFile(version *datatypes.ModrinthVersion) (url string, filename string) {
	primary := primaryFile(version.Files)
	return primary.Url, primary.Filename
}

func primaryFile(files []datatypes.ModrinthVersionFile) (primary datatypes.ModrinthVersionFile) {
	for _, file := range files {
		if file.Primary {
			return file
		}
	}
	return files[0]
}
