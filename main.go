package main
import (
	"basic-server/core/config"
	"basic-server/routers"
	"fmt"
	"github.com/gin-gonic/gin"
)
func main(){
	router := gin.Default()

	// 初始化配置
	var conf config.Config
	conf.InitConf()

	routers.InitRouter(router)

	router.Run(fmt.Sprintf(":%d",conf.Server.HttpPort))
}
