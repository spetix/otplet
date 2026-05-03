// Package secretmanager handles local persistence of OTP generator URLs.
package secretmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	otplib "github.com/pquerna/otp"
	"github.com/proglottis/gpgme"
)

// SecretManager manages a file path used to save and load OTP key URLs.
type SecretManager struct {
	location   string
	recipients []string
}

func NewSecretManager(location string, recipients ...string) SecretManager {
	return SecretManager{location: location, recipients: recipients}
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

func newGPGMECtx() (*gpgme.Context, error) {
	ctx, err := gpgme.New()
	if err != nil {
		return nil, err
	}
	if err := ctx.SetProtocol(gpgme.ProtocolOpenPGP); err != nil {
		ctx.Release()
		return nil, err
	}
	ctx.SetArmor(true)
	return ctx, nil
}

func (sm SecretManager) decryptIfNeeded(ciphertext []byte) ([]byte, error) {
	ctx, err := newGPGMECtx()
	if err != nil {
		return nil, err
	}
	defer ctx.Release()

	cipherData, err := gpgme.NewDataBytes(ciphertext)
	if err != nil {
		return nil, err
	}
	defer cipherData.Close()

	plainData, err := gpgme.NewData()
	if err != nil {
		return nil, err
	}
	defer plainData.Close()

	if err := ctx.Decrypt(cipherData, plainData); err != nil {
		return nil, err
	}
	if _, err := plainData.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	return io.ReadAll(plainData)
}

func (sm SecretManager) encryptData(plaintext []byte) ([]byte, error) {
	if len(sm.recipients) == 0 {
		return nil, fmt.Errorf("missing GPG recipient")
	}

	ctx, err := newGPGMECtx()
	if err != nil {
		return nil, err
	}
	defer ctx.Release()

	recipients := make([]*gpgme.Key, 0, len(sm.recipients))
	for _, recipient := range sm.recipients {
		keys, err := gpgme.FindKeys(recipient, false)
		if err != nil {
			return nil, fmt.Errorf("failed to search for recipient %q: %w", recipient, err)
		}
		if len(keys) == 0 {
			return nil, fmt.Errorf("recipient %q not found", recipient)
		}
		key := keys[0]
		if !key.CanEncrypt() {
			key.Release()
			return nil, fmt.Errorf("recipient %q cannot encrypt", recipient)
		}
		recipients = append(recipients, key)
	}
	defer func() {
		for _, key := range recipients {
			key.Release()
		}
	}()

	plainData, err := gpgme.NewDataBytes(plaintext)
	if err != nil {
		return nil, err
	}
	defer plainData.Close()

	cipherData, err := gpgme.NewData()
	if err != nil {
		return nil, err
	}
	defer cipherData.Close()

	if err := ctx.Encrypt(recipients, gpgme.EncryptAlwaysTrust, plainData, cipherData); err != nil {
		return nil, err
	}
	if _, err := cipherData.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	return io.ReadAll(cipherData)
}

// LoadTokenFromFile reads the stored OTP generator URL from disk and returns the parsed key.
func (sm SecretManager) LoadTokenFromFile() (*otplib.Key, error) {
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
		if looksLikePGPMessage(data) {
			return nil, err
		}
		plaintext = data
	}

	var keyUrl string
	if err := json.NewDecoder(bytes.NewReader(plaintext)).Decode(&keyUrl); err != nil {
		return nil, err
	}

	key, err := otplib.NewKeyFromURL(keyUrl)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (sm SecretManager) SaveGenerator(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	log.Default().Printf("Importing OTP generator with URL: %s to path: %s\n", url, sm.location)
	key, err := otplib.NewKeyFromURL(url)
	if err != nil {
		return err
	}
	return sm.saveTokenToFile(key)
}

func (sm SecretManager) saveTokenToFile(key *otplib.Key) error {
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
		log.Default().Printf("Encrypted OTP generator saved to %s\n", sm.location)
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
	log.Default().Printf("OTP generator saved to %s\n", sm.location)

	return nil
}
