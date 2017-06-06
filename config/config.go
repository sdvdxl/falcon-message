package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config 配置
type Config struct {
	Addr     string   `json:"addr"`
	DingTalk DingTalk `json:"dingTalk"`
}

// DingTalk 钉钉配置
type DingTalk struct {
	Enable bool `json:"enable"`
	// Level 等级， 只发送level 及其以下的消息
	Level uint `json:"level"`
}

// Read 读取配置
func Read() Config {
	bytes, err := ioutil.ReadFile("cfg.json")
	if err != nil {
		log.Fatalln("need file: cfg.json")
	}
	var cfg Config
	if err = json.Unmarshal(bytes, &cfg); err != nil {
		log.Fatalln("config file error", err.Error())
	}

	return cfg
}
