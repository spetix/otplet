// Package onetimepass defines the one-time password value and its validity
// metadata used by the blocklet runtime.
package onetimepass

import "time"

// PassCode holds a generated OTP, its expiration time, and any generation error.
type PassCode struct {
	Code       string
	ValidUntil time.Time
	Remaining  time.Duration
	Error      error
}

// IsValid returns true when the OTP code is present, has no error, and is still valid.
func (p PassCode) IsValid() bool {
	return p.Error == nil && time.Now().Before(p.ValidUntil)
}
