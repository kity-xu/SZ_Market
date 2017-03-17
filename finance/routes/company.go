package routes

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/controllers/company"
)

func RegCompany(rg *gin.RouterGroup) {
	// 获取广告列表
	rg.GET("/company/info", company.NewCompanyInfo().GetInfo)
	rg.GET("/company/dividend", company.NewDividendInfo().GetDiv)
	rg.GET("/company/refinance", company.NewDividendInfo().GetSEO)
	rg.GET("/company/ration", company.NewDividendInfo().GetRO)
}
