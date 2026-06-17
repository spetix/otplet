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

func decryptIfNeeded(ciphertext []byte) ([]byte, error) {
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

func encryptData(plaintext []byte, recipients []string) ([]byte, error) {
	if len(recipients) == 0 {
		return nil, fmt.Errorf("missing GPG recipient")
	}

	// Build gpg command with recipients
	args := []string{"--encrypt", "--armor", "--quiet", "--batch", "--no-tty"}
	for _, recipient := range recipients {
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
func LoadTokenFromFile(location string) (provider.OtpProviderItf, error) {
	expandedPath, err := expandPath(location)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return nil, err
	}

	plaintext, err := decryptIfNeeded(data)
	if err != nil {
		// unlock required sending dummy provider
		slog.Warn("Warning: failed to decrypt OTP store, returning dummy provider", "error", err)
		return provider.NewDummyProvider(), nil
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
		slog.Warn("Error loading OTP key", "error", err)
		return provider.NewDummyProvider(), nil
	}
	otpProvider := provider.NewOtpProvider(key)

	return otpProvider, nil
}

func SaveGenerator(url string, location string, recipients []string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	slog.Info("Importing OTP generator", "location", location, "recipients", recipients)
	slog.Debug("Seed",
		"url", url,
	)
	key, err := otplib.NewKeyFromURL(url)
	if err != nil {
		return err
	}
	expandedPath, err := expandPath(location)
	if err != nil {
		return err
	}

	if len(recipients) > 0 {
		payload, err := json.Marshal(key.URL())
		if err != nil {
			return err
		}
		ciphertext, err := encryptData(payload, recipients)
		if err != nil {
			return err
		}
		if err := os.WriteFile(expandedPath, ciphertext, 0600); err != nil {
			return err
		}
		slog.Info("Encrypted OTP generator saved",
			"location", location,
		)
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
	slog.Info("OTP generator saved",
		"location", location,
	)

	return nil

}
