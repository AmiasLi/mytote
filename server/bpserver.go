package server

import (
	"database/sql"
	"errors"
	"github.com/AmiasLi/mytote/config"
	"os"
	"time"

	"github.com/AmiasLi/mytote/db"
	"github.com/AmiasLi/mytote/utils"
	"github.com/sirupsen/logrus"
)

type BpServer config.Server

func (s *BpServer) GenSubPath() (string, error) {
	// Generate the backup directory path

	// Check if the backup directory exists
	IsExists := utils.CheckDirExists(s.BackupDir)
	if !IsExists {
		return "", errors.New("backup directory does not exist")
	}

	SubDirName := s.HostName + "_backup_" + time.Now().Format("20060102_150405")
	TarGetDir := s.BackupDir + "/" + SubDirName
	err := os.Mkdir(TarGetDir, 0755)
	if err != nil {
		return "", errors.New("error creating backup directory")
	}
	return TarGetDir, nil
}

func (s *BpServer) GetServerVersion() interface{} {
	return nil
}

// GetServerStatus MySQL online status
func (s *BpServer) GetServerStatus() (bool, error) {
	// Check if the server is online
	dbs, err := db.GetBackupConnection()
	defer func(dbs *sql.DB) {
		err := dbs.Close()
		if err != nil {
			logrus.Errorf("Error closing the backup database: %s\n", err)
		}
	}(dbs)
	if err != nil {
		return false, err
	}
	_, err = dbs.Exec("SELECT 1")
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *BpServer) EstimateDatabaseSize() (int64, error) {
	dbs, err := db.GetBackupConnection()
	if err != nil {
		return 0, err
	}
	defer func(dbs *sql.DB) {
		err := dbs.Close()
		if err != nil {
			logrus.Errorf("Error closing the backup database: %s\n", err)
		}
	}(dbs)

	// Get the database size
	var SizeDataBase int64
	err = dbs.QueryRow("select sum(FILE_SIZE) from information_schema.INNODB_TABLESPACES;").Scan(&SizeDataBase)

	if err != nil {
		return 0, err
	}
	return SizeDataBase, nil
}

func (s *BpServer) SpaceAllow() (bool, error) {
	// Check if the disk space is enough
	// Get the disk space

	FreeSpace, err := utils.GetDiskFreeSpace(s.BackupDir)
	if err != nil {
		return false, err
	}

	// Get the database size
	SizeDataBase, err := s.EstimateDatabaseSize()
	if err != nil {
		logrus.Errorf("Error getting the database size: %s\n", err)
		return false, err
	}

	// Check if the disk space is enough
	if FreeSpace < SizeDataBase+s.ReserveSpace {
		return false, errors.New("the disk space is not enough")
	}
	return true, nil
}
