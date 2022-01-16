package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/ocurr/cryorio/roborio"
	"github.com/spf13/cobra"
)

func NewBackupCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "backup [fileName]",
		Short: "Creates a copy of FRCUserProgram as a backup.",
		Long: "Creates a copy of FRCUserProgram as a backup.\n" +
			"The backup will be named fileName.backup if it is provided otherwise it will be named timestamp.backup.\n" +
			"The timestamp will be in the format %M_%D_%Y_%H_%M_%S",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				fmt.Println("Incorrect number of arguments.")
				fmt.Println(cmd.Use)
				return
			}
			rio, err := roborio.NewRoborio(roborio.DialSsh, "admin", "")
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

			backupName := defaultBackupName()
			if len(args) == 1 {
				backupName = args[0]
			}

			err = rio.BackupFile("FRCUserProgram", backupName)
			if err != nil {
				if errors.Is(err, roborio.ErrorFileNotExist) {
					fmt.Printf("*** FRCUserProgram does not exist ***\nDid you deploy the code?\n")
				} else {
					fmt.Printf("\n*** Unable to backup FRCUserProgram ***\n%s\n", err)
				}
				return
			}
			fmt.Println("Successfully created backup!")
		},
	}
}

func defaultBackupName() string {
	return time.Now().Format("01_02_06_03_04_05")
}
