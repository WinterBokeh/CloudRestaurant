package tool

import (
	"bufio"
	"encoding/json"
	"os"
)

type Config struct {
	AppName string     `json:"app_name"`
	AppHost string     `json:"app_host"`
	AppPort string     `json:"app_port"`
	AppMode string     `json:"app-mode"`
	Sms     SmsConfig  `json:"sms"`
	Database DatabaseConfig `json:"database"`
	Redis RedisConfig `json:"redis_config"`
}

type RedisConfig struct {
	Addr string `json:"addr"`
	Port string `json:"port"`
	Password string `json:"password"`
	Db int `json:"db"`
}

type SmsConfig struct {
	SignName	 string `json:"sign_name"`
	TemplateCode string `json:"template_code"`
	RegionId	 string `json:"region_id"`
	AppKey		 string `json:"app_key"`
	AppSecret 	 string `json:"app_secret"`
}

type DatabaseConfig struct {
	Driver   string `json:"driver"`
	User 	 string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port	 string `json:"port"`
	DbName 	 string `json:"db_name"`
	ShowSql	 bool   `json:"show_sql"`
	Charset  string `json:"charset"`
}

var _cfg *Config

func GetConfig() *Config {
	return _cfg
}

func PraseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	read := bufio.NewReader(file)
	decoder := json.NewDecoder(read)
	err = decoder.Decode(&_cfg)
	if err != nil {
		return nil, err
	}
	return _cfg, nil
}


