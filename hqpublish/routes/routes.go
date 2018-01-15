package routes

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/controllers/publish"
)

func Register(engine *gin.Engine) {
	engine.GET("/", publish.NewTest().Test)

	// 行情
	rg := engine.Group("/api/hq")

	// publish
	RegPublish(rg)

	RegPublish2(rg)

	// 财务相关
	RegFinance(engine)
}
