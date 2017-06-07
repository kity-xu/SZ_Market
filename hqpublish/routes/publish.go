package routes

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/controllers/publish"
	"haina.com/market/hqpublish/controllers/publish/kline"
	"haina.com/market/hqpublish/controllers/publish/security"
)

func RegPublish(rg *gin.RouterGroup) {

	// 市场状态
	rg.POST("/market", publish.NewMarketStatus().POST) //默认pb模式

	// 分钟K线
	rg.POST("/min", publish.NewMinKLine().POST)

	//历史分钟K线
	rg.POST("/hismin", publish.NewHisMinLine().POST) //默认pb模式

	// 证券快照
	rg.POST("/snap", publish.NewStockSnapshot().POST)

	//市场、证券信息、股票代码表
	rg.POST("/sntab", security.NewSecurityTable().POST) //默认pb模式
	rg.POST("/sn", security.NewSecurityInfo().POST)     //默认pb模式
	rg.POST("/ssta", security.NewSecurityStatic().POST) //默认pb模式

	//历史K线
	rg.POST("/kline", kline.NewKline().POST) //默认pb模式

	//分笔成交正序
	rg.POST("/tradeet", publish.NewTradeEveryTime().POST)
}
