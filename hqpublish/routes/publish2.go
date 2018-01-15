package routes

import (
	"github.com/gin-gonic/gin"

	"haina.com/market/hqpublish/controllers/publish/f10"
	"haina.com/market/hqpublish/controllers/publish2"
)

func RegPublish2(rg *gin.RouterGroup) {

	// 涨跌统计
	rg.GET("/pxchg", publish2.NewStatistics().GET)

	// 融资融券
	rg.POST("/smt", publish2.NewSmt().POST)

	// 个股资金流向
	rg.POST("/capflow", publish2.NewCapitalflow().CapFlowSecuritySingle)

	// 市场分类资金流向
	rg.POST("/mkflow", publish2.NewCapitalflow().CapFlowMarket)

	// 资金趋势
	rg.POST("/captrend", publish2.NewCapTendency().POST)

	// K线上的除权除息
	rg.POST("/sdr", publish2.NewKlineXDXR().POST)

	// 攻击力度和攻击人气
	rg.POST("/gjx", publish2.NewLDRQ().POST)

	// F10首页
	rg.POST("/f10/home", f10.NewHN_F10_Mobile().GetF10_Mobile)

	// 公司详细信息
	rg.POST("/f10/comInfo", f10.NewCompany().GetF10_ComInfo)

	// 历史股本变动
	rg.POST("/f10/capstock/change", f10.NewCapitalStock().GetF10_CapitalStock)

	// 十大股东+十大流通股东
	rg.POST("/f10/holder/top10", f10.NewShareholderslTop10().GetShareholdersTop10)

}
