package subcommands

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spetix/bar-out-adapters/pkg/barout/models"
	"github.com/spetix/otplet/internal/events"
	"github.com/spetix/otplet/internal/provider"
	"github.com/spetix/otplet/internal/render"
	"github.com/spetix/otplet/internal/secretmanager"
	"github.com/spf13/cobra"
)

func NewShowCommand(setupBlocklet models.SetupBlocklet) *cobra.Command {
	var recipient string

	show := &cobra.Command{
		Use:   "show",
		Short: "show current otp",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Handle i3blocks click events
			eventId := os.Getenv("BLOCK_BUTTON")

			cm := events.ClickManager{}
			otpStore := cmd.Flag("store").Value.String()
			slog.Info("Using seed in store", "store", otpStore)
			err := cm.HandleClickEvents(eventId, otpStore, recipient)
			var pv provider.OtpProviderItf
			if err != nil {
				pv = provider.NewDummyProvider()
			} else {
				pv, err = secretmanager.LoadTokenFromFile(otpStore)
			}

			b := setupBlocklet.GetOutput()

			if pv == nil {
				return fmt.Errorf("OTP provider not found")
			}
			data := render.OtpDataNew(pv, setupBlocklet.Options())

			b.Print(data)
			return nil
		},
	}

	showFlags := show.Flags()
	recipient = os.Getenv("RECIPIENT")
	showFlags.StringVarP(&recipient, "recipient", "r", recipient, "GPG recipient to decrypt OTP store")

	slog.Info("Show command registered")
	return show
}
