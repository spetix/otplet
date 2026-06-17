package provider

import (
	"testing"
	"time"

	otplib "github.com/pquerna/otp"
)

func TestGetTime(t *testing.T) {
	now, remaining := getTime(30)
	if remaining <= 0 || remaining > 30*time.Second {
		t.Fatalf("remaining time should be between 0 and 30 seconds, got %v", remaining)
	}

	if now.IsZero() {
		t.Fatalf("expected non-zero current time")
	}
}

func TestGetOtpCode(t *testing.T) {
	secret := "JBSWY3DPEHPK3PXP"
	key, err := otplib.NewKeyFromURL("otpauth://totp/Example:aa?secret=" + secret + "&issuer=TestIssuer&period=30")
	if err != nil {
		t.Fatalf("failed to generate OTP key: %v", err)
	}

	p := NewOtpProvider(key)
	if p == nil {
		t.Fatal("expected non-nil OtpProvider")
	}

	code := p.GetCode()
	if code == nil {
		t.Fatal("expected non-nil PassCode from provider")
	}
	if code.Error != nil {
		t.Fatalf("failed to generate OTP code: %v", code.Error)
	}
	if len(code.Code) != 6 {
		t.Fatalf("expected OTP code length of 6, got %d", len(code.Code))
	}
}
