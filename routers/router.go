package routers

import (
	"basic-server/api"
	"basic-server/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine){
	apiv1 := r.Group("/api/v1")
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiv1.Use(middleware.Test(),middleware.Test2())
	apiv1.GET("/logs", api.GetLogOne)
}