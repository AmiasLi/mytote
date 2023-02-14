package config

import "github.com/spf13/viper"

type config struct {
	Host             string
	User             string
	Password         string
	BackupHour       int
	BackupMin        int
	Port             int
	Socket           string
	BackupDir        string
	BackupLog        string
	BackupRetain     string
	MysqlLogUser     string
	MysqlLogPassword string
	MysqlLogPort     int
	MysqlLogHost     string
	MysqlLogDb       string
	MysqlLogTable    string
	BackupType       string
	ReserveSpace     int64
	RetryDuration    int
	Compress         bool
	CompressThreads  int
}

var Conf *config

func GetConfig() {
	Conf = &config{
		User:             viper.Get("server_backup.user").(string),
		Password:         viper.Get("server_backup.password").(string),
		BackupHour:       viper.Get("server_backup.backup_hour").(int),
		BackupMin:        viper.Get("server_backup.backup_min").(int),
		Port:             viper.Get("server_backup.port").(int),
		Socket:           viper.Get("server_backup.socket").(string),
		BackupDir:        viper.Get("server_backup.backup_dir").(string),
		BackupLog:        viper.Get("server_backup.backup_log").(string),
		BackupRetain:     viper.Get("server_backup.backup_retain").(string),
		MysqlLogUser:     viper.Get("server_backup.mysql_log_user").(string),
		MysqlLogPassword: viper.Get("server_backup.mysql_log_password").(string),
		MysqlLogPort:     viper.Get("server_backup.mysql_log_port").(int),
		MysqlLogHost:     viper.Get("server_backup.mysql_log_host").(string),
		MysqlLogDb:       viper.Get("server_backup.mysql_log_db").(string),
		MysqlLogTable:    viper.Get("server_backup.mysql_log_table").(string),
		BackupType:       viper.Get("server_backup.backup_type").(string),
		ReserveSpace:     viper.Get("server_backup.reserve_space").(int64) * 1024 * 1024 * 1024,
		RetryDuration:    viper.Get("server_backup.retry_duration").(int),
		Host:             viper.Get("server_backup.host").(string),
		Compress:         viper.Get("server_backup.compress").(bool),
		CompressThreads:  viper.Get("server_backup.compress_threads").(int),
	}
}
