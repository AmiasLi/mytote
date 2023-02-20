package cmd

import (
	"github.com/AmiasLi/mytote/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the backup immediately",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.ConfigFileUsed() == "" {
			logrus.Fatal("config file not found")
		} else {
			config.Conf.BpServer.ManageBackup()
		}
	},
}
