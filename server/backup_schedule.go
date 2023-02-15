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
	s.HostName = utils.GetHostName()

	s.StartTime = time.Now()
	s.EndTime = time.Now()

	// Check if the server is online
	online, err := s.GetServerStatus()
	if err != nil {
		logrus.Errorf("can not connect to backup server[%s:%d] : %s", s.Host, s.Port, err)
		LogContentsObj := logs.NewLogContents(err.Error(), "ERROR",
			s.StartTime, time.Now(),
			fmt.Sprintf("%v", time.Now().Sub(s.StartTime)),
			"", 0, s.Host, s.Port, s.BackupType)

		logs.LogToMySQL(LogContentsObj, s.LogTable)
		return err
	}

	xtrExec, err := utils.CheckXtrabackupInstalled()
	if err != nil {
		logrus.Errorf("xtrabackup not found: %s\n", err)

		LogContentsObj := logs.NewLogContents(err.Error(), "ERROR",
			s.StartTime, time.Now(),
			fmt.Sprintf("%v", time.Now().Sub(s.StartTime)),
			"", 0, s.Host, s.Port, s.BackupType)

		logs.LogToMySQL(LogContentsObj, s.LogTable)
		return err
	}

	spaceAllow, err := s.SpaceAllow()

	if err != nil {
		logrus.Errorf("error checking free disk: %s\n", err)

		LogContentsObj := logs.NewLogContents(err.Error(), "ERROR",
			s.StartTime, time.Now(),
			fmt.Sprintf("%v", time.Now().Sub(s.StartTime)),
			"", 0, s.Host, s.Port, s.BackupType)

		logs.LogToMySQL(LogContentsObj, s.LogTable)
		return err
	}

	if TarGetDirectory, err := s.GenSubPath(); err == nil {
		s.SubDataPath = TarGetDirectory
	} else {
		logrus.Errorf("Error creating backup directory: %s\n", err)

		LogContentsObj := logs.NewLogContents(err.Error(), "ERROR",
			s.StartTime, time.Now(),
			fmt.Sprintf("%v", time.Now().Sub(s.StartTime)),
			"", 0, s.Host, s.Port, s.BackupType)

		logs.LogToMySQL(LogContentsObj, s.LogTable)

		return err
	}

	if online && xtrExec != "" && s.SubDataPath != "" && spaceAllow {
		// Run the backup
		s.XtrBin = xtrExec
		err := s.Backup()
		if err != nil {
			logrus.Errorf("Error running backup: %s\n", err)
			logrus.Infof("Retrying afer %d.\n", s.RetryDuration)
			time.Sleep(time.Duration(s.RetryDuration) * time.Minute)

			err := s.Backup()
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
