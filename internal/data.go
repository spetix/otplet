package otp

import (
	"fmt"

	"github.com/spetix/bar-out-adapters/pkg/barout/data"
)

type RenderOptions struct {
	Label           string
	Format          string
	ForegroundColor string
	BackgroundColor string
}

type OtpData struct {
	data.Data
	otp     *OtpProvider
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

func (d *OtpData) Long() string {
	code := d.otp.GetCode()
	if code == nil || code.Error != nil {
		return "Error generating OTP"
	}
	return fmt.Sprintf(code.Code, code.Remaining.Seconds())
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

func OtpDataNew(otp *OtpProvider, renderOptions *RenderOptions) *OtpData {
	return &OtpData{
		options: renderOptions,
		otp:     otp,
	}
}
