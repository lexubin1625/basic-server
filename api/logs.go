package api

import (
	"basic-server/core/response"
	"basic-server/dao"
	"fmt"
	"github.com/gin-gonic/gin"
)


// GetAllLogs is a function to get a slice of record(s) from logs table in the task database
// @Summary Get list of Logs
// @Tags Logs
// @Description GetAllLogs is a handler to get a slice of record(s) from logs table in the task database
// @Accept  json
// @Produce  json
// @Param   id     query    int     false        "page requested (defaults to 0)"
// @Success 200
// @Failure 400
// @Failure 404
// @Router /api/v1/logs [get]
// http "http://localhost:8080/api/v1/logs?id=6198" X-Api-User:user123

func GetLogOne(c *gin.Context){
	id := c.GetInt("id")
	fmt.Println(id)

	logs,_ := dao.GetLogsFirst(fmt.Sprintf("id = %d",id),[]string{"id desc"})

	response.JsonResponse(c,response.SUCCESS,logs)
}