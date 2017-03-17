package routes

import (
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	// 不需要用户认证
	rg := engine.Group("/api/finance")

	// 注册公司信息获取路径
	RegCompany(rg)
}
