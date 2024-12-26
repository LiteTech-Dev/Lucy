package main

import (
	"context"
	"lucy/cmd"
	"lucy/logger"
	"lucy/sources"
	"lucy/syntax"
	"os"
)

func main() {
	sources.SelectSource(syntax.Fabric)
	cmd.Cli.Run(context.Background(), os.Args)
	logger.WriteAll()
}
