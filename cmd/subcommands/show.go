package subcommands

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spetix/bar-out-adapters/pkg/barout"
	"github.com/spetix/otplet/internal/events"
	"github.com/spetix/otplet/internal/provider"
	"github.com/spetix/otplet/internal/render"
	"github.com/spetix/otplet/internal/secretmanager"
	"github.com/spf13/cobra"
)

func NewShowCommand() *cobra.Command {
	var recipient string
	var proto string
	var renderOptions render.RenderOptions

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

			b := barout.New(proto)

			if pv == nil {
				return fmt.Errorf("OTP provider not found")
			}
			data := render.OtpDataNew(pv, &renderOptions)

			b.Print(data)
			return nil
		},
	}

	showFlags := show.Flags()
	showFlags.StringVarP(&proto, "proto", "p", "raw", "proto")
	lbl := os.Getenv("label")
	if lbl == "" {
		lbl = "🎄"
	}
	showFlags.StringVarP(&renderOptions.Label, "label", "L", lbl, "label")
	showFlags.StringVarP(&renderOptions.Format, "format", "R", "text", "format")
	clr := os.Getenv("color")
	if clr == "" {
		clr = "#ff0000"
	}
	showFlags.StringVarP(&renderOptions.ForegroundColor, "color", "c", clr, "foreground color")
	bgclr := os.Getenv("background")
	if bgclr == "" {
		bgclr = "#000000"
	}
	showFlags.StringVarP(&renderOptions.BackgroundColor, "background", "b", "#000000", "background color")

	showFlags.StringVarP(&recipient, "recipient", "r", "", "GPG recipient to decrypt OTP store")

	slog.Info("Show command registered")
	return show
}
