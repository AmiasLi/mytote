package server

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/AmiasLi/mytote/db"
	"github.com/AmiasLi/mytote/utils"
	"github.com/sirupsen/logrus"
)

type FilesInfo struct {
	ID         int
	BackupFile string
}

func (s *BpServer) RemoveFiles() {

	// Remove the backup files that are older than the retention period

	if Files, err := s.NeedRemoveFiles(); err == nil {
		if FileList, err := s.RemoveFilesList(Files); err == nil {
			if len(FileList) > 0 {
				logrus.Infof("Remove files: %s\n", FileList)
			} else {
				logrus.Infof("No files to remove\n")
			}
		} else {
			logrus.Errorf("Error removing files: %s\n", err)
		}
	} else {
		logrus.Errorf("Error getting files: %s\n", err)
	}
}

func (s *BpServer) NeedRemoveFiles() ([]FilesInfo, error) {
	connection, err := db.GetLogConnection()
	if err != nil {
		return nil, err
	}
	defer func(connection *sql.DB) {
		err := connection.Close()
		if err != nil {
			logrus.Errorf("Error closing the MySQL connection: %s\n", err)
		}
	}(connection)

	// Insert the log into the MySQL database
	rows, err := connection.Query("select id,backup_file " +
		"from " + s.LogTable + " where backup_status = 'SUCCESS' " +
		" and backup_type = 'full' and backup_file_status = 1" + " and backup_date < now()" +
		" - interval " + s.BackupRetain + " day;")

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logrus.Errorf("Error closing the MySQL statement: %s\n", err)
		}
	}(rows)

	var filesInfoSlice []FilesInfo
	for rows.Next() {
		var filesInfo FilesInfo
		if err := rows.Scan(&filesInfo.ID, &filesInfo.BackupFile); err != nil {
			return nil, err
		}
		filesInfoSlice = append(filesInfoSlice, filesInfo)
	}
	return filesInfoSlice, nil
}

func (s *BpServer) RemoveFilesList(backupFileList []FilesInfo) ([]string, error) {
	//Remove the backup files
	// update the file info in the MySQL database
	var RemoveFileList []string

	if len(backupFileList) == 0 {
		return RemoveFileList, nil
	}

	for _, backupFile := range backupFileList {
		if filepath.Dir(backupFile.BackupFile) == s.BackupDir {

			// Check if the file is a backup file, if not, do not delete it
			RemoveFlag, err := utils.MatchBackupFile(backupFile.BackupFile)
			if err != nil {
				logrus.Errorf("Error matching backup file: %s\n", err)
			} else {
				if RemoveFlag == false {
					logrus.Errorf("The file %s is not match the backup file type, "+
						"do not delete it\n", backupFile.BackupFile)
					continue
				}
			}

			err = os.RemoveAll(backupFile.BackupFile)
			if err == nil {
				RemoveFileList = append(RemoveFileList, backupFile.BackupFile)
				err = s.UpdateFileStatus(backupFile.ID)
				if err != nil {
					logrus.Errorf("Error updating file status: %s\n", err)
				}
			} else {
				logrus.Errorf("Error removing %s: %s\n", backupFile.BackupFile, err)
			}
		}
	}
	return RemoveFileList, nil
}

func (s *BpServer) UpdateFileStatus(backupFileId int) error {
	connection, err := db.GetLogConnection()
	if err != nil {
		return err
	}
	defer func(connection *sql.DB) {
		err := connection.Close()
		if err != nil {
			logrus.Errorf("Error closing the MySQL connection: %s\n", err)
		}
	}(connection)

	_, err = connection.Exec("update "+s.LogTable+
		" set backup_file_status = 0 ,file_drop_time = now()"+""+
		"where id = ?;", backupFileId)

	if err != nil {
		return err
	}

	return nil
}
