package main

import (
	"context"
	"lucy/cmd"
	"lucy/logger"
	"os"
)

func main() {
	cmd.Cli.Run(context.Background(), os.Args)
	logger.WriteAll()
}
