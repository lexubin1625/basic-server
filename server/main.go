package main

import (
	"basic-server/core/config"
	"basic-server/core/db"
	"basic-server/core/log"
	"basic-server/routers"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 初始化配置
	var conf config.Config
	conf.InitConf()

	// 数据库初始化
	db.New()

	// 日志初始化
	log.New()

	// 路由初始化
	routers.InitRouter(router)

	router.Run(fmt.Sprintf(":%d", conf.Server.HttpPort))
}
