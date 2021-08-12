package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type AppConfig struct {
	Server   *ServerConfig   `json:"server"`
	Auth     *AuthConfig     `json:"auth"`
	Database *DatabaseConfig `json:"database"`
}

type ServerConfig struct {
	Name    string `json:"name"`
	Port    int    `json:"port"`
	Version string `json:"version"`
}

type DatabaseConfig struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Schema      string `json:"schema"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	MaxIdleConn int    `json:"maxIdleConn"`
	MaxOpenConn int    `json:"maxOpenConn"`
}

type AuthConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoadFile() *AppConfig {
	path := os.Getenv("CONFIG_PATH")
	if len(path) == 0 {
		path = "./config.json"
	}

	file, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	appConfig := AppConfig{}
	err = json.Unmarshal(b, &appConfig)
	if err != nil {
		panic(err)
	}

	return &appConfig
}

func (c ServerConfig) AppAddress() string {
	return fmt.Sprintf(":%v", c.Port)
}
