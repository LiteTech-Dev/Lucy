package sources

import (
	"fmt"
	"io"
	"lucy/syntax"
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

var AvailableSources = map[syntax.Platform][]Source{
	syntax.Fabric: {CurseForge, Modrinth},
	syntax.Forge:  {CurseForge, Modrinth},
	syntax.Mcdr:   {Mcdr},
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
func SelectSource(platform syntax.Platform) Source {
	min := slow
	fastestSource := None
	wg := sync.WaitGroup{}
	for _, source := range AvailableSources[platform] {
		wg.Add(1)
		go func() {
			defer wg.Done()
			speed := testDownloadSpeed(SpeedTestUrls[source])
			if speed < min {
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
