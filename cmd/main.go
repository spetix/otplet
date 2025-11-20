package main

import (
	"os"

	"github.com/spetix/bar-out-adapters/pkg/barout"
	"github.com/spf13/cobra"
	"fmt"
)

// type renderOptions RenderOptions {
// 	Label string
// 	Format string
// 	Color string
// 	BackgroundColor string
// }

func main() {
	var renderOptions struct {
		Label string
		Format string
		ForegroundColor string
		BackgroundColor string
	}

	var proto string

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
			// setup := token.New(name,secret)
			// setup.store()
			fmt.Println("Not yet implemented")
		},
	}
			
	show := &cobra.Command{
		Use:   "show",
		Short: "show current otp",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			b := barout.New("proto")
			b.Print(nil)
		},
	}

	rootCmd.AddCommand(show)
	rootCmd.AddCommand(create)

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

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
