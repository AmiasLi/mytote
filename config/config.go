package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	BackupServer   BackupServer   `yaml:"backup_server"`
	LogMySQLServer LogMySQLServer `yaml:"mysql_log"`
}

type BackupServer struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Socket          string `yaml:"socket"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	BackupType      string `yaml:"backup_type"`
	Compress        bool   `yaml:"compress"`
	CompressThreads int    `yaml:"compress_threads"`
	BackupHour      int    `yaml:"backup_hour"`
	BackupMin       int    `yaml:"backup_minute"`
	RetryDuration   int    `yaml:"retry_duration"`
	BackupRetain    string `yaml:"backup_retain"`
	BackupDir       string `yaml:"backup_dir"`
	BackupLog       string `yaml:"backup_log"`
	ReserveSpace    int64  `yaml:"reserve_space"`
}

type LogMySQLServer struct {
	MysqlLogHost     string `yaml:"host"`
	MysqlLogPort     string `yaml:"port"`
	MysqlLogUser     string `yaml:"user"`
	MysqlLogPassword string `yaml:"password"`
	MysqlLogDb       string `yaml:"db"`
	MysqlLogTable    string `yaml:"table"`
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
