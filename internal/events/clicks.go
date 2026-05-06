package events

import (
	"crypto/rand"
	"log"
	"os"
	"os/exec"

	"github.com/spetix/otplet/internal/secretmanager"
)

type ClickManager struct{}

// HandleClickEvents handles mouse click events from i3blocks
// BLOCK_BUTTON environment variable: 1=left, 3=right
func (cm *ClickManager) HandleClickEvents(eventId string, otpStore string, recipient string) (sm secretmanager.SecretManagerItf) {
	switch eventId {
	case "1":
		// Left click - unlock GPG key
		if err := triggerGPGPopup(recipient); err != nil {
			log.Printf("Warning: failed to unlock GPG key: %v\n", err)
			return &secretmanager.DummySecretManager{}
		}
	case "3":
		// Right click - lock GPG key
		cmd := exec.Command("gpgconf", "--kill", "gpg-agent")
		if err := cmd.Run(); err != nil {
			log.Printf("Warning: failed to lock GPG key: %v\n", err)
		}
	}
	return secretmanager.NewSecretManager(otpStore, recipient)
}

func triggerGPGPopup(recipient string) error {
	// create temp files
	plain, err := os.CreateTemp("", "gpg-plain-*")
	if err != nil {
		return err
	}
	defer os.Remove(plain.Name())

	enc, err := os.CreateTemp("", "gpg-enc-*")
	if err != nil {
		return err
	}
	defer os.Remove(enc.Name())

	// write a few random bytes
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return err
	}
	if _, err := plain.Write(buf); err != nil {
		return err
	}
	plain.Close()

	// encrypt (no passphrase needed)
	encrypt := exec.Command(
		"gpg",
		"--yes",
		"--batch",
		"-r", recipient,
		"--encrypt",
		"-o", enc.Name(),
		plain.Name(),
	)
	if err := encrypt.Run(); err != nil {
		return err
	}

	// decrypt → this triggers pinentry popup
	decrypt := exec.Command(
		"gpg",
		"--decrypt",
		enc.Name(),
	)
	decrypt.Stdout = nil
	decrypt.Stderr = nil

	return decrypt.Run()
}
