package output

import (
	"fmt"
	"io"
	"lucy/lucytypes"
	"lucy/tools"
	"os"
	"path/filepath"
	"runtime"
)

var LogFile = getLogFile()

func WriteLogItem(message *lucytypes.LogItem) {
	fmt.Println()
	_, _ = io.WriteString(os.Stderr, formatLogItem(message))
	_, _ = io.WriteString(LogFile, formatUncolorLogItem(message)+"\n")
	fmt.Println()
}

func formatLogItem(message *lucytypes.LogItem) string {
	switch message.Level {
	case 0:
		return tools.Cyan("[INFO] ") + message.Content.Error()
	case 1:
		return tools.Yellow("[WARNING] ") + message.Content.Error()
	case 2:
		return tools.Red("[ERROR] ") + message.Content.Error()
	case 3:
		return tools.Red("[FATAL] ") + message.Content.Error()
	}
	return tools.Dim("[UNKNOWN] ") + message.Content.Error()
}

func formatUncolorLogItem(message *lucytypes.LogItem) string {
	switch message.Level {
	case 0:
		return "[INFO] " + message.Content.Error()
	case 1:
		return "[WARNING] " + message.Content.Error()
	case 2:
		return "[ERROR] " + message.Content.Error()
	case 3:
		return "[FATAL] " + message.Content.Error()
	}
	return "[UNKNOWN] " + message.Content.Error()
}

func getLogDir() string {
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

func getLogFile() *os.File {
	logDir := getLogDir()
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		devNull, _ := os.Open(os.DevNull)
		return devNull
	}

	logFilePath := filepath.Join(logDir, "lucy.log")
	logFile, err := os.OpenFile(
		logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0755,
	)
	if err != nil {
		println(err.Error())
		devNull, _ := os.Open(os.DevNull)
		return devNull
	}

	return logFile
}
