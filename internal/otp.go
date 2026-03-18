package otp

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	otplib "github.com/pquerna/otp"
	totp "github.com/pquerna/otp/totp"
)

type PassCode struct {
	Code       string
	ValidUntil time.Time
	Remaining  time.Duration
	Error      error
}

func (p PassCode) IsValid() bool {
	return p.Error == nil && time.Now().Before(p.ValidUntil)
}

func LoadTokenFromFile(path string) (*otplib.Key, error) {
	hdlr, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer hdlr.Close()

	var keyUrl string
	if err := json.NewDecoder(hdlr).Decode(&keyUrl); err != nil {
		return nil, err
	}

	key, err := otplib.NewKeyFromURL(keyUrl)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func SaveGenerator(url string, path string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	log.Default().Printf("Importing OTP generator with URL: %s to path: %s\n", url, path)
	key, err := otplib.NewKeyFromURL(url)
	if err != nil {
		return err
	}
	return SaveTokenToFile(path, key)
}

func SaveTokenToFile(path string, key *otplib.Key) error {
	hdlr, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer hdlr.Close()

	enc := json.NewEncoder(hdlr)

	if err := enc.Encode(key.URL()); err != nil {
		return err
	}
	log.Default().Printf("OTP generator saved to %s\n", path)

	return nil
}

func getTime(period uint64) (time.Time, time.Duration) {
	now := time.Now().UTC()
	remaining := period - (uint64(now.Second()) % period)
	return now, time.Duration(remaining) * time.Second
}

func (p *OtpProvider) otpCode() {
	now, remaining := getTime(p.Key.Period())
	validUntil := now.Add(remaining)
	code, err := totp.GenerateCodeCustom(p.Key.Secret(), time.Now().UTC(), totp.ValidateOpts{
		Period:    uint(p.Key.Period()),
		Skew:      1,
		Digits:    p.Key.Digits(),
		Algorithm: p.Key.Algorithm(),
	})
	codeVal := &PassCode{
		Code:       code,
		ValidUntil: validUntil,
		Remaining:  remaining,
		Error:      err,
	}

	p.Code.Store(codeVal)
}

type OtpProvider struct {
	Key  *otplib.Key
	Code atomic.Pointer[PassCode]
}

func NewOtpProvider(key *otplib.Key) *OtpProvider {
	provider := &OtpProvider{
		Key: key,
		// atomic.Pointer zero value is fine, but we populate it immediately
		Code: atomic.Pointer[PassCode]{},
	}

	// generate an initial code synchronously so callers don't get nil
	provider.otpCode()

	// start the background refresher now that we have a valid code
	go provider.start()
	return provider
}

func (p *OtpProvider) start() {
	period := time.Duration(p.Key.Period()) * time.Second
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for range ticker.C {
		p.otpCode()
	}
}

func (p *OtpProvider) GetCode() *PassCode {
	return p.Code.Load()
}
