package subcommands

import (
	"log/slog"

	"github.com/spetix/otplet/internal/events"
	"github.com/spf13/cobra"
)

func NewLockCommand() *cobra.Command {
	var recipient string

	lock := &cobra.Command{
		Use:   "lock",
		Short: "lock GPG key",
		Run: func(cmd *cobra.Command, args []string) {
			otpStore := cmd.Flag("store").Value.String()
			slog.Info("Locking GPG key for OTP store", "store", otpStore, "recipient", recipient)
			events.TriggerLock()
			slog.Info("GPG key locked")
		},
	}

	lock.Flags().StringVarP(&recipient, "recipient", "r", "", "GPG recipient to decrypt OTP store")

	slog.Info("Lock command registered")
	return lock
}
