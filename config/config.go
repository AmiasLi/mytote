package config

import (
	"github.com/AmiasLi/mytote/db"
	"github.com/AmiasLi/mytote/utils"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"log"
)

var (
	CfgFile string
	Conf    Config
)

type Config struct {
	Server      `mapstructure:"server"`
	LogMySQL    `mapstructure:"mysql_log"`
	LogMongoDB  `mapstructure:"mongodb_log"`
	LogDingTalk `mapstructure:"ding_talk_log"`
}

type Server struct {
	BusinessName      string    `mapstructure:"business_name"`
	HostName          string    `mapstructure:"host_name"`
	Host              string    `mapstructure:"host"`
	Port              int       `mapstructure:"port"`
	Socket            string    `mapstructure:"socket"`
	User              string    `mapstructure:"user"`
	Password          string    `mapstructure:"password"`
	BackupMethod      string    `mapstructure:"backup_method"`
	BackupType        string    `mapstructure:"backup_type"`
	Compress          bool      `mapstructure:"compress"`
	CompressThreads   int       `mapstructure:"compress_threads"`
	BackupFullWeekday int       `mapstructure:"backup_full_weekday"`
	BackupHour        int       `mapstructure:"backup_hour"`
	BackupMin         int       `mapstructure:"backup_minute"`
	RetryDuration     int       `mapstructure:"retry_duration"`
	RetryTimes        int       `mapstructure:"retry_times"`
	BackupRetain      string    `mapstructure:"backup_retain"`
	BackupDir         string    `mapstructure:"backup_dir"`
	BackupStatus      bool      `mapstructure:"backup_status"`
	BackupSize        int64     `mapstructure:"backup_size"`
	SubDataPath       string    `mapstructure:"sub_data_path"`
	BackupLog         string    `mapstructure:"backup_log"`
	ReserveSpace      int64     `mapstructure:"reserve_space"`
	StartTime         time.Time `mapstructure:"start_time"`
	EndTime           time.Time `mapstructure:"end_time"`
	XtrBin            string    `mapstructure:"xtrabackup_bin"`
	LogTable          string    `mapstructure:"log_table"`
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
	Token    string `mapstructure:"token"`
	ProxyUrl string `mapstructure:"proxy_url"`
	Secret   string `mapstructure:"secret"`
}

func defaultConfig() {
	viper.SetDefault("server.host_name", utils.GetHostName())
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
		Host:     Conf.Server.Host,
		Port:     Conf.Server.Port,
		User:     Conf.Server.User,
		Password: Conf.Server.Password,
		Db:       "information_schema",
	}
}

func InitConfig() {
	defaultConfig()
	readConfig()
}
