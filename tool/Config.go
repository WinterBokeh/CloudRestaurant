package tool

import (
	"bufio"
	"encoding/json"
	"os"
)

//从app.json文件中读取参数
type Config struct {
	AppName string `json:"app_name"`
	AppHost string `json:"app_host"`
	AppMode string `json:"app_mode"`
	AppPort string `json:"app_port"`
}

var _cfg *Config = nil

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)           //new一个解析器
	if err = decoder.Decode(&_cfg); err != nil { //使用decode把json中的参数赋值给_cfg，记得加取地址符
		return nil, err
	}
	return _cfg, nil
}
