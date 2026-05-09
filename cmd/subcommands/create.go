package subcommands

import (
	"log/slog"

	"github.com/spetix/otplet/internal/secretmanager"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	var url string
	var recipient string

	create := &cobra.Command{
		Use:   "create",
		Short: "create otp setup",
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Info("Creating secure seed storage")
			otpStore := cmd.Flag("store").Value.String()
			slog.Info("Using OTP store", "store", otpStore)
			err := secretmanager.SaveGenerator(url, otpStore, []string{recipient})
			if err != nil {
				slog.Debug("used url as seed", "url", url)
				return err
			}
			slog.Info("Imported generator seed encrypted", "recipient", recipient, "store", otpStore)
			return nil
		},
	}

	createFlags := create.Flags()
	createFlags.StringVarP(&url, "url", "u", "", "qr code decoded url")
	createFlags.StringVarP(&recipient, "recipient", "r", "", "GPG recipient to encrypt OTP store")
	slog.Info("Create command registered")
	return create
}
