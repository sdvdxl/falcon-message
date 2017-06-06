package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sdvdxl/dinghook"
	"github.com/sdvdxl/falcon-message/config"
	"github.com/sdvdxl/falcon-message/sender"
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

func main() {
	cfg := config.Read()
	dingQueue := dinghook.NewQueue(cfg.DingTalk.Token, "告警", 1, 0)

	engine := echo.New()
	engine.Use(middleware.Recover())
	// engine.Use(middleware.Logger())

	engine.Server.Addr = cfg.Addr
	server := &graceful.Server{Timeout: time.Second * 10, Server: engine.Server, Logger: graceful.DefaultLogger()}
	api := engine.Group("/api")
	api.POST("/v1", func(c echo.Context) error {
		tos := c.FormValue("tos")
		var persons []sender.Person
		personsText := strings.Split(tos, ",")
		for _, v := range personsText {
			if err := json.Unmarshal([]byte(tos), &persons); err != nil {
				log.Println("parse person info error:", err, "info:", v)
			}
		} 

		content := c.FormValue("content")
		log.Println("sss", tos, content)
		if tos == "" || content == "" {
			msg := "tos or content is empty"
			log.Println(msg)
			return c.JSON(http.StatusBadRequest, msg)
		}
		dingQueue.PushWithTitle(tos, content)
		return nil
	})

	go func() {
		dingQueue.Start()
	}()

	log.Println("listening on ", cfg.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
