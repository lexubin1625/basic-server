package main

import (
	"basic-server/core/config"
	"basic-server/core/db"
	"basic-server/core/log"
	"basic-server/routers"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "basic-server/docs"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"os"
	"os/signal"
	"syscall"
)

var (
	OsSignal chan os.Signal
)

func main() {
	// 初始化配置
	var conf config.Config
	_, err := conf.InitConf()
	if err != nil {
		panic(err)
	}

	// 数据库初始化
	err = db.New()
	if err != nil {
		panic(err)
	}

	// 日志初始化
	err = log.New()
	if err != nil {
		panic(err)
	}

	OsSignal = make(chan os.Signal, 1)

	go ginServer(conf)
	LoopForever()

}

// ginServer launch gin http server
func ginServer(conf config.Config) {
	router := gin.Default()

	// 是否开启swagger
	if conf.Server.SwaggerEnable == true {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	// 路由初始化
	routers.InitRouter(router)
	router.Run(fmt.Sprintf(":%d", conf.Server.HttpPort))
}

func LoopForever() {
	fmt.Printf("Entering infinite loop\n")

	signal.Notify(OsSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	_ = <-OsSignal

	fmt.Printf("Exiting infinite loop received OsSignal\n")

}