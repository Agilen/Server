package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/tkanos/gonfig"
)

type Config struct {
	DatabaseURl  string
	BaseUrl      string
	Port         string
	SmtpHost     string
	SmtpPort     int
	MailSender   string
	MailPassword string
	TokenTTL     int
}

func DefualtConfig() *Config {
	return &Config{
		Port:         ":10000",
		DatabaseURl:  "host=127.0.0.1 user=avtor password=12QWaszx dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Kiev",
		TokenTTL:     24,
		SmtpHost:     "smtp.gmail.com",
		SmtpPort:     587,
		MailSender:   "mail",
		MailPassword: "password",
		BaseUrl:      "http://localhost",
	}
}

func NewConfig(configPath string) *Config {
	var config *Config
	var err error
	if IsConfigExist(configPath) {
		config, err = GetConfigFromFile(configPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		config = DefualtConfig()
		err := SaveConfig(config, configPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	return config
}

func GetConfigFromFile(path string) (*Config, error) {
	conf := Config{}
	err := gonfig.GetConf(path, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func SaveConfig(conf *Config, path string) error {
	a, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, a, os.ModePerm)
}

func IsConfigExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false

		}
	}
	return true
}
