package subcommands

import (
	"log/slog"

	"github.com/spetix/otplet/internal/events"
	"github.com/spf13/cobra"
)

func NewLockCommand() *cobra.Command {
	var otpStore string
	var recipient string

	lock := &cobra.Command{
		Use:   "lock",
		Short: "lock GPG key",
		Run: func(cmd *cobra.Command, args []string) {
			cm := events.ClickManager{}
			cm.HandleClickEvents("3", otpStore, recipient)
			slog.Info("GPG key locked")
		},
	}

	lock.Flags().StringVarP(&recipient, "recipient", "r", "", "GPG recipient to decrypt OTP store")
	lock.Flags().StringVarP(&otpStore, "store", "s", "otp.json", "path to otp store")

	slog.Info("Lock command registered")
	return lock
}
