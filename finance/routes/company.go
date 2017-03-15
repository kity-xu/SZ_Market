package routes

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/controllers/company"
)

func RegCompany(rg *gin.RouterGroup) {
	// 获取广告列表
	rg.GET("/info", company.NewCompanyInfo().GetInfo)
}
