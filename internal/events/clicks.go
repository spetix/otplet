package events

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
)

type ClickManager struct{}

// BLOCK_BUTTON environment variable: 1=left, 3=right
func (cm *ClickManager) HandleClickEvents(eventId string, otpStore string, recipient string) error {
	switch eventId {
	case "1":
		return triggerGPGPopup(recipient)
	case "3":
		// Right click - lock GPG key
		cmd := exec.Command("gpgconf", "--kill", "gpg-agent")
		cmd.Run()
		return fmt.Errorf("GPG key locked")
	}
	return nil
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
