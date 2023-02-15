package cmd

import (
	"log"

	"github.com/AmiasLi/mytote/config"
	"github.com/AmiasLi/mytote/db"
	"github.com/AmiasLi/mytote/logs"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	Conf    config.Config
)

func init() {
	// cobra.OnInitialize(initConfig)
	// default configration
	cobra.OnInitialize(func() {
		initConfig()
		configDefault()
		Conf.BpServer.LogTable = Conf.LogMySQL.Table
		logs.InitLog(Conf.BpServer.BackupLog)
		InitMySQL()
	})

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "default ./config.yml")

}

func configDefault() {
	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", 3306)
	viper.SetDefault("server.socket", "mysql.sock")
	viper.SetDefault("server.user", "root")
	viper.SetDefault("server.backup_type", "full")
	viper.SetDefault("server.backup_retain", "7")
	viper.SetDefault("server.backup_dir", "./backup")
	viper.SetDefault("server.backup_log", "./backup/backup.log")
	viper.SetDefault("server.backup_hour", 0)
	viper.SetDefault("server.backup_minute", 0)
	viper.SetDefault("server.reserve_space", 5)
	viper.SetDefault("server.retry_duration", 15)
	viper.SetDefault("server.retry_times", 1)
	viper.SetDefault("server.compress", true)
	viper.SetDefault("server.compress_threads", 1)
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

	if err := viper.ReadInConfig(); err == nil {
		err := viper.Unmarshal(&Conf, func(conf *mapstructure.DecoderConfig) {
			conf.ErrorUnused = true
		})
		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
		}

	}
}

func InitMySQL() {
	db.ConnLogMySQL = db.ConnString(Conf.LogMySQL)
	db.ConnBackupMySQL = db.ConnString{
		Host:     Conf.BpServer.Host,
		Port:     Conf.BpServer.Port,
		User:     Conf.BpServer.User,
		Password: Conf.BpServer.Password,
		Db:       "information_schema",
	}
}
