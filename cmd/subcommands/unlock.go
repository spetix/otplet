package subcommands

import (
	"log/slog"
	"os"

	"github.com/spetix/bar-out-adapters/pkg/barout/models"
	"github.com/spetix/otplet/internal/events"
	"github.com/spf13/cobra"
)

func NewUnlockCommand(setupBlocklet models.SetupBlocklet) *cobra.Command {
	var recipient string
	unlock := &cobra.Command{
		Use:   "unlock",
		Short: "unlock GPG key",
		Run: func(cmd *cobra.Command, args []string) {
			otpStore := cmd.Flag("store").Value.String()
			slog.Info("Unlocking GPG key for OTP store", "store", otpStore, "recipient", recipient)
			if err := events.TriggerGPGPopup(recipient); err != nil {
				slog.Error("Failed to unlock GPG key", "error", err)
				return
			}
			slog.Info("GPG key unlocked")
		},
	}

	recipient = os.Getenv("RECIPIENT")
	unlock.Flags().StringVarP(&recipient, "recipient", "r", recipient, "GPG recipient to decrypt OTP store")
	slog.Info("Unlock command registered")
	return unlock
}
