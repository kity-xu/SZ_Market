package routes

import (
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	// 不需要用户认证
	rg := engine.Group("/api/business")

	// 注册启动页广告
	RegCompany(rg)
}
