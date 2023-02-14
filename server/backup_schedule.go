package server

import (
	"fmt"
	"github.com/AmiasLi/mytote/config"
	"github.com/AmiasLi/mytote/logs"
	"github.com/AmiasLi/mytote/utils"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func BackupEntry() error {
	// Create a new backup server
	HostName := utils.GetHostName()
	Port := config.Conf.Port
	bpServer := NewBpServer(HostName, Port)
	bpServer.StartTime = time.Now()
	bpServer.EndTime = time.Now()

	// Check if the server is online
	online, err := bpServer.GetServerStatus()
	if err != nil {
		logrus.Errorf("Server is offline: %s", err)
		LogContentsObj := logs.NewLogContents(err.Error(), "ERROR",
			bpServer.StartTime, time.Now(),
			fmt.Sprintf("%v", time.Now().Sub(bpServer.StartTime)),
			"", 0)

		logs.LogToMySQL(LogContentsObj)
		return err
	}

	xtrExec, err := utils.CheckXtrabackupInstalled()
	if err != nil {
		logrus.Errorf("xtrabackup not found: %s\n", err)

		LogContentsObj := logs.NewLogContents(err.Error(), "ERROR",
			bpServer.StartTime, time.Now(),
			fmt.Sprintf("%v", time.Now().Sub(bpServer.StartTime)),
			"", 0)

		logs.LogToMySQL(LogContentsObj)
		return err
	}

	spaceAllow, err := bpServer.SpaceAllow()
	if err != nil {
		logrus.Errorf("Error checking disk space: %s\n", err)

		LogContentsObj := logs.NewLogContents(err.Error(), "ERROR",
			bpServer.StartTime, time.Now(),
			fmt.Sprintf("%v", time.Now().Sub(bpServer.StartTime)),
			"", 0)

		logs.LogToMySQL(LogContentsObj)
		return err
	}

	if TarGetDirectory, err := bpServer.GenSubPath(); err == nil {
		bpServer.SubDataPath = TarGetDirectory
	} else {
		logrus.Errorf("Error creating backup directory: %s\n", err)

		LogContentsObj := logs.NewLogContents(err.Error(), "ERROR",
			bpServer.StartTime, time.Now(),
			fmt.Sprintf("%v", time.Now().Sub(bpServer.StartTime)),
			"", 0)

		logs.LogToMySQL(LogContentsObj)

		return err
	}

	if online && xtrExec != "" && bpServer.SubDataPath != "" && spaceAllow {
		// Run the backup
		bpServer.xtrBin = xtrExec
		err := bpServer.Backup()
		if err != nil {
			logrus.Errorf("Error running backup: %s\n", err)
			logrus.Infof("Retrying afer %d.\n", config.Conf.RetryDuration)
			time.Sleep(time.Duration(config.Conf.RetryDuration) * time.Minute)

			err := bpServer.Backup()
			if err != nil {
				logrus.Errorf("Error running backup again: %s\n", err)
				return err
			} else {
				return nil
			}
		}
	}
	return nil
}

func DoBackup() {
	RemoveFiles()
	err := BackupEntry()
	if err != nil {
		logrus.Errorf("Error running backup: %s\n", err)
	} else {
		logrus.Infof("Backup completed successfully.\n")
	}
}

func CombineSchedule() string {
	cronSchedule := "0 " + strconv.Itoa(config.Conf.BackupMin) +
		" " + strconv.Itoa(config.Conf.BackupHour) + " * * *"
	return cronSchedule
}

func BackupCron() {
	logrus.Infof("mytote is running")
	cronSchedule := CombineSchedule()
	c := cron.New()
	err := c.AddFunc(cronSchedule, DoBackup)
	if err != nil {
		logrus.Fatalf("Error adding cron job: %s\n", err)
	}
	c.Start()

	select {}
}
