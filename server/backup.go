package server

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/AmiasLi/mytote/logs"
	"github.com/AmiasLi/mytote/utils"
	"github.com/sirupsen/logrus"
)

func (s *BpServer) BackupCommand() []string {
	var cmdXtra []string
	cmdXtra = append(cmdXtra,
		"--user="+s.User,
		"--host="+s.Host,
		"--password="+s.Password,
		"--port="+strconv.Itoa(s.Port),
		"--socket="+s.Socket,
		"--backup",
		"--target-dir="+s.SubDataPath)
	if s.Compress {
		cmdXtra = append(cmdXtra, "--compress",
			"--compress-threads="+strconv.Itoa(s.CompressThreads))
	}
	return cmdXtra
}

func (s *BpServer) Backup() error {
	// Run the xtrabackup command
	logrus.Info("Running xtrabackup...")

	out, err := exec.Command(s.XtrBin, s.BackupCommand()...).CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"output": string(out),
			"status": "ERROR",
		}).Error("Error running xtrabackup")

		LogContentsObj := logs.NewLogContents(err.Error(), "ERROR",
			s.StartTime, time.Now(),
			fmt.Sprintf("%v", time.Now().Sub(s.StartTime)),
			s.SubDataPath, 0, s.Host, s.Port, s.BackupType)

		logs.LogToMySQL(LogContentsObj, s.LogTable)

		return err

	} else {
		logrus.WithFields(logrus.Fields{
			//"output": string(out),
			"status": "SUCCESS",
		}).Info("Backup complete")

		s.BackupStatus = true

		s.EndTime = time.Now()
		backupSize, err := utils.GetDirectorySize(s.SubDataPath)
		s.BackupSize = backupSize

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Error getting backup size")
			backupSize = 0
		}
		LogContentsObj := logs.NewLogContents("", "SUCCESS",
			s.StartTime, time.Now(),
			fmt.Sprintf("%v", s.EndTime.Sub(s.StartTime)),
			s.SubDataPath, s.BackupSize, s.Host, s.Port, s.BackupType)

		logs.LogToMySQL(LogContentsObj, s.LogTable)

		return nil
	}
}
