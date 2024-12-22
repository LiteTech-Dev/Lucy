package util

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/term"
	"io"
	"lucy/probe"
	"lucy/tools"
	"net/http"
	"os"
	"path"
)

// DownloadFile
// All downloaded files are stored in .lucy/downloads/{subdir}/{filename}
// Current policy for path is the slug of the package
func DownloadFile(
	url string,
	subdir string,
	filename string,
) (outFile *os.File) {
	serverInfo := probe.GetServerInfo()
	if !serverInfo.HasLucy {
		// This is a very bad implementation
		// Not sure whether I should check for Lucy's existence here
		// Maybe we should assume all callers have checked it??
		// However, that might be redundant
		// Another drawback is that we cannot provide specific error messages here
		// Maybe a last check here is necessary
		panic("Lucy is not installed")
	}

	out, err := os.Create(path.Join(LucyDownloadDir, subdir, filename))
	if os.IsNotExist(err) {
		os.MkdirAll(path.Join(LucyDownloadDir, subdir), os.ModePerm)
		out, _ = os.Create(path.Join(LucyDownloadDir, subdir, filename))
	}
	defer out.Close()

	res, _ := http.Get(url)
	defer res.Body.Close()

	fmt.Println("Downloading", url)

	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	bar := progressbar.NewOptions64(
		res.ContentLength,
		progressbar.OptionShowCount(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(
			progressbar.Theme{
				Saucer:        "[bold][magenta]█[reset]",
				SaucerHead:    "[bold][magenta]█[reset]",
				SaucerPadding: " ",
				BarStart:      "[bold][ [reset]",
				BarEnd:        "[bold] ][reset]",
			},
		),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(
			tools.Trenary(
				func() bool { return termWidth/3 > 40 },
				termWidth/3,
				40,
			),
		),
	)
	writer := io.MultiWriter(out, bar)
	io.Copy(writer, res.Body)
	fmt.Println()

	return out
}
