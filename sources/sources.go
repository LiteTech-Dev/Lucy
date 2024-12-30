package sources

import (
	"fmt"
	"io"
	"lucy/syntax"
	"net/http"
	"os"
	"strconv"
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

// MultiSourceDownload expects the urls hosts the same file. However, it does
// not verify the checksums to allow more loose file recognition policies in its
// callers.
//
// Download is concurrent. Other threads will be cancelled when one thread
// complete downloaded winThreshold of the file.
//
// Note that if the urls' speed are close, urls[0] will be selected since its
// goroutine is started first.
//
// Pros:
//   - Guaranteed to download the file from the fastest source.
//
// Cons:
//   - Wastes bandwidth
func MultiSourceDownload(urls []string, path string) {
	const winThreshold = 0.2 // 20% of the file
	var wg sync.WaitGroup
	var mu sync.Mutex
	var win bool
	var data *[]byte
	var winUrl string

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			// TODO: totalSize might be -1 when size it not known, handle this case
			totalSize := resp.ContentLength
			thresholdSize := int64(float64(totalSize) * winThreshold)
			buffer := make([]byte, 2048)
			var downloadedSize int64

			for {
				n, err := resp.Body.Read(buffer)
				if err != nil && err != io.EOF {
					return
				}
				if n == 0 {
					break
				}
				downloadedSize += int64(n)
				if win && winUrl != url {
					println(
						"canceling:",
						url,
						"("+strconv.FormatInt(downloadedSize, 10)+"/"+
							strconv.FormatInt(totalSize, 10), "bytes)",
					)
					return
				}
				if downloadedSize >= thresholdSize {
					mu.Lock()
					if !win {
						println(
							"winning:",
							url,
							"("+strconv.FormatInt(downloadedSize, 10)+"/"+
								strconv.FormatInt(totalSize, 10), "bytes)",
						)
						win = true
						data = &buffer
						winUrl = url
					}
					mu.Unlock()
				}
			}
		}(url)
	}

	wg.Wait()
	println("winning url: ", winUrl)

	file, _ := os.Create(path)
	defer file.Close()
	_, err := file.Write(*data)
	if err != nil {
		panic(err)
	}

	println("Downloaded to", path)
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
