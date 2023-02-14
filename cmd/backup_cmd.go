package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Print the version number of mytote",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start backup")
	},
}
