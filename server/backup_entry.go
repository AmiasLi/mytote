package server

import (
	"errors"
	"fmt"
	"github.com/AmiasLi/mytote/config"
	"path/filepath"
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

func (s *BpServer) ManageBackup() {
	err := s.ServerBackupProcess()
	if err != nil {
		logrus.Errorf("Error running backup: %s\n", err)
	} else {
		logrus.Infof("Backup completed successfully.\n")
		s.RemoveFiles()
		err = errors.New("")
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
		BusinessName: config.Conf.Server.BusinessName,
		Instance:     s.HostName + ":" + strconv.Itoa(s.Port),
		StartTime:    s.StartTime,
		CostTime:     fmt.Sprintf("%v", s.EndTime.Sub(s.StartTime)),
		EndTime:      s.EndTime,
		FileName:     filepath.Base(s.SubDataPath),
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

func (s *BpServer) ManualBackup() {
	err := s.ServerBackupProcess()
	if err != nil {
		logrus.Errorf("Error running backup: %s\n", err)
	} else {
		logrus.Infof("Backup completed successfully.\n")
	}
}
