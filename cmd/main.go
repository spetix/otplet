package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spetix/otplet/cmd/subcommands"
	"github.com/spetix/otplet/internal/logger"
	"github.com/spf13/cobra"
)

func NewMainCommand() *cobra.Command {
	main := &cobra.Command{
		Use:   "otplet",
		Short: "manage your otps in a smart way",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logLevel, err := cmd.Flags().GetString("logLevel")
			if err != nil {
				return err
			}
			logFile := cmd.Flag("logFile").Value.String()

			logFormat, err := cmd.Flags().GetString("logFormat")
			if err != nil {
				return err
			}
			logger.SetupLogger(logLevel, logFile, logFormat)
			slog.Info("logger setup")
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Use --help to see available commands.")
		},
	}

	main.PersistentFlags().StringP("logLevel", "l", "info", "log level (default is info)")
	main.PersistentFlags().StringP("logFile", "f", "", "log file path (default is stdout)")
	main.PersistentFlags().StringP("logFormat", "F", "text", "log format (default is text, options are text and json)")
	main.PersistentFlags().StringP("store", "s", "otp.json", "path to otp store")

	return main
}

func main() {

	mainCmd := NewMainCommand()
	mainCmd.AddCommand(subcommands.NewCreateCommand())
	mainCmd.AddCommand(subcommands.NewShowCommand())
	mainCmd.AddCommand(subcommands.NewLockCommand())
	mainCmd.AddCommand(subcommands.NewUnlockCommand())

	if err := mainCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
