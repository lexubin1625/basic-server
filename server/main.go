package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"x-server/core/config"
	"x-server/core/db"
	"x-server/core/log"
	_ "x-server/docs"
	"x-server/routers"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
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
