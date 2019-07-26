package main

import (
    "bytes"
    "errors"
    "log"
    "net/http"
    "path"
    "strings"
    "text/template"
    "time"

    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "github.com/sdvdxl/falcon-message/config"
    "github.com/sdvdxl/falcon-message/sender"
    "github.com/sdvdxl/falcon-message/util"
    "github.com/tylerb/graceful"
)

// OK
// P3
// Endpoint:SMARTMATRIX_MONITOR
// Metric:api
// Tags:act=webapi.bindDevice,loc=GZ
// all(#2): 1==0
// Note:act=webapi.bindDevice,loc=GZ
// Max:3, Current:1
// Timestamp:2017-06-02 08:02:00
// http://127.0.0.1:8081/portal/template/view/37

const (
    // IMDingPrefix 钉钉 前缀
    IMDingPrefix = "[ding]:"
    // IMWexinPrefix 微信前缀
    // IMWexinPrefix = "[wexin]:"
)

var (
    cfg  config.Config
    ding *sender.DingTalk
    wx   *sender.Weixin

    tpl *template.Template
)

func main() {

    cfg = config.Read()
    funcMap := template.FuncMap{"elapse": func(count, reportInterval, triggerCount, postpone int) int {
        // 都使用 秒计算
        // 超过1次，要计算推迟时间
        if count > 1 {
            return reportInterval*triggerCount + postpone*(count-1)
        }

        // 第一次，直接返回上报间隔
        return reportInterval * triggerCount
    }, "divide": func(a, b int) int { return a / b },
        "timeFormat": func(t time.Time, format string) string {
            return t.Format(format)
        },
        "timeDiffFormat": func(t time.Time, format string, seconds int) string {
            return t.Add(-(time.Second * time.Duration(seconds))).Format(format)
        }}

    tpl = template.Must(template.New(path.Base(cfg.DingTalk.TemplateFile)).Funcs(funcMap).ParseFiles(cfg.DingTalk.TemplateFile))
    if cfg.DingTalk.Enable {
        ding = sender.NewDingTalk()
    }

    if cfg.Weixin.Enable {
        wx = sender.NewWeixin(cfg.Weixin.CorpID, cfg.Weixin.Secret)
        go wx.GetAccessToken()
    }

    engine := echo.New()
    engine.Server.Addr = cfg.Addr
    server := &graceful.Server{Timeout: time.Second * 10, Server: engine.Server, Logger: graceful.DefaultLogger()}
    engine.Use(middleware.Recover())
    // engine.Use(middleware.Logger())
    api := engine.Group("/api/v1")
    api.GET("/wechat/auth", wxAuth)
    api.POST("/message", func(c echo.Context) error {
        log.Println("message comming")
        tos := c.FormValue("tos")
        content := c.FormValue("content")
        log.Println("tos:", tos, " content:", content)
        if content == "" {
            return echo.NewHTTPError(http.StatusBadRequest, "content is requied")
        }

        msg, err := util.HandleContent(content)
        if err != nil {
            return err
        }

        var buffer bytes.Buffer
        if err := tpl.Execute(&buffer, msg); err != nil {
            return err
        }
        content = buffer.String()

        if strings.HasPrefix(tos, IMDingPrefix) { //是钉钉
            tokens := tos[len(IMDingPrefix):]

            if cfg.DingTalk.Enable {
                for _, v := range strings.Split(tokens, ";") {
                    go func(token string) {
                        if err := ding.Send(token, content, cfg.DingTalk.MessageType); err != nil {
                            log.Println("ERR:", err)
                        }
                    }(v)
                }
            }
        } else { //微信
            if cfg.Weixin.Enable {
                if err := wx.Send(tos, content); err != nil {
                    return echo.NewHTTPError(500, err.Error())
                }
            }
        }

        return nil
    })

    log.Println("listening on ", cfg.Addr)
    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}

// WxAuth 开启回调模式验证
func wxAuth(context echo.Context) error {
    if cfg.Weixin.Enable {
        echostr := context.FormValue("echostr")
        if echostr == "" {
            return errors.New("无法获取请求参数, echostr 为空")
        }
        var buf []byte
        var err error
        if buf, err = wx.Auth(echostr); err != nil {
            return err
        }

        return context.JSONBlob(200, buf)
    }

    return context.String(200, "微信没有启用")
}
