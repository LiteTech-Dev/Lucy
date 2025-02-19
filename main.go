package main

import (
	"context"
	"os"

	"lucy/cmd"
	"lucy/logger"
)

func main() {
	defer func() {
		logger.Debug("program finished with exit code 0")
		logger.WriteAll()
	}()
	cmd.Cli.Run(context.Background(), os.Args)
}
