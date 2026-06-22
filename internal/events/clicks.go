package events

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
)

func TriggerLock() error {
	cmd := exec.Command("gpgconf", "--kill", "gpg-agent")
	cmd.Run()
	return fmt.Errorf("GPG key locked")
}

func TriggerGPGPopup(recipient string) error {
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
