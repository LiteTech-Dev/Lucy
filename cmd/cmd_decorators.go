package cmd

import (
	"context"
	"github.com/urfave/cli/v3"
	"lucy/logger"
)

// globalFlagsDecorator is a high-order function that appends global flag actions
// to the action function.
func globalFlagsDecorator(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() == 0 && len(cmd.FlagNames()) == 0 {
			cli.ShowAppHelpAndExit(cmd, 0)
		}
		if cmd.Bool("log-file") {
			println("Log file at", logger.LogFile.Name())
		}
		if cmd.Bool("verbose") {
			logger.UseConsoleOutput()
		}
		if cmd.Bool("debug") {
			logger.UseDebug()
		}
		return f(ctx, cmd)
	}
}

// helpOnNoInputDecorator is a high-order function that takes a cli.ActionFunc and
// returns a cli.ActionFunc that prints help and exit when there's no args and flags.
//
// This function is not necessarily used for every action function, as some
// action functions are expected to have no args and flags. E.g., `lucy status`.
func helpOnNoInputDecorator(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() == 0 && len(cmd.LocalFlagNames()) == 0 {
			cli.ShowAppHelpAndExit(cmd, 0)
		}
		return f(ctx, cmd)
	}
}

func helpOnErrorDecorator(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		err := f(ctx, cmd)
		if err != nil {
			cli.ShowAppHelpAndExit(cmd, 1)
		}
		return err
	}
}
