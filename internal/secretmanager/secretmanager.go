// Package secretmanager handles local persistence of OTP generator URLs.
package secretmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	otplib "github.com/pquerna/otp"
	"github.com/spetix/otplet/internal/provider"
)

type SecretManagerItf interface {
	LoadTokenFromFile() (provider.OtpProviderItf, error)
	SaveGenerator(url string) error
}

type DummySecretManager struct{}

func (d *DummySecretManager) LoadTokenFromFile() (provider.OtpProviderItf, error) {
	return &provider.DummyOtpProvider{}, nil
}
func (d *DummySecretManager) SaveGenerator(url string) error {
	slog.Debug("DummySecretManager: SaveGenerator called with URL: %s", url)
	return nil
}

// SecretManager manages a file path used to save and load OTP key URLs, with optional GPG encryption.

// SecretManager manages a file path used to save and load OTP key URLs.
type SecretManager struct {
	location   string
	recipients []string
}

func NewSecretManager(location string, recipients ...string) *SecretManager {
	return &SecretManager{location: location, recipients: recipients}
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = strings.Replace(path, "~", home, 1)
	}
	return filepath.Abs(path)
}

func looksLikePGPMessage(data []byte) bool {
	return bytes.Contains(data, []byte("-----BEGIN PGP MESSAGE-----"))
}

func (sm *SecretManager) decryptIfNeeded(ciphertext []byte) ([]byte, error) {
	cmd := exec.Command("gpg", "--decrypt", "--quiet", "--batch", "--pinentry-mode", "error")
	cmd.Stdin = bytes.NewReader(ciphertext)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("gpg decrypt failed: %w (stderr: %s)", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

func (sm *SecretManager) encryptData(plaintext []byte) ([]byte, error) {
	if len(sm.recipients) == 0 {
		return nil, fmt.Errorf("missing GPG recipient")
	}

	// Build gpg command with recipients
	args := []string{"--encrypt", "--armor", "--quiet", "--batch", "--no-tty"}
	for _, recipient := range sm.recipients {
		args = append(args, "--recipient", recipient)
	}
	args = append(args, "--trust-model", "always")

	cmd := exec.Command("gpg", args...)
	cmd.Stdin = bytes.NewReader(plaintext)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("gpg encrypt failed: %w (stderr: %s)", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

// LoadTokenFromFile reads the stored OTP generator URL from disk and returns the parsed key.
func (sm *SecretManager) LoadTokenFromFile() (provider.OtpProviderItf, error) {
	expandedPath, err := expandPath(sm.location)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return nil, err
	}

	plaintext, err := sm.decryptIfNeeded(data)
	if err != nil {
		// unlock required sending dummy provider
		slog.Warn("Warning: failed to decrypt OTP store, returning dummy provider: %v", err)
		return &provider.DummyOtpProvider{}, nil
	}

	var keyUrl string
	if err := json.NewDecoder(bytes.NewReader(plaintext)).Decode(&keyUrl); err != nil {
		return nil, err
	}

	key, err := otplib.NewKeyFromURL(keyUrl)
	if err != nil {
		return nil, err
	}

	if err != nil {
		slog.Warn("Error loading OTP key: %v", err)
		return &provider.DummyOtpProvider{}, nil
	}
	otpProvider := provider.NewOtpProvider(key)

	return otpProvider, nil
}

func (sm *SecretManager) SaveGenerator(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	slog.Info("Importing OTP generator with URL: %s to path: %s", url, sm.location)
	key, err := otplib.NewKeyFromURL(url)
	if err != nil {
		return err
	}
	return sm.saveTokenToFile(key)
}

func (sm *SecretManager) saveTokenToFile(key *otplib.Key) error {
	expandedPath, err := expandPath(sm.location)
	if err != nil {
		return err
	}

	if len(sm.recipients) > 0 {
		payload, err := json.Marshal(key.URL())
		if err != nil {
			return err
		}
		ciphertext, err := sm.encryptData(payload)
		if err != nil {
			return err
		}
		if err := os.WriteFile(expandedPath, ciphertext, 0600); err != nil {
			return err
		}
		slog.Info("Encrypted OTP generator saved to %s", sm.location)
		return nil
	}

	hdlr, err := os.OpenFile(expandedPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer hdlr.Close()

	enc := json.NewEncoder(hdlr)
	if err := enc.Encode(key.URL()); err != nil {
		return err
	}
	slog.Info("OTP generator saved to %s", sm.location)

	return nil
}
