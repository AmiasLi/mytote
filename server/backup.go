package server

import (
	"os/exec"
	"strconv"
	"time"

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

		s.EndTime = time.Now()
		s.BackupSize = 0
		s.BackupStatus = false
		return err

	} else {
		logrus.WithFields(logrus.Fields{
			//"output": string(out),
			"status": "SUCCESS",
		}).Info("Backup complete")

		backupSize, err := utils.GetDirectorySize(s.SubDataPath)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Error getting backup size")
			s.BackupSize = 0
		}

		s.BackupSize = backupSize
		s.BackupStatus = true
		s.EndTime = time.Now()

		return nil
	}
}
