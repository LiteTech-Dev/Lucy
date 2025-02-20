package cmd

import (
	"errors"
	"github.com/urfave/cli/v3"
	"lucy/lucytypes"
)

var flagJsonOutput = &cli.BoolFlag{
	Name:  "json",
	Usage: "Print raw JSON response",
	Value: false,
}

var flagLongOutput = &cli.BoolFlag{
	Name:    "long",
	Usage:   "Show hidden or collapsed output",
	Value:   false,
	Aliases: []string{"a"},
}

func sourceFlag(absent lucytypes.Source) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "source",
		Aliases: []string{"s"},
		Usage:   "To fetch info from `SOURCE`",
		Value:   absent.String(),
		Validator: func(s string) error {
			if s != lucytypes.UnknownSource.String() {
				return errors.New("unknown source " + s)
			}
			return nil
		},
	}
}
