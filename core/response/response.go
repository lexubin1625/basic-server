package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func JsonResponse(c *gin.Context,code int,data interface{}){
	msg,ok := msgMap[code]
	if !ok {
		msg = msgMap[FAIL]
	}
	c.JSON(http.StatusOK,Response{
		Code:code,
		Msg:msg,
		Data:data,
	})
}