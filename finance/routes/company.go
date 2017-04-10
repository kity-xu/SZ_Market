package routes

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/controllers/company"
)

func RegCompany(rg *gin.RouterGroup) {

	// 公司简介
	rg.GET("/company/info", company.NewCompany().GetInfo)

	//公司高管信息
	rg.GET("/company/manager/info", company.NewCompany().GetManagreInfo)

	// 融资分红
	rg.GET("/company/dividend", company.NewDividendInfo().GetDiv)
	rg.GET("/company/refinance", company.NewDividendInfo().GetSEO)
	rg.GET("/company/ration", company.NewDividendInfo().GetRO)

	/////////////////////////股东权益

	//股东人数
	rg.GET("/company/equity/shareholder", company.NewEquityInfo().GetShareholderJson)

	//十大流通股东
	rg.GET("/company/equity/top10", company.NewEquityInfo().GetTop10Json)

	//机构持股
	rg.GET("/company/equity/organization", company.NewEquityInfo().GetOrganizationJson)

	//股本结构
	rg.GET("/company/equity/structure", company.NewCapitalizationInfo().GetStructureJson)

	//股本变动
	rg.GET("/company/equity/changes", company.NewCapitalizationInfo().GetChangesJson)

	// 财务分析
	// 关键指标表
	rg.GET("/company/indicators", company.NewIndicators().GET)
	// 利润表
	rg.GET("/company/profits", company.NewProfits().GET)
	// 资产负债表
	rg.GET("/company/liabilities", company.NewLiabilities().GET)
	// 现金流量表
	rg.GET("/company/cashflow", company.NewCashflow().GET)
}
