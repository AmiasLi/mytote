package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type APP struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type MYSQL struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type YamlConfig struct {
	App   APP   `yaml:"app"`
	Mysql MYSQL `yaml:"mysql"`
}

// 解析yaml配置文件
func main() {
	viperConfig := viper.New()
	// 设置配置文件名，没有后缀
	viperConfig.SetConfigName("app")
	// 设置读取文件格式为: yaml
	viperConfig.SetConfigType("yaml")
	// 设置配置文件目录(可以设置多个,优先级根据添加顺序来)
	viperConfig.AddConfigPath(".")
	// 读取解析
	if err := viperConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("配置文件未找到！%v\n", err)
			return
		} else {
			fmt.Printf("找到配置文件,但是解析错误,%v\n", err)
			return
		}
	}
	// 映射到结构体
	var yamlConfig YamlConfig
	if err := viperConfig.Unmarshal(&yamlConfig); err != nil {
		fmt.Printf("配置映射错误,%v\n", err)
	}
	fmt.Printf("App: %+v\n", yamlConfig.App)
	fmt.Printf("Mysql: %+v\n", yamlConfig.Mysql)
}
