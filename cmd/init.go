package cmd

import (
	"github.com/AmiasLi/mytote/config"
	"github.com/AmiasLi/mytote/logs"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(func() {
		config.InitConfig()
		logs.InitLog(config.Conf.Server.BackupLog)
		config.InitMySQL()
	})

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.PersistentFlags().StringVarP(&config.CfgFile,
		"config", "c", "", "default ./config.yml")

}
