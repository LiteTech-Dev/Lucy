package main

import (
	"context"
	"log"
	"lucy/cmd"
	"os"
)

const LucyDir = ".lucy"

func main() {
	if err := cmd.Cli.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
