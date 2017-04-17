package routes

import (
	"github.com/gin-gonic/gin"
	"haina.com/share/lib"
)

func RegPublish(rg *gin.RouterGroup) {

	// test
	rg.GET("/ping", pong)
}

func pong(c *gin.Context) {
	lib.WriteString(c, 200, "pong")
}
