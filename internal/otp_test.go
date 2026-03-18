package otp

import (
	"fmt"
	"os"
	"testing"
	"time"

	otplib "github.com/pquerna/otp"
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
	key, err := otplib.NewKeyFromURL("otpauth://totp/Example:aa?secret=" + secret + "&issuer=TestIssuer&period=30")
	if err != nil {
		t.Fatalf("Failed to generate OTP key: %v", err)
	}
	p := NewOtpProvider(key)

	// GetCode should never be nil given our provider seeds the value
	code := p.GetCode()
	if code == nil {
		t.Fatalf("expected non-nil PassCode from provider")
	}
	if code.Error != nil {
		t.Errorf("Failed to generate OTP code: %v", code.Error)
	}
	fmt.Printf("Generated OTP Code: %s\n", code.Code)
	if len(code.Code) != 6 {
		t.Errorf("Expected OTP code length of 6, got %d", len(code.Code))
	}
}

func TestSaveAndLoadToken(t *testing.T) {
	// Example secret key (base32 encoded
	secret := "JBSWY3DPEHPK3PXP"
	key, err := otplib.NewKeyFromURL("otpauth://totp/Example:aa?secret=" + secret + "&issuer=TestIssuer&period=30")
	if err != nil {
		t.Fatalf("Failed to generate OTP key: %v", err)
	}

	tmpdir, err := os.MkdirTemp("", "example-otp")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpdir) // Clean up the temporary directory

	tmpFile := tmpdir + "/test_otp_key.json"
	err = SaveTokenToFile(tmpFile, key)
	if err != nil {
		t.Fatalf("Failed to save OTP key to file: %v", err)
	}
	defer os.Remove(tmpFile)

	// Load the key back from the file
	loadedKey, err := LoadTokenFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load OTP key from file: %v", err)
	}

	// Verify that the loaded key matches the original key
	if loadedKey.Secret() != key.Secret() {
		t.Errorf("Loaded key secret does not match original key secret")
	}
}
