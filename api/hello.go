package api

import (
	"basic-server/core/response"
	"fmt"
	"github.com/gin-gonic/gin"
)



func Hello(c *gin.Context){
	fmt.Println(c.GetInt("key1"))
	response.JsonResponse(c,response.SUCCESS,nil)
	//c.String(http.StatusOK, "Hello World")
}