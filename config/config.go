package config

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "time"
)

// Config 配置
type Config struct {
    Addr     string   `json:"addr"`
    DingTalk DingTalk `json:"dingTalk"`

    Weixin Weixin `json:"weixin"`
}

// Weixin 微信配置
type Weixin struct {
    Enable         bool
    CorpID         string `json:"corpID"`
    AgentID        string `json:"agentId"`
    Secret         string `json:"secret"`
    EncodingAESKey string `json:"encodingAESKey"`
}

type AlarmMessage struct {
    Level        string // 告警等级 P1
    Type         string // 类型 PROBLEM，OK
    Endpoint     string // 主机host或者ip
    Desc         string // 告警描述
    Counter      string // 告警指标
    Tags         string // tags
    TriggerCount int    // 间隔
    Count        int    // 当前告警次数
    Time         time.Time
    Expression   string
    // 告警时间
}

// DingTalk 钉钉配置
type DingTalk struct {
    Enable bool `json:"enable"`
    // Level 等级， 只发送level 及其以下 的消息

    Level        uint `json:"level"`
    TemplateFile string
    MessageType  string // markdown ，text
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
