package sender

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/sdvdxl/dinghook"
)

type DingTalk struct {
}

func (d *DingTalk) Send(token string, content string) error {
	if token == "" {
		return errors.New("need dingding token")
	}

	// 发送钉钉
	ding := dinghook.NewDing(token)
	result := ding.SendMessage(dinghook.Message{Content: content})
	log.Println(result)
	if !result.Success {
		log.Println("token:", token)
		return echo.NewHTTPError(http.StatusBadRequest, result.ErrMsg)
	}

	return nil
}

func NewDingTalk() *DingTalk {
	return &DingTalk{}
}
