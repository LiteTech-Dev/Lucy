package main

import (
	"context"
	"log"
	"lucy/cli"
	"os"
)

func main() {
	if err := cli.Cli.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
