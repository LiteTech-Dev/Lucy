package main

import (
	"context"
	"lucy/cmd"
	"lucy/output"
	"os"
)

func main() {
	cmd.Cli.Run(context.Background(), os.Args)
	output.PrintMessagesAndExit(0)
}
