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
	Mail     Mail     `json:"mail"`
}

// DingTalk 钉钉配置
type DingTalk struct {
	Enable bool `json:"enable"`
	// Level 等级， 只发送level 及其以下的消息
	Level  uint        `json:"level"`
	Groups []DingGroup `json:"groups"`
}

// DingGroup 钉钉配置
type DingGroup struct {
	Key   string `json:"key"` // 关键词，要和 note 里面匹配，如果没有配置，则全部下发
	Token string `json:"token"`
}

// Mail 邮件配置
type Mail struct {
	Manager  string `json:"manager"` // 管理人员的邮件，如果发送失败会发送消息给管理人员
	SMTP     string `json:"smtp"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
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

	checkParams(cfg)
	return cfg
}

func checkParams(cfg Config) {
	for _, v := range cfg.DingTalk.Groups {
		if v.Token == "" {
			log.Fatal("dingTalk group key:", v.Key, " need token")
		}
	}
}
