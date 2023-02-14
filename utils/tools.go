package utils

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"syscall"
)

func CheckXtrabackupInstalled() (string, error) {
	xtrExecPath := "/usr/local/xtrabackup/bin/xtrabackup"
	xtrExec := "xtrabackup"
	_, err := exec.LookPath(xtrExec)

	if err == nil {
		return xtrExec, nil
	} else {
		_, err = exec.LookPath(xtrExecPath)
		if err != nil {
			return "", err
		} else {
			return xtrExecPath, nil
		}
	}
}

func GetDiskFreeSpace(Path string) (int64, error) {
	if runtime.GOOS == "linux" {
		fs := syscall.Statfs_t{}
		err := syscall.Statfs(Path, &fs)
		if err != nil {
			return 0, err
		}

		// calculate the free space in gigabytes
		freeSpace := fs.Bfree * uint64(fs.Bsize)
		return int64(freeSpace), nil
	} else {
		return 0, errors.New("not support OS")
	}
}

func GetHostName() string {
	HostName, _ := os.Hostname()
	return HostName
}

func CheckDirExists(dirName string) bool {
	_, err := os.Stat(dirName)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func GetDirectorySize(directory string) (int64, error) {
	var size int64
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

func MatchBackupFile(FileString string) (bool, error) {
	match, err := regexp.MatchString(`^_\d{8}_$`,
		FileString[len(FileString)-16:len(FileString)-6])
	if err != nil {
		return false, err
	}
	return match, nil
}
