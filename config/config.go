package config

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

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
	IP               string
	BackupType       string
	ReserveSpace     int64
	RetryDuration    int
	Compress         bool
	CompressThreads  int
}

var Conf *config

func readConfig() {
	// Read the config file

	cfg, err := ini.Load("config.ini")

	if err != nil {
		logrus.Fatalf("Error opening config file: %s\n", err)
	}

	// Read the values from the INI file
	user := cfg.Section("server_backup").Key("user").String()
	password := cfg.Section("server_backup").Key("password").String()
	backupRetain := cfg.Section("server_backup").Key("backup_retain").MustString("2")
	port := cfg.Section("server_backup").Key("port").MustInt(3306)
	backupLog := cfg.Section("server_backup").Key("backup_log").MustString("./backup.log")
	backupHour := cfg.Section("server_backup").Key("backup_hour").MustInt(0)
	backupMin := cfg.Section("server_backup").Key("backup_minute").MustInt(0)
	socket := cfg.Section("server_backup").Key("socket").String()
	backupDir := cfg.Section("server_backup").Key("backup_dir").MustString("./")
	mysqlLogUser := cfg.Section("mysql_log").Key("mysql_log_user").String()
	mysqlLogPassword := cfg.Section("mysql_log").Key("mysql_log_password").String()
	mysqlLogPort := cfg.Section("mysql_log").Key("mysql_log_port").MustInt(3306)
	mysqlLogHost := cfg.Section("mysql_log").Key("mysql_log_host").String()
	mysqlLogDb := cfg.Section("mysql_log").Key("mysql_log_db").String()
	mysqlLogTable := cfg.Section("mysql_log").Key("mysql_log_table").String()
	ip := cfg.Section("server_backup").Key("host").MustString("127.0.0.1")
	backupType := cfg.Section("server_backup").Key("backup_type").MustString("full")
	reserveSpace := cfg.Section("server_backup").Key("reserve_space").MustInt64(5)
	retryDuration := cfg.Section("server_backup").Key("retry_duration").MustInt(15)
	host := cfg.Section("server_backup").Key("host").MustString("127.0.0.1")
	compress := cfg.Section("server_backup").Key("compress").MustBool(true)
	compressThreads := cfg.Section("server_backup").Key("compress_threads").MustInt(1)

	if user == "" || password == "" {
		logrus.Fatalf("Error reading user or password from config file")
	}

	Conf = &config{
		User:             user,
		Password:         password,
		BackupHour:       backupHour,
		BackupMin:        backupMin,
		Port:             port,
		Socket:           socket,
		BackupDir:        backupDir,
		BackupLog:        backupLog,
		BackupRetain:     backupRetain,
		MysqlLogUser:     mysqlLogUser,
		MysqlLogPassword: mysqlLogPassword,
		MysqlLogPort:     mysqlLogPort,
		MysqlLogHost:     mysqlLogHost,
		MysqlLogDb:       mysqlLogDb,
		MysqlLogTable:    mysqlLogTable,
		IP:               ip,
		BackupType:       backupType,
		ReserveSpace:     reserveSpace * 1024 * 1024 * 1024,
		RetryDuration:    retryDuration,
		Host:             host,
		Compress:         compress,
		CompressThreads:  compressThreads,
	}
}

func init() {
	readConfig()
}
