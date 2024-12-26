package sources

import (
	"fmt"
	"io"
	"lucy/syntax"
	"net/http"
	"time"
)

type Source string

const (
	CurseForge Source = "curseforge"
	Modrinth   Source = "modrinth"
	GitHub     Source = "github"
	None       Source = "none"
)

var AvailableSources = map[syntax.Platform][]Source{
	syntax.Fabric: {CurseForge, Modrinth},
	syntax.Forge:  {CurseForge, Modrinth},
}

var SpeedTestUrls = map[Source]string{
	CurseForge: "https://mediafilez.forgecdn.net/files/4834/896/fabric-api-0.87.2%2B1.19.4.jar",
	Modrinth:   "https://cdn.modrinth.com/data/P7dR8mSH/versions/nyAmoHlr/fabric-api-0.87.2%2B1.19.4.jar",
}

const slow = 2e31

func SelectSource(platform syntax.Platform) Source {
	min := slow
	fastestSource := None
	for _, source := range AvailableSources[platform] {
		speed := testDownloadSpeed(SpeedTestUrls[source])
		if speed < min {
			fastestSource = source
		}
		fmt.Printf("Speed for %s: %f\n", source, speed)
	}
	if fastestSource == None {
		panic("No available source")
	}
	return fastestSource
}

func testDownloadSpeed(url string) (elapsedTime float64) {
	resp, err := http.Get(url)
	if err != nil {
		return slow
	}
	defer resp.Body.Close()

	chunkSize := 2048
	startTime := time.Now()

	buffer := make([]byte, chunkSize)
	totalDownloaded := 0
	for i := 0; i < 10; i++ {
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			return slow
		}
		totalDownloaded += n
		if n == 0 {
			break
		}
	}

	elapsedTime = time.Since(startTime).Seconds()
	return elapsedTime
}
