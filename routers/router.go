package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(r *gin.Engine){
	apiv1 := r.Group("/api/v1")
	apiv1.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	//var conf config.Config
	//config.Viper.AllKeys()
	//fmt.Println(config.Viper.AllKeys())
}