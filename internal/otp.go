package otp

import (
	"encoding/json"
	"os"
	"time"

	otp "github.com/pquerna/otp"
	totp "github.com/pquerna/otp/totp"
)

type PassCode struct {
	Code       string
	ValidUntil time.Time
	Error      error
}

func (p PassCode) IsValid() bool {
	return p.Error == nil && time.Now().Before(p.ValidUntil)
}

func LoadTokenFromFile(path string) (*otp.Key, error) {
	hdlr, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer hdlr.Close()

	var key *otp.Key
	if err := json.NewDecoder(hdlr).Decode(&key); err != nil {
		return nil, err
	}

	return key, nil
}

func SaveTokenToFile(path string, key *otp.Key) error {
	hdlr, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer hdlr.Close()

	if err := json.NewEncoder(hdlr).Encode(key); err != nil {
		return err
	}

	return nil
}

func GetOtpCode(key *otp.Key) PassCode {
	now := time.Now().UTC()
	remaining := key.Period() - (uint64(now.Second()) % key.Period())
	validUntil := now.Add(time.Duration(remaining) * time.Second)
	code, err := totp.GenerateCodeCustom(key.Secret(), time.Now().UTC(), totp.ValidateOpts{
		Period:    uint(key.Period()),
		Skew:      1,
		Digits:    key.Digits(),
		Algorithm: key.Algorithm(),
	})
	return PassCode{
		Code:       code,
		ValidUntil: validUntil,
		Error:      err,
	}
}
