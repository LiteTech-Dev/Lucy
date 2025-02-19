package logger

import (
	"os"
	"path/filepath"
	"runtime"
)

var LogFile = logFile()

func logDir() string {
	var logDir string

	switch runtime.GOOS {
	case "windows":
		logDir = filepath.Join(os.Getenv("APPDATA"), "lucy", "logs")
	case "darwin":
		logDir = filepath.Join(os.Getenv("HOME"), "Library", "Logs", "lucy")
	case "linux":
		logDir = filepath.Join(
			os.Getenv("HOME"),
			".local",
			"share",
			"lucy",
			"logs",
		)
	default:
		logDir = "./logs"
	}

	return logDir
}

func logFile() *os.File {
	logDir := logDir()
	err := os.MkdirAll(logDir, 0o755)
	if err != nil {
		devNull, _ := os.Open(os.DevNull)
		return devNull
	}

	logFilePath := filepath.Join(logDir, "lucy.log")
	logFile, err := os.OpenFile(
		logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o755,
	)
	if err != nil {
		println(err.Error())
		devNull, _ := os.Open(os.DevNull)
		return devNull
	}

	return logFile
}
