package server

import (
	"fmt"
	"strconv"
	"time"

	"github.com/AmiasLi/mytote/logs"
	"github.com/AmiasLi/mytote/utils"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

func (s *BpServer) ServerBackupProcess() error {
	// Create a new backup server
	HostName := utils.GetHostName()
	bpServer := NewBpServer(HostName, s.Port)
	bpServer.StartTime = time.Now()
	bpServer.EndTime = time.Now()

	// Check if the server is online
	online, err := bpServer.GetServerStatus()
	if err != nil {
		logrus.Errorf("Server is offline: %s", err)
		LogContentsObj := logs.NewLogContents(err.Error(), "ERROR",
			bpServer.StartTime, time.Now(),
			fmt.Sprintf("%v", time.Now().Sub(bpServer.StartTime)),
			"", 0, s.Host, s.Port, s.BackupType)

		logs.LogToMySQL(LogContentsObj, s.LogTable)
		return err
	}

	xtrExec, err := utils.CheckXtrabackupInstalled()
	if err != nil {
		logrus.Errorf("xtrabackup not found: %s\n", err)

		LogContentsObj := logs.NewLogContents(err.Error(), "ERROR",
			bpServer.StartTime, time.Now(),
			fmt.Sprintf("%v", time.Now().Sub(bpServer.StartTime)),
			"", 0, s.Host, s.Port, s.BackupType)

		logs.LogToMySQL(LogContentsObj, s.LogTable)
		return err
	}

	spaceAllow, err := bpServer.SpaceAllow()
	if err != nil {
		logrus.Errorf("Error checking disk space: %s\n", err)

		LogContentsObj := logs.NewLogContents(err.Error(), "ERROR",
			bpServer.StartTime, time.Now(),
			fmt.Sprintf("%v", time.Now().Sub(bpServer.StartTime)),
			"", 0, s.Host, s.Port, s.BackupType)

		logs.LogToMySQL(LogContentsObj, s.LogTable)
		return err
	}

	if TarGetDirectory, err := bpServer.GenSubPath(); err == nil {
		bpServer.SubDataPath = TarGetDirectory
	} else {
		logrus.Errorf("Error creating backup directory: %s\n", err)

		LogContentsObj := logs.NewLogContents(err.Error(), "ERROR",
			bpServer.StartTime, time.Now(),
			fmt.Sprintf("%v", time.Now().Sub(bpServer.StartTime)),
			"", 0, s.Host, s.Port, s.BackupType)

		logs.LogToMySQL(LogContentsObj, s.LogTable)

		return err
	}

	if online && xtrExec != "" && bpServer.SubDataPath != "" && spaceAllow {
		// Run the backup
		bpServer.XtrBin = xtrExec
		err := bpServer.Backup()
		if err != nil {
			logrus.Errorf("Error running backup: %s\n", err)
			logrus.Infof("Retrying afer %d.\n", bpServer.RetryDuration)
			time.Sleep(time.Duration(bpServer.RetryDuration) * time.Minute)

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

func (s *BpServer) ManageBackup() {
	s.RemoveFiles()
	err := s.ServerBackupProcess()
	if err != nil {
		logrus.Errorf("Error running backup: %s\n", err)
	} else {
		logrus.Infof("Backup completed successfully.\n")
	}
}

func (s *BpServer) BackupCron() {
	logrus.Infof("mytote is running")
	cronSchedule := "0 " + strconv.Itoa(s.BackupMin) +
		" " + strconv.Itoa(s.BackupHour) + " * * *"
	c := cron.New()
	err := c.AddFunc(cronSchedule, s.ManageBackup)
	if err != nil {
		logrus.Fatalf("Error adding cron job: %s\n", err)
	}
	c.Start()

	select {}
}
