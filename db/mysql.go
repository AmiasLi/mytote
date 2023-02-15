package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type ConnString struct {
	Host     string
	Port     int
	User     string
	Password string
	Db       string
	Table    string
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

func GetLogConnection(connLogMySQL ConnString) (*sql.DB, error) {
	return GetConnection(&connLogMySQL)
}

func GetBackupConnection(connBackupServer ConnString) (*sql.DB, error) {
	return GetConnection(&connBackupServer)
}
