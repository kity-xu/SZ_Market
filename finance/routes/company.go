package routes

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/controllers/company"
)

func RegCompany(rg *gin.RouterGroup) {
	// 获取广告列表
	rg.GET("/info", company.NewCompanyInfo().GetInfo)

	// 获取关键指标
	rg.GET("/indicators", company.NewIndicators().GET)
	// 获取利润表
	rg.GET("/profits", company.NewProfits().GET)
	// 获取资产负债表
	rg.GET("/liabilities", company.NewLiabilities().GET)
	// 获取现金流量表
	rg.GET("/cashflow", company.NewCashflow().GET)
}
