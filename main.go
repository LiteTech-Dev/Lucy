package main

import (
	"context"
	"lucy/cmd"
	"lucy/logger"
	"os"
)

func main() {
	defer func() {
		logger.Debug("program finished with exit code 0")
		logger.WriteAll()
	}()
	cmd.Cli.Run(context.Background(), os.Args)
}
