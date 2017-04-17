package routes

import (
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	// 行情
	rg := engine.Group("/api/hq")

	// publish
	RegPublish(rg)
}
