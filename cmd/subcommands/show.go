package subcommands

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spetix/bar-out-adapters/pkg/barout"
	"github.com/spetix/otplet/internal/events"
	"github.com/spetix/otplet/internal/render"
	"github.com/spf13/cobra"
)

func NewShowCommand() *cobra.Command {
	var recipient string
	var proto string
	var otpStore string
	var renderOptions render.RenderOptions

	show := &cobra.Command{
		Use:   "show",
		Short: "show current otp",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Handle i3blocks click events
			eventId := os.Getenv("BLOCK_BUTTON")

			cm := events.ClickManager{}
			sm := cm.HandleClickEvents(eventId, otpStore, recipient)

			b := barout.New(proto)
			otpProvider, err := sm.LoadTokenFromFile()
			if err != nil {
				return err
			}

			if otpProvider == nil {
				return fmt.Errorf("OTP provider not found")
			}
			data := render.OtpDataNew(otpProvider, &renderOptions)

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
	showFlags.StringVarP(&renderOptions.Label, "label", "l", lbl, "label")
	showFlags.StringVarP(&renderOptions.Format, "format", "f", "text", "format")
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

	showFlags.StringVarP(&otpStore, "store", "s", "otp.json", "path to otp store")
	showFlags.StringVarP(&recipient, "recipient", "r", "", "GPG recipient to decrypt OTP store")

	slog.Info("Show command registered")
	return show
}
