package main

import (
	"EditSync/middleware"
	"EditSync/router"
	"EditSync/server/task"
	"EditSync/server/websocket"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"os"
)

func main() {
	initConfig()

	initLogoFile()

	initRedis()

	router.WebSocketInit()

	go websocket.StartWebSocket()
	// 定时任务
	task.Init()
	// 服务注册
	task.ServerInit()

	r := router.InitRouter()
	base := fmt.Sprintf("%s:%d", viper.GetString("app.httpIp"), viper.GetInt("app.httpPort"))
	fmt.Println("start server @", base)
	err := r.Run(base)
	if err != nil {
		fmt.Println("failed to start: ", err.Error())
	}

}

func initLogoFile() {
	gin.DisableConsoleColor()
	logFile := viper.GetString("app.logFile")
	fp, _ := os.Create(logFile)
	gin.DefaultWriter = io.MultiWriter(fp)
}

func initConfig() {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file:%s \n", err))
	}
	fmt.Println("config app", viper.Get("app"))
	fmt.Println("config socket", viper.Get("socket"))
	fmt.Println("config redis", viper.Get("redis"))
}

func initRedis() {
	middleware.NewClient()
}
