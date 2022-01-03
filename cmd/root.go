package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "cryorio",
		Short: "cryorio is a utility to enable backing up executables on the FRC roboRIO",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
