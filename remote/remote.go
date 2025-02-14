package remote

import (
	"fmt"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/remote/modrinth"
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
		remote = &lucytypes.PackageRemote{
			Source:           lucytypes.Modrinth,
			RemoteId:         modrinth.GetProjectId(id.Name),
			FileUrl:          modrinth.GetNewestProjectVersion(id.Name).Files[0].Url,
			LatestVersionUrl: modrinth.GetNewestProjectVersion(id.Name).Files[0].Url,
		}
	default:
		logger.CreateFatal(fmt.Errorf("source fetch not supported yet:" + source.String()))
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
) []*lucytypes.PackageName {
	return nil
}
