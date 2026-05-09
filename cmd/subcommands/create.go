package subcommands

import (
	"log/slog"

	"github.com/spetix/otplet/internal/secretmanager"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	var url string
	var path string
	var recipient string

	create := &cobra.Command{
		Use:   "create",
		Short: "create otp setup",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Info("Creating secure seed storage")
			sm := secretmanager.NewSecretManager(path, recipient)
			err := sm.SaveGenerator(url)
			if err != nil {
				slog.Debug("used uri %s", url)
				return err
			}
			slog.Info("Imported generator seed encrypted by %s in %s", recipient, path)
			return nil
		},
	}

	createFlags := create.Flags()
	createFlags.StringVarP(&url, "url", "u", "", "qr code decoded url")
	createFlags.StringVarP(&path, "path", "p", "~/.secrets", "path to save otp")
	createFlags.StringVarP(&recipient, "recipient", "r", "", "GPG recipient to encrypt OTP store")
	slog.Info("Create command registered")
	return create
}
