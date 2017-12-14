package routes

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/controllers/publish/f10"
)

func RegF10(rg *gin.RouterGroup) {
	// F10首页
	rg.POST("/f10/home", f10.NewHN_F10_Mobile().GetF10_Mobile)

	// 公司详细信息
	rg.POST("/f10/comInfo", f10.NewCompany().GetF10_ComInfo)

	// 历史股本变动
	rg.POST("/f10/capstock/change", f10.NewCapitalStock().GetF10_CapitalStock)

	// 十大股东+十大流通股东
	rg.POST("/f10/holder/top10", f10.NewShareholderslTop10().GetShareholdersTop10)
}
