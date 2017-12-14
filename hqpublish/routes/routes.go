package routes

import (
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	// 行情
	rg := engine.Group("/api/hq")

	// F10 2.2 zxw
	RegF10(rg)

	// publish
	RegPublish(rg)

	// 财务相关
	RegFinance(engine)
}
