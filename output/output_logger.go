package output

import (
	"fmt"
	"io"
	"lucy/types"
	"os"
	"path/filepath"
	"runtime"
)

var LogWriter = io.MultiWriter(os.Stderr, getLogFile())

func WriteLogItem(message *types.LogItem) {
	fmt.Println()
	_, _ = io.WriteString(LogWriter, formatLogItem(message))
	fmt.Println()
}

func formatLogItem(message *types.LogItem) string {
	switch message.Level {
	case 0:
		return cyan("[INFO] ") + message.Content.Error()
	case 1:
		return yellow("[WARNING] ") + message.Content.Error()
	case 2:
		return red("[ERROR] ") + message.Content.Error()
	case 3:
		return red("[FATAL] ") + message.Content.Error()
	}
	return "Wrong level"
}

func getLogDir() string {
	var logDir string

	switch runtime.GOOS {
	case "windows":
		logDir = filepath.Join(os.Getenv("APPDATA"), "MyApp", "logs")
	case "darwin":
		logDir = filepath.Join(os.Getenv("HOME"), "Library", "Logs", "MyApp")
	case "linux":
		logDir = filepath.Join(
			os.Getenv("HOME"),
			".local",
			"share",
			"MyApp",
			"logs",
		)
	default:
		logDir = "./logs"
	}

	return logDir
}

func getLogFile() *os.File {
	logDir := getLogDir()
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		devNull, _ := os.Open(os.DevNull)
		return devNull
	}

	logFilePath := filepath.Join(logDir, "app.log")
	logFile, err := os.OpenFile(
		logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		os.ModeAppend,
	)
	if err != nil {
		devNull, _ := os.Open(os.DevNull)
		return devNull
	}

	return logFile
}
