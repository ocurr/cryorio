package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ocurr/cryorio/roborio"
	"github.com/spf13/cobra"
)

func NewRestoreCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "restore fileName",
		Short: "Restores a backup file as FRCUserProgram. The user will still need to restart the robot for the program to take effect.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Incorrect number of arguments.")
				fmt.Println(cmd.Use)
				return
			}
			rio, err := roborio.NewRoborio(roborio.DialSsh, frcUser, frcPassword)
			if err != nil {
				fmt.Println(err)
				return
			}
			for err := rio.Connect(); err != nil; err = rio.Connect() {
				if errors.Is(err, roborio.ErrorNoConnection) {
					fmt.Printf("*** Unable to connect to the roborio ***\nIs the robot on?\n")
					return
				}
				fmt.Printf("Connection attempt failed: %s\n", err)
			}
			defer rio.Disconnect()

			backupName := strings.TrimSuffix(args[0], ".backup")

			err = rio.BackupFile(backupName, frcUserProgram)
			if err != nil {
				if errors.Is(err, roborio.ErrorFileNotExist) {
					fmt.Printf("*** %s does not exist ***\nDid you type the name correctly\n", backupName)
				} else {
					fmt.Printf("\n*** Unable to restore %s ***\n%s\n", backupName, err)
				}
				return
			}
			fmt.Println("Successfully created backup!")
		},
	}
}
