package otp

import (
	"fmt"
	"testing"
	"time"
)

func TestOtpSetup(t *testing.T) {
	timeNow := time.Now().UTC()

	remaining := timeNow.Second() % 30
	validUntil := timeNow.Add(time.Duration(30-remaining) * time.Second)
	fmt.Print(timeNow, remaining, validUntil)
}
