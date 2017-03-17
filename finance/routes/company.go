package routes

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/controllers/company"
)

func RegCompany(rg *gin.RouterGroup) {
	rg.GET("/info", company.NewCompanyInfo().GetInfo)

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

}
