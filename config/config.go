package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Server     `mapstructure:"server"`
	LogMySQL   `mapstructure:"mysql_log"`
	LogMongoDB `mapstructure:mongodb_log`
}

type Server struct {
	Host              string `mapstructure:"host"`
	Port              int    `mapstructure:"port"`
	Socket            string `mapstructure:"socket"`
	User              string `mapstructure:"user"`
	Password          string `mapstructure:"password"`
	BakcupMethod      string `mapstructure:"backup_method"`
	BackupType        string `mapstructure:"backup_type"`
	Compress          bool   `mapstructure:"compress"`
	CompressThreads   int    `mapstructure:"compress_threads"`
	BackupFullWeekday int    `mapstructure:"backup_full_weekday"`
	BackupHour        int    `mapstructure:"backup_hour"`
	BackupMin         int    `mapstructure:"backup_minute"`
	RetryDuration     int    `mapstructure:"retry_duration"`
	RetryTimes        int    `mapstructure:"retry_times"`
	BackupRetain      string `mapstructure:"backup_retain"`
	BackupDir         string `mapstructure:"backup_dir"`
	BackupLog         string `mapstructure:"backup_log"`
	ReserveSpace      int64  `mapstructure:"reserve_space"`
}

type LogMySQL struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
	Table    string `yaml:"table"`
}

type LogMongoDB struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Db         string `yaml:"db"`
	Collection string `yaml:"table"`
}

var Conf Config
var logMySQLServer *LogMySQLServer

func main() {
	viper := viper.New()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if viper.ReadInConfig() != nil {
		logrus.Error("config file not found")
	} else {
		logrus.Info("Using config file:", viper.ConfigFileUsed())
		//fmt.Println(viper.Get("backup_server"))
		err := viper.Unmarshal(&logMySQLServer)
		fmt.Println(logMySQLServer)

		if err != nil {
			logrus.Error("unable to decode into struct, %v", err)
		}
	}
}
