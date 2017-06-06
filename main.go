package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sdvdxl/dinghook"
	"github.com/sdvdxl/falcon-message/config"
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
)

func main() {
	cfg := config.Read()

	engine := echo.New()
	engine.Use(middleware.Recover())
	// engine.Use(middleware.Logger())

	engine.Server.Addr = cfg.Addr
	server := &graceful.Server{Timeout: time.Second * 10, Server: engine.Server, Logger: graceful.DefaultLogger()}
	api := engine.Group("/api")
	api.POST("/v1", func(c echo.Context) error {
		log.Println("message comming")
		tos := c.FormValue("tos")
		if strings.HasPrefix(tos, IMDingPrefix) { //是钉钉
			token := tos[len(IMDingPrefix):]
			if token == "" {
				log.Println("ERROR: ding token is blank")
				return echo.NewHTTPError(http.StatusBadRequest, "need dingding token")
			}

			// 发送钉钉
			ding := dinghook.NewDing(token)
			content := c.FormValue("content")
			result := ding.SendMessage(dinghook.Message{Content: content})
			log.Println(result)
			if !result.Success {
				log.Println("token:", token)
				return echo.NewHTTPError(http.StatusBadRequest, result.ErrMsg)
			}
		}

		return nil
	})

	log.Println("listening on ", cfg.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
