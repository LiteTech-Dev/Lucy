package remote

import (
	"fmt"
	"lucy/logger"
	"lucy/lucytypes"
)

func FetchSource(
	source lucytypes.Source,
	id lucytypes.PackageId,
) (remote *lucytypes.PackageRemote) {
	if source == lucytypes.Auto {
		source = SelectSource(id.Platform)
	}

	switch source {
	case lucytypes.Modrinth:
		// return modrinth.Fetch(id)
	default:
		logger.Fatal(fmt.Errorf("source fetch not supported yet:" + source.String()))
	}

	return nil // unreachable
}

func GetDependencies(
	source lucytypes.Source,
	id lucytypes.PackageId,
) *lucytypes.PackageDependencies {
	return nil
}

func GetInformation(
	source lucytypes.Source,
	id lucytypes.PackageId,
) *lucytypes.PackageInformation {
	return nil
}

func SearchForProject(
	source lucytypes.Source,
	query string,
) []lucytypes.PackageName {
	return nil
}
