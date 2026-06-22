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

			em := setupBlocklet.GetEventManager()
			em.Register(models.LeftButton, func() error {
				slog.Info("Left click event triggered")
				return events.TriggerGPGPopup(recipient)
			})
			em.Register(models.RightButton, func() error {
				slog.Info("Right click event triggered")
				return events.TriggerLock()
			})

			otpStore := cmd.Flag("store").Value.String()
			var pv provider.OtpProviderItf
			err := em.Run()
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
