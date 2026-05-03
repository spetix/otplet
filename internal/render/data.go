// Package render exposes OTP values to the bar-out adapter used by the blocklet.
package render

import (
	"fmt"

	"github.com/spetix/bar-out-adapters/pkg/barout/data"
	"github.com/spetix/otplet/internal/provider"
)

// RenderOptions contains display configuration for the OTP blocklet output.
type RenderOptions struct {
	Label           string
	Format          string
	ForegroundColor string
	BackgroundColor string
}

// OtpData implements the bar-out data interface for OTP display.
type OtpData struct {
	data.Data
	otp     *provider.OtpProvider
	options *RenderOptions
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

	if d.options.Format != "" {
		return fmt.Sprintf(d.options.Format, code.Code, int(code.Remaining.Seconds()))
	}

	return code.Code
}

func (d *OtpData) BackgroundColor() string {
	return d.options.BackgroundColor
}

func (d *OtpData) ForegroundColor() string {
	return d.options.ForegroundColor
}

func (d *OtpData) Label() string {
	return d.options.Label
}

// OtpDataNew constructs a new render data instance for the given provider and
// render options.
func OtpDataNew(otp *provider.OtpProvider, renderOptions *RenderOptions) *OtpData {
	return &OtpData{
		options: renderOptions,
		otp:     otp,
	}
}
