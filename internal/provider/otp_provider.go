// Package provider generates and refreshes OTP codes from an OTP key.
package provider

import (
	"sync/atomic"
	"time"

	otplib "github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/spetix/otplet/internal/onetimepass"
)

type OtpProviderItf interface {
	GetCode() *onetimepass.PassCode
	Secret() string
}

type dummyOtpProvider struct{}

func NewDummyProvider() OtpProviderItf {
	return &dummyOtpProvider{}
}

func (d *dummyOtpProvider) GetCode() *onetimepass.PassCode {
	return &onetimepass.PassCode{
		Code:       "󱧚",
		ValidUntil: time.Now().Add(30 * time.Second),
		Remaining:  30 * time.Second,
		Error:      nil,
	}
}

func (d *dummyOtpProvider) Secret() string {
	return "dummy"
}

// OtpProvider generates time-based one-time passwords and keeps a current code
// available through atomic storage so render consumers can read updates safely.
type otpProvider struct {
	Key  *otplib.Key
	Code atomic.Pointer[onetimepass.PassCode]
}

func NewOtpProvider(key *otplib.Key) *otpProvider {
	provider := &otpProvider{
		Key: key,
		// atomic.Pointer zero value is fine, but we populate it immediately
		Code: atomic.Pointer[onetimepass.PassCode]{},
	}

	// generate an initial code synchronously so callers don't get nil
	provider.otpCode()

	// start the background refresher now that we have a valid code
	go provider.start()
	return provider
}

func (p *otpProvider) otpCode() {
	now, remaining := getTime(p.Key.Period())
	validUntil := now.Add(remaining)
	code, err := totp.GenerateCodeCustom(p.Key.Secret(), time.Now().UTC(), totp.ValidateOpts{
		Period:    uint(p.Key.Period()),
		Skew:      1,
		Digits:    p.Key.Digits(),
		Algorithm: p.Key.Algorithm(),
	})
	codeVal := &onetimepass.PassCode{
		Code:       code,
		ValidUntil: validUntil,
		Remaining:  remaining,
		Error:      err,
	}

	p.Code.Store(codeVal)
}

func (p *otpProvider) start() {
	period := time.Duration(p.Key.Period()) * time.Second
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for range ticker.C {
		p.otpCode()
	}
}

// GetCode returns the latest generated OTP code, or nil if generation failed.
func (p *otpProvider) GetCode() *onetimepass.PassCode {
	return p.Code.Load()
}

func (p *otpProvider) Secret() string {
	return p.Key.Secret()
}

// getTime returns the current UTC time and the remaining duration for the
// current OTP period.
func getTime(period uint64) (time.Time, time.Duration) {
	now := time.Now().UTC()
	remaining := period - (uint64(now.Second()) % period)
	return now, time.Duration(remaining) * time.Second
}
