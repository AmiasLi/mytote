package db

import (
	"database/sql"
	"fmt"
	"github.com/AmiasLi/mytote/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type ConnString struct {
	Host     string
	Port     int
	User     string
	Password string
	Db       string
}

var stringLogServer = ConnString{
	Host:     config.Conf.MysqlLogHost,
	Port:     config.Conf.MysqlLogPort,
	User:     config.Conf.MysqlLogUser,
	Password: config.Conf.MysqlLogPassword,
	Db:       config.Conf.MysqlLogDb,
}

var stringBackupServer = ConnString{
	Host:     config.Conf.Host,
	Port:     config.Conf.Port,
	User:     config.Conf.User,
	Password: config.Conf.Password,
	Db:       "information_schema",
}

func GetConnection(conn *ConnString) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		conn.User,
		conn.Password,
		conn.Host,
		conn.Port,
		conn.Db,
	))

	if err != nil {
		logrus.Errorf("Error connecting to log database: %s\n", err)
		return db, err
	}
	return db, nil
}

func GetLogConnection() (*sql.DB, error) {
	return GetConnection(&stringLogServer)
}

func GetBackupConnection() (*sql.DB, error) {
	return GetConnection(&stringBackupServer)
}
