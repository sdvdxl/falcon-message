package sender

import "time"
import "github.com/sdvdxl/dinghook"

// Sender 发送接口
type Sender interface {
	Send(tos []string, msg Message) error
}

type Message struct {
	Status     string // PROBLEM OK
	Level      string //  告警等级
	Endpoint   string
	Metric     string
	Tags       []Tag
	Expression string //表达式
	Note       string //注释
	Max        uint
	Current    uint
	Timestamp  time.Time
	URL        string
}

type Tag struct {
	Key   string
	Value interface{}
}

