package config

import "github.com/AmiasLi/mytote/server"

type Config struct {
	server.BpServer `mapstructure:"server"`
	LogMySQL        `mapstructure:"mysql_log"`
	LogMongoDB      `mapstructure:"mongodb_log"`
	LogES           `mapstructure:"es_log"`
}

type LogMySQL struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
	Table    string `yaml:"table"`
}

type LogMongoDB struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Db         string `yaml:"db"`
	Collection string `yaml:"table"`
}

type LogES struct {
}
