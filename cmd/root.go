package cmd

import (
	"fmt"
	"github.com/AmiasLi/mytote/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Conf    config.Config
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "mytote",
		Short: "mytote is backup Cron",
		Long: "mytote is backup Cron for your MySQL " +
			"using xtrabackup, " +
			"and manage backup file.",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	configDefault()
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file(default ./config.yaml)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(backupCmd)

}

func configDefault() {
	viper.SetDefault("backup_server.host", "127.0.0.1")
	viper.SetDefault("backup_server.port", 3306)
	viper.SetDefault("backup_server.socket", "mysql.sock")
	viper.SetDefault("backup_server.user", "root")
	viper.SetDefault("backup_server.backup_type", "full")
	viper.SetDefault("backup_server.backup_retain", "7")
	viper.SetDefault("backup_server.backup_dir", "./backup")
	viper.SetDefault("backup_server.backup_log", "./backup/backup.log")
	viper.SetDefault("backup_server.backup_hour", 0)
	viper.SetDefault("backup_server.backup_minute", 0)
	viper.SetDefault("backup_server.reserve_space", 5)
	viper.SetDefault("backup_server.retry_duration", 15)
	viper.SetDefault("backup_server.compress", true)
	viper.SetDefault("backup_server.compress_threads", 1)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}
	viper.AutomaticEnv()

	if viper.ReadInConfig() != nil {
		logrus.Error("config file not found")
	} else {
		logrus.Info("Using config file:", viper.ConfigFileUsed())
		fmt.Println(viper.Get("backup_server"))
		err := viper.Unmarshal(&Conf)
		if err != nil {
			logrus.Error("unable to decode into struct, %v", err)
		}
	}
	fmt.Println(Conf)
	fmt.Println(viper.Get("backup_server"))
}
