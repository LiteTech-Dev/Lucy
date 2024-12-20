package util

import (
	"github.com/schollz/progressbar/v3"
	"io"
	"lucy/probe"
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

	out, err := os.Create(path.Join(LucyPath, "downloads", subdir, filename))
	if os.IsNotExist(err) {
		os.MkdirAll(path.Join(LucyPath, "downloads", subdir), os.ModePerm)
		out, _ = os.Create(path.Join(LucyPath, "downloads", subdir, filename))
	}
	defer out.Close()

	res, _ := http.Get(url)
	defer res.Body.Close()

	bar := progressbar.DefaultBytes(
		res.ContentLength,
		"Downloading",
	)

	writer := io.MultiWriter(out, bar)
	io.Copy(writer, res.Body)

	return out
}
