// Package render exposes OTP values to the bar-out adapter used by the blocklet.
package render

import (
	"fmt"

	"github.com/spetix/bar-out-adapters/pkg/barout/data"
	"github.com/spetix/bar-out-adapters/pkg/barout/models"
	"github.com/spetix/otplet/internal/provider"
)

// OtpData implements the bar-out data interface for OTP display.
type OtpData struct {
	data.Data
	otp       provider.OtpProviderItf
	options   models.RenderOptions
	formatter models.Formatter
}

func (d *OtpData) Short() string {
	if d.otp == nil {
		return "No OTP configured"
	}
	code := d.otp.GetCode()
	if code == nil || code.Error != nil {
		return "Error generating OTP"
	}
	return code.Code
}

// Long returns a detailed OTP string for the bar output, optionally using a
// configured format string.
func (d *OtpData) Long() string {
	if d.otp == nil {
		return "No OTP configured"
	}
	code := d.otp.GetCode()
	if code == nil || code.Error != nil {
		return "Error generating OTP"
	}

	if d.formatter != nil {
		return d.formatter.Render(code.Code, fmt.Sprint(int(code.Remaining.Seconds())))
	}

	return code.Code
}

func (d *OtpData) BackgroundColor() string {
	return d.options.BackgroundColor()
}

func (d *OtpData) ForegroundColor() string {
	return d.options.ForegroundColor()
}

func (d *OtpData) Label() string {
	return d.options.Label()
}

// OnClick returns the command to execute when the blocklet is left-clicked.
// This unlocks the GPG key so it won't prompt for passphrase on next use.
func (d *OtpData) OnClick() string {
	// The command runs gpg-connect-agent to cache the passphrase
	return "gpg-connect-agent updatestatus /bye"
}

// OnRightClick returns the command to execute when the blocklet is right-clicked.
// This blocks the GPG key by killing the agent, clearing the cached passphrase.
func (d *OtpData) OnRightClick() string {
	// The command kills the gpg-agent to clear cached passphrases
	return "gpgconf --kill gpg-agent"
}

// OtpDataNew constructs a new render data instance for the given provider and
// render options.
func OtpDataNew(otp provider.OtpProviderItf, renderOptions models.RenderOptions) *OtpData {
	return &OtpData{
		options: renderOptions,
		otp:     otp,
	}
}
