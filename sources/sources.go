package sources

import (
	"fmt"
	"io"
	"lucy/apitypes"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/sources/modrinth"
	"lucy/syntaxtypes"
	"net/http"
	"sync"
	"time"
)

type Source string

const (
	CurseForge Source = "curseforge"
	Modrinth   Source = "modrinth"
	GitHub     Source = "github"
	Mcdr       Source = "mcdr"
	None       Source = "none"
)

var AvailableSources = map[syntaxtypes.Platform][]Source{
	syntaxtypes.Fabric: {CurseForge, Modrinth},
	syntaxtypes.Forge:  {CurseForge, Modrinth},
	syntaxtypes.Mcdr:   {Mcdr},
}

var SpeedTestUrls = map[Source]string{
	CurseForge: "https://mediafilez.forgecdn.net/files/4834/896/fabric-api-0.87.2%2B1.19.4.jar",
	Modrinth:   "https://cdn.modrinth.com/data/P7dR8mSH/versions/nyAmoHlr/fabric-api-0.87.2%2B1.19.4.jar",
}

const slow float64 = 0x7FF0000000000000 // inf

// SelectSource is an alternative to MultiSourceDownload. It fetches a fixed url
// from SpeedTestUrls and measures the download speed of each source. The source
// with the fastest download speed is returned.
//
// Pros:
//   - Fastest source can be stored for later use.
//   - Saves bandwidth
//
// Cons:
//   - Speed test might not be representative
func SelectSource(platform syntaxtypes.Platform) Source {
	slowest := slow
	fastestSource := None
	wg := sync.WaitGroup{}
	for _, source := range AvailableSources[platform] {
		wg.Add(1)
		go func() {
			defer wg.Done()
			speed := testDownloadSpeed(SpeedTestUrls[source])
			if speed < slowest {
				fastestSource = source
			}
			fmt.Printf("Speed for %s: %f\n", source, speed)
		}()
	}

	wg.Wait()
	if fastestSource == None {
		panic("No available source")
	}

	return fastestSource
}

func testDownloadSpeed(url string) (elapsedTime float64) {
	startTime := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return slow
	}
	defer resp.Body.Close()

	chunkSize := 2048

	buffer := make([]byte, chunkSize)
	for i := 0; i < 10; i++ {
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			return slow
		}
		if n == 0 {
			break
		}
	}

	elapsedTime = time.Since(startTime).Seconds()
	return elapsedTime
}

// Maybe conversion functions are not a bad idea (and they are widely applied over
// the project). I may refactor them into methods when structs from apitypes are
// moved to their local packages.

func CInfoFromModrinth(res *apitypes.ModrinthSearchResults) []lucytypes.Package {
	packages := make([]lucytypes.Package, len(res.Hits))
	for _, hit := range res.Hits {
		info := lucytypes.Package{
			Id: syntaxtypes.PackageId{
				Platform: syntaxtypes.AllPlatform, // Add some algorithm to determine the platform
				Name:     syntaxtypes.PackageName(hit.Slug),
				Version:  "",
			},
			Path:               "",    // don't need
			Installed:          false, // don't need
			Urls:               nil,   // don't need
			Name:               "",    // don't need
			Description:        "",    // don't need
			SupportedVersions:  nil,
			SupportedPlatforms: nil,
		}
		packages = append(packages, info)
	}
	return packages
}

// TODO: More param should be added here to implement search options However,
//  we don't have a unified search option data structure yet. It should be designed
//  first.

func Search(
source Source,
keyword syntaxtypes.PackageName,
) []lucytypes.Package {
	switch source {
	case Modrinth:
		res := modrinth.Search(
			syntaxtypes.AllPlatform,
			keyword,
			false,
			"relevance",
		)
		return CInfoFromModrinth(res)
	case CurseForge:
		// TODO: curseforge search
		logger.CreateError(fmt.Errorf("curseforge not yet supported"))
	case GitHub:
		logger.CreateError(fmt.Errorf("github search not yet supported"))
	}
	return nil
}
