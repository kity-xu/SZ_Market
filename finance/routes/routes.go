package routes

import (
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	// F10
	rg := engine.Group("/api/finance/company")

	// 注册公司信息获取路径
	RegCompany(rg)
}
