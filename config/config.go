package config

import (
	"github.com/AmiasLi/mytote/db"
	"github.com/AmiasLi/mytote/server"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"log"
)

var (
	CfgFile string
	Conf    Config
)

type Config struct {
	server.BpServer `mapstructure:"server"`
	LogMySQL        `mapstructure:"mysql_log"`
	LogMongoDB      `mapstructure:"mongodb_log"`
	LogDingTalk     `mapstructure:"ding_talk_log"`
}

type LogMySQL struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
	Table    string `yaml:"table"`
}

type LogMongoDB struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Db         string `yaml:"db"`
	Collection string `yaml:"table"`
}

type LogDingTalk struct {
	Token    string `yaml:"token"`
	ProxyUrl string `yaml:"proxy_url"`
	Secret   string `yaml:"secret"`
}

func defaultConfig() {
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

func readConfig() {
	if CfgFile != "" {
		viper.SetConfigFile(CfgFile)
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

func InitConfig() {
	readConfig()
	defaultConfig()
}
