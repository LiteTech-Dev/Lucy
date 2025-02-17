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

func primaryFile(files []datatypes.ModrinthVersionFile) (primary datatypes.ModrinthVersionFile) {
	for _, file := range files {
		if file.Primary {
			return file
		}
	}
	return files[0]
}
