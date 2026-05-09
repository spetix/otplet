package subcommands

import (
	"fmt"
	"log/slog"

	"github.com/spetix/otplet/internal/events"
	"github.com/spf13/cobra"
)

func NewUnlockCommand() *cobra.Command {
	var otpStore string
	var recipient string
	unlock := &cobra.Command{
		Use:   "unlock",
		Short: "unlock GPG key",
		Run: func(cmd *cobra.Command, args []string) {
			cm := events.ClickManager{}
			cm.HandleClickEvents("1", otpStore, recipient)
			fmt.Println("GPG key unlocked")
		},
	}

	unlock.Flags().StringVarP(&recipient, "recipient", "r", "", "GPG recipient to decrypt OTP store")
	unlock.Flags().StringVarP(&otpStore, "store", "s", "otp.json", "path to otp store")
	slog.Info("Unlock command registered")
	return unlock
}
