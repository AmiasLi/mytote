package logs

import (
	"database/sql"
	"github.com/AmiasLi/mytote/config"
	"github.com/AmiasLi/mytote/db"
	"github.com/AmiasLi/mytote/utils"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

type LogContents struct {
	HostName     string
	IP           string
	Port         int
	BackupType   string
	BackupFile   string
	BackupSize   int64
	StartTime    time.Time
	EndTime      time.Time
	Duration     string
	BackupStatus string
	ErrMessage   string
}

func LogToMySQL(contents *LogContents) {
	// Connect to the MySQL server
	// Insert the log into the MySQL database
	DbObj, err := db.GetLogConnection()
	if err != nil {
		logrus.Errorf("Error connecting to the log database: %s\n", err)
		return
	}

	// Insert the log into the MySQL database
	stmt, err := DbObj.Prepare("insert into" + " " + config.Conf.MysqlLogTable +
		" (" + "host_name, " +
		"ip, port, backup_type, " +
		"backup_file, backup_size," +
		" start_time, end_time,duration,backup_status, err_message) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		logrus.Errorf("Error preparing the MySQL statement: %s\n", err)
		return
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			logrus.Errorf("Error closing the MySQL statement: %s\n", err)
		}
	}(stmt)

	if _, err = stmt.Exec(contents.HostName,
		contents.IP, contents.Port,
		contents.BackupType, contents.BackupFile,
		contents.BackupSize, contents.StartTime,
		contents.EndTime, contents.Duration, contents.BackupStatus,
		contents.ErrMessage); err != nil {
		logrus.Errorf("Error inserting the log into the MySQL database: %s\n", err)
	}
}

func NewLogContents(errMessage string, backupStatus string, startTime time.Time, endTime time.Time,
	duration string, backupFile string, backupSize int64) *LogContents {
	logContentsObj := LogContents{
		HostName:     utils.GetHostName(),
		IP:           config.Conf.IP,
		Port:         config.Conf.Port,
		BackupType:   config.Conf.BackupType,
		StartTime:    startTime,
		EndTime:      endTime,
		Duration:     duration,
		BackupSize:   backupSize,
		BackupFile:   backupFile,
		ErrMessage:   errMessage,
		BackupStatus: backupStatus,
	}
	return &logContentsObj
}

func InitLog() {
	// Log to the file
	//WrOutput1 := os.Stdout

	f, err := os.OpenFile(config.Conf.BackupLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("error opening log file: %v", err)
	}

	//logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(io.MultiWriter(f))
}

func init() {
	InitLog()
}
