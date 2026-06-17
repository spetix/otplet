package render

import (
	"testing"

	otplib "github.com/pquerna/otp"
	"github.com/spetix/bar-out-adapters/pkg/barout"
	"github.com/spetix/otplet/internal/provider"
	"github.com/spf13/cobra"
)

func TestOtpDataShortAndLong(t *testing.T) {
	secret := "JBSWY3DPEHPK3PXP"
	key, err := otplib.NewKeyFromURL("otpauth://totp/Example:aa?secret=" + secret + "&issuer=TestIssuer&period=30")
	if err != nil {
		t.Fatalf("failed to create OTP key: %v", err)
	}
	blkSetup := barout.NewSetupBlocklet(nil)
	c := &cobra.Command{Use: "test"}
	blkSetup.Setup(c)
	c.PersistentFlags().Set("label", "OTP")

	p := provider.NewOtpProvider(key)
	data := OtpDataNew(p, blkSetup.Options())

	short := data.Short()
	if short == "" {
		t.Fatal("expected non-empty short OTP output")
	}

	long := data.Long()
	if long == "" {
		t.Fatal("expected non-empty long OTP output")
	}
}

func TestOtpDataNoOtpConfigured(t *testing.T) {
	blkSetup := barout.NewSetupBlocklet(nil)

	data := OtpDataNew(nil, blkSetup.Options())

	if got := data.Short(); got != "No OTP configured" {
		t.Fatalf("expected no OTP configured message, got %q", got)
	}
	if got := data.Long(); got != "No OTP configured" {
		t.Fatalf("expected no OTP configured message, got %q", got)
	}
}
