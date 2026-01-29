package otp

import (
	"fmt"
	"testing"
	"time"

	"github.com/pquerna/otp"
)

func TestGetTime(t *testing.T) {
	now, remaining := getTime(30)
	if remaining <= 0 || remaining > 30*time.Second {
		t.Errorf("Remaining time should be between 0 and 30 seconds, got %v", remaining)
	}

	fmt.Printf("Current time: %v, Remaining seconds: %v\n", now, remaining)
}

func TestGetOtpCode(t *testing.T) {
	// Example secret key (base32 encoded)
	secret := "JBSWY3DPEHPK3PXP"
	key, err := otp.NewKeyFromURL("otpauth://totp/Example:aa?secret=" + secret + "&issuer=TestIssuer&period=30")
	if err != nil {
		t.Fatalf("Failed to generate OTP key: %v", err)
	}
	p := NewOtpProvider(key)
	code := p.GetCode()
	if code.Error != nil {
		t.Errorf("Failed to generate OTP code: %v", code.Error)
	}
	fmt.Printf("Generated OTP Code: %s\n", code.Code)
}
