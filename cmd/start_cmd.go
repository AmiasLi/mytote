package cmd

import (
	"github.com/AmiasLi/mytote/config"
	"github.com/AmiasLi/mytote/server"
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
			Instance := TransInstance()
			Instance.ManualBackup()
		}
	},
}

func TransInstance() server.BpServer {
	config.Conf.Server.LogTable = config.Conf.LogMySQL.Table
	return server.BpServer(config.Conf.Server)
}
