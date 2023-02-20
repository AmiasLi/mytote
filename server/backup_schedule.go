package server

import (
	"errors"
	"fmt"
	"github.com/AmiasLi/mytote/config"
	"strconv"
	"time"

	"github.com/AmiasLi/mytote/logs"
	"github.com/AmiasLi/mytote/utils"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

func (s *BpServer) BackupFailConditionStage() {
	s.BackupStatus = false
	s.EndTime = time.Now()
	s.SubDataPath = ""
	s.BackupSize = 0
}

func (s *BpServer) ServerBackupProcess() error {
	// Create a new backup server
	s.HostName = utils.GetHostName()

	s.StartTime = time.Now()
	s.EndTime = time.Now()

	// Check if the server is online
	online, err := s.GetServerStatus()
	if err != nil {
		logrus.Errorf("can not connect to backup server[%s:%d] : %s", s.Host, s.Port, err)
		s.BackupFailConditionStage()
		return err
	}

	xtrExec, err := utils.CheckXtrabackupInstalled()
	if err != nil {
		logrus.Errorf("xtrabackup not found: %s\n", err)
		s.BackupFailConditionStage()
		return err
	}

	spaceAllow, err := s.SpaceAllow()

	if err != nil {
		logrus.Errorf("error checking free disk: %s\n", err)
		s.BackupFailConditionStage()
		return err
	}

	if TarGetDirectory, err := s.GenSubPath(); err == nil {
		s.SubDataPath = TarGetDirectory
	} else {
		logrus.Errorf("Error creating backup directory: %s\n", err)
		s.BackupFailConditionStage()
		return err
	}

	if online && xtrExec != "" && s.SubDataPath != "" && spaceAllow {
		// Run the backup
		s.XtrBin = xtrExec
		err := s.Backup()
		if err != nil {
			logrus.Errorf("Error running backup: %s\n", err)
			var errRoll error

			for i := 0; i < s.RetryTimes; i++ {
				logrus.Warnf("Retrying backup after %d minutes.\n", s.RetryDuration)
				time.Sleep(time.Duration(s.RetryDuration) * time.Minute)
				logrus.Infof("Retrying %d afer %d.\n", i, s.RetryDuration)

				errRoll = s.Backup()
				if errRoll != nil {
					logrus.Errorf("Error running backup again: %s\n", errRoll)
					continue

				} else {
					break
				}
			}

			if errRoll != nil {
				return errors.New("backup failed after " +
					"retrying " + strconv.Itoa(s.RetryTimes) + " times" + errRoll.Error())
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	return nil
}

func (s *BpServer) ManualBackup() {
	err := s.ServerBackupProcess()
	if err != nil {
		logrus.Errorf("Error running backup: %s\n", err)
	} else {
		logrus.Infof("Backup completed successfully.\n")
		s.RemoveFiles()
	}
}

func (s *BpServer) ManageBackup() {
	err := s.ServerBackupProcess()
	if err != nil {
		logrus.Errorf("Error running backup: %s\n", err)
	} else {
		logrus.Infof("Backup completed successfully.\n")
		s.RemoveFiles()
	}

	LogContentsObj := logs.NewLogContents(err.Error(), s.BackupStatus,
		s.StartTime, s.EndTime,
		fmt.Sprintf("%v", s.EndTime.Sub(s.StartTime)),
		s.SubDataPath, s.BackupSize, s.Host, s.Port, s.BackupType)

	logs.LogToMySQL(LogContentsObj, s.LogTable)

	FormatFileSize := utils.ByteHumanRead(s.BackupSize)
	logDing := logs.LogContentDingTalk{
		Token:        config.Conf.LogDingTalk.Token,
		ProxyUrl:     config.Conf.LogDingTalk.ProxyUrl,
		Secret:       config.Conf.LogDingTalk.Secret,
		BusinessName: config.Conf.BpServer.BusinessName,
		Instance:     s.HostName + "-" + s.Host + ":" + strconv.Itoa(s.Port),
		StartTime:    s.StartTime,
		CostTime:     fmt.Sprintf("%v", s.EndTime.Sub(s.StartTime)),
		EndTime:      s.EndTime,
		FileName:     s.SubDataPath,
		FileSize:     FormatFileSize,
		Status:       s.BackupStatus,
	}

	err = logDing.ResultToDingTalkGroup()

	if err != nil {
		logrus.Errorf("Error sendding dingtalk message: %s\n", err)
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
