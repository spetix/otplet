package main

import (
	"log"
	"os"

	"fmt"

	"github.com/spetix/bar-out-adapters/pkg/barout"
	otp "github.com/spetix/otplet/internal/provider"
	"github.com/spetix/otplet/internal/render"
	"github.com/spetix/otplet/internal/secretmanager"
	"github.com/spf13/cobra"
)

func main() {
	var renderOptions render.RenderOptions

	var proto string
	var otpStore string
	var url string
	var path string
	var recipient string

	var rootCmd = &cobra.Command{
		Use:   "otplet",
		Short: "manage your otps in a smart way",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Use --help to see available commands.")
		},
	}

	create := &cobra.Command{
		Use:   "create",
		Short: "create otp setup",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Default().Print("Import generator token")
			sm := secretmanager.NewSecretManager(path, recipient)
			err := sm.SaveGenerator(url)
			if err != nil {
				log.Default().Printf("Error saving OTP generator: %v\n", err)
				os.Exit(1)
			}

		},
	}

	show := &cobra.Command{
		Use:   "show",
		Short: "show current otp",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			b := barout.New(proto)
			sm := secretmanager.NewSecretManager(otpStore, recipient)
			key, err := sm.LoadTokenFromFile()
			if err != nil {
				fmt.Printf("Error loading OTP key: %v\n", err)
				os.Exit(1)
			}
			otpProvider := otp.NewOtpProvider(key)
			if otpProvider == nil {
				os.Exit(1)
			}
			data := render.OtpDataNew(otpProvider, &renderOptions)

			b.Print(data)
		},
	}

	rootCmd.AddCommand(create)

	createFlags := create.Flags()
	createFlags.StringVarP(&url, "url", "u", "", "qr code decoded url")
	createFlags.StringVarP(&path, "path", "p", "~/.secrets", "path to save otp")
	createFlags.StringVarP(&recipient, "recipient", "r", "", "GPG recipient to encrypt OTP store")
	rootCmd.AddCommand(show)

	showFlags := show.Flags()
	showFlags.StringVarP(&proto, "proto", "p", "raw", "proto")
	lbl := os.Getenv("label")
	if lbl == "" {
		lbl = "🎄"
	}
	showFlags.StringVarP(&renderOptions.Label, "label", "l", lbl, "label")
	showFlags.StringVarP(&renderOptions.Format, "format", "f", "text", "format")
	clr := os.Getenv("color")
	if clr == "" {
		clr = "#ff0000"
	}
	showFlags.StringVarP(&renderOptions.ForegroundColor, "color", "c", clr, "foreground color")
	bgclr := os.Getenv("background")
	if bgclr == "" {
		bgclr = "#000000"
	}
	showFlags.StringVarP(&renderOptions.BackgroundColor, "background", "b", "#000000", "background color")

	showFlags.StringVarP(&otpStore, "store", "s", "otp.json", "path to otp store")
	showFlags.StringVarP(&recipient, "recipient", "r", "", "GPG recipient to decrypt OTP store")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
