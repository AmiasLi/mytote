package server

import (
	"fmt"
	"github.com/AmiasLi/mytote/config"
	"github.com/AmiasLi/mytote/logs"
	"github.com/AmiasLi/mytote/utils"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strconv"
	"time"
)

func (s *BpServer) BackupCommand() []string {
	var cmd []string
	cmd = append(cmd,
		"--user="+config.Conf.User,
		"--host="+config.Conf.Host,
		"--password="+config.Conf.Password,
		"--port="+strconv.Itoa(config.Conf.Port),
		"--socket="+config.Conf.Socket,
		"--backup",
		"--target-dir="+s.SubDataPath)
	if config.Conf.Compress {
		cmd = append(cmd, "--compress",
			"--compress-threads="+strconv.Itoa(config.Conf.CompressThreads))
	}
	return cmd
}

func (s *BpServer) Backup() error {
	// Run the xtrabackup command
	logrus.Info("Running xtrabackup...")

	out, err := exec.Command(s.xtrBin, s.BackupCommand()...).CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"output": string(out),
			"status": "ERROR",
		}).Error("Error running xtrabackup")

		LogContentsObj := logs.NewLogContents(err.Error(), "ERROR",
			s.StartTime, time.Now(),
			fmt.Sprintf("%v", time.Now().Sub(s.StartTime)),
			s.SubDataPath, 0)

		logs.LogToMySQL(LogContentsObj)

		return err

	} else {
		logrus.WithFields(logrus.Fields{
			//"output": string(out),
			"status": "SUCCESS",
		}).Info("Backup complete")

		s.EndTime = time.Now()
		backupSize, err := utils.GetDirectorySize(s.SubDataPath)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Error getting backup size")
			backupSize = 0
		}
		LogContentsObj := logs.NewLogContents("", "SUCCESS",
			s.StartTime, time.Now(),
			fmt.Sprintf("%v", s.EndTime.Sub(s.StartTime)),
			s.SubDataPath, backupSize)

		logs.LogToMySQL(LogContentsObj)

		return nil
	}
}
