package sender

import (
    "errors"
    "github.com/labstack/echo"
    "github.com/sdvdxl/dinghook"
    "log"
    "net/http"
)

type DingTalk struct {
}

func (d *DingTalk) Send(token string, content, msgType string) error {
    if token == "" {
        return errors.New("need dingding token")
    }

    // 发送钉钉
    ding := dinghook.NewDing(token)
    var result dinghook.Result
    if msgType == dinghook.MsgTypeMarkdown {
        result = ding.SendMarkdown(dinghook.Markdown{Title: "告警", Content: content})
    } else {
        result = ding.SendMessage(dinghook.Message{Content: content})
    }
    log.Println(result)
    if !result.Success {
        log.Println("token:", token, " send result:", result)
        return echo.NewHTTPError(http.StatusBadRequest, result.ErrMsg)
    }

    return nil
}

func NewDingTalk() *DingTalk {
    return &DingTalk{}
}
