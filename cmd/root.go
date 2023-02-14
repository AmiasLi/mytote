package cmd

import (
	"github.com/AmiasLi/mytote/config"
	"github.com/AmiasLi/mytote/logs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
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
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(config.GetConfig)
	cobra.OnInitialize(logs.InitLog)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file(default ./config.yaml)")

	configDefault()

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(backupCmd)

}

func configDefault() {
	viper.SetDefault("server_backup.backup_retain", "7")
	viper.SetDefault("server_backup.backup_dir", "./backup")
	viper.SetDefault("server_backup.backup_log", "./backup/backup.log")
	viper.SetDefault("server_backup.backup_hour", 0)
	viper.SetDefault("server_backup.backup_minute", 0)
	viper.SetDefault("server_backup.port", 3306)
	viper.SetDefault("server_backup.host", "127.0.0.1")
	viper.SetDefault("server_backup.backup_type", "full")
	viper.SetDefault("server_backup.reserve_space", 5)
	viper.SetDefault("server_backup.retry_duration", 15)
	viper.SetDefault("server_backup.compress", true)
	viper.SetDefault("server_backup.compress_threads", 1)
	viper.SetDefault("mysql_log.mysql_log_port", 3306)
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
}
