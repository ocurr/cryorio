package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ocurr/cryorio/roborio"
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists the *.backup files available",
		Run: func(cmd *cobra.Command, args []string) {
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

			out, err := rio.ListDir()
			if err != nil {
				fmt.Printf("*** Unable to get contents of directory ***\n%s\n", err)
			}
			outStr := strings.TrimSpace(string(out))
			files := strings.Split(outStr, "\n")
			for _, f := range files {
				if strings.HasSuffix(f, ".backup") {
					fmt.Println(f)
				}
			}
		},
	}
}
