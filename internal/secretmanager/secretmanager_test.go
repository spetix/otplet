package secretmanager

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveAndLoadToken(t *testing.T) {
	secret := "JBSWY3DPEHPK3PXP"
	url := "otpauth://totp/Example:aa?secret=" + secret + "&issuer=TestIssuer&period=30"

	tmpDir, err := os.MkdirTemp("", "example-otp")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "test_otp_key.json")
	sm := NewSecretManager(filePath)

	if err := sm.SaveGenerator(url); err != nil {
		t.Fatalf("failed to save OTP generator: %v", err)
	}

	loadedKey, err := sm.LoadTokenFromFile()
	if err != nil {
		t.Fatalf("failed to load OTP key from file: %v", err)
	}

	if loadedKey.Secret() != secret {
		t.Fatalf("loaded key secret does not match expected secret")
	}
}

func TestSaveAndLoadTokenEncrypted(t *testing.T) {
	if _, err := exec.LookPath("gpg"); err != nil {
		t.Skip("gpg not installed")
	}

	secret := "JBSWY3DPEHPK3PXP"
	url := "otpauth://totp/Example:aa?secret=" + secret + "&issuer=TestIssuer&period=30"

	tmpDir, err := os.MkdirTemp("", "example-otp")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	gpgHome, err := os.MkdirTemp("", "gpg-home")
	if err != nil {
		t.Fatalf("failed to create temporary GNUPGHOME: %v", err)
	}
	defer os.RemoveAll(gpgHome)
	if err := os.Setenv("GNUPGHOME", gpgHome); err != nil {
		t.Fatalf("failed to set GNUPGHOME: %v", err)
	}

	env := append(os.Environ(), "GNUPGHOME="+gpgHome)
	cmd := exec.Command("gpg", "--batch", "--yes", "--pinentry-mode", "loopback", "--passphrase", "", "--homedir", gpgHome, "--quick-generate-key", "testuser@example.com", "default", "default", "never")
	cmd.Env = env
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to generate GPG key: %v\n%s", err, string(output))
	}

	filePath := filepath.Join(tmpDir, "encrypted_otp_key.gpg")
	sm := NewSecretManager(filePath, "testuser@example.com")

	if err := sm.SaveGenerator(url); err != nil {
		t.Fatalf("failed to save encrypted OTP generator: %v", err)
	}

	ciphertext, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read encrypted file: %v", err)
	}
	if !strings.Contains(string(ciphertext), "-----BEGIN PGP MESSAGE-----") {
		t.Fatalf("expected encrypted file contents to include PGP message header")
	}

	loadedKey, err := sm.LoadTokenFromFile()
	if err != nil {
		t.Fatalf("failed to load encrypted OTP key from file: %v", err)
	}

	if loadedKey.Secret() != secret {
		t.Fatalf("loaded encrypted key secret does not match expected secret")
	}
}
