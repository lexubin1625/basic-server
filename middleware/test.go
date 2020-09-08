package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
)

func Test()gin.HandlerFunc{
	return func(c *gin.Context) {
		log.Println("MiddlewareA before request")
		c.Set("key1", 123)
		c.Next()
		log.Println("MiddlewareA after request")
	}
}

func Test2()gin.HandlerFunc{
	return func(c *gin.Context) {
		log.Println("Middleware Test2 before request")
		c.Set("key1", 123)
		log.Println("Middleware Test2  after request")
	}
}