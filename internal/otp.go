package otp

import (
	"encoding/json"
	"os"
	"sync/atomic"
	"time"

	"github.com/pquerna/otp"
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

	var key *otp.Key
	if err := json.NewDecoder(hdlr).Decode(&key); err != nil {
		return nil, err
	}

	return key, nil
}

func SaveTokenToFile(path string, key *otplib.Key) error {
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
		Key:  key,
		Code: atomic.Pointer[PassCode]{},
	}

	go provider.start()
	return provider
}

func (p *OtpProvider) start() {
	for {
		p.otpCode()
		time.Sleep(p.Code.Load().Remaining)
	}
}

func (p *OtpProvider) GetCode() *PassCode {
	return p.Code.Load()
}
