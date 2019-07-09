package app

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	FTP ftp `yaml:"ftp"`
}

type ftp struct {
	Host            string `yaml:"host"`
	Password        string `yaml:"password"`
	Login           string `yaml:"login"`
	SourcePath      string `yaml:"sourcePath"`
	DestinationPath string `yaml:"destinationPath"`
	IsRewrite       bool   `yaml:"isRewrite"`
	LogPath         string `yaml:"logPath"`
	Connections     int    `yaml:"connections"`
}

func NewConfig() *Config {
	data, err := ioutil.ReadFile("../config.yml")
	if err != nil {
		panic(err)
	}
	cnf := Config{}
	errs := yaml.Unmarshal([]byte(data), &cnf)
	if errs != nil {
		panic(errs)
	}
	return &cnf
}
