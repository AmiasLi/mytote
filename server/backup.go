package server

import (
	"errors"
	"os/exec"
	"strconv"
	"time"

	"github.com/AmiasLi/mytote/utils"
	"github.com/sirupsen/logrus"
)

func (s *BpServer) BackupCommand() []string {
	var cmdExtra []string
	cmdExtra = append(cmdExtra,
		"--user="+s.User,
		"--host="+s.Host,
		"--password="+s.Password,
		"--port="+strconv.Itoa(s.Port),
		"--socket="+s.Socket,
		"--backup",
		"--target-dir="+s.SubDataPath)
	if s.Compress {
		cmdExtra = append(cmdExtra, "--compress",
			"--compress-threads="+strconv.Itoa(s.CompressThreads))
	}
	return cmdExtra
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
