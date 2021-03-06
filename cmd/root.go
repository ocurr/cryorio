package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.1"
const frcUserProgram = "FRCUserProgram"
const frcUser = "lvuser"
const frcPassword = ""

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "cryorio",
		Short: "cryorio is a utility to enable backing up executables on the FRC roboRIO",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("cryorio version v%s\n", version)
		},
	}

	rootCmd.AddCommand(
		NewBackupCommand(),
		NewRestoreCommand(),
		NewListCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
