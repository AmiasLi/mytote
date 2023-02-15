package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mytote",
		Short: "mytote is backup Cron",
		Long: "mytote is backup Cron for your MySQL " +
			"using xtrabackup, " +
			"and manage backup file.",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
