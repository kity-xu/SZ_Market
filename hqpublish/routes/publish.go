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

	//排序
	rg.POST("/sort", publish.NewSort().POST)

	//资金流向
	rg.POST("/fundflow", publish.NewFundflow().POST)

	//板块及板块成分
	rg.POST("/block", publish.NewStockBlock().POST)
	rg.POST("/element", publish.NewStockElement().POST)

	//A股市场代码表、市场代码表、证券基本信息、股票静态数据
	rg.GET("/sntab/astock", security.NewSecurityTable().GET) //默认pb模式

	rg.POST("/sntab", security.NewSecurityTable().POST) //默认pb模式
	rg.POST("/sn", security.NewSecurityInfo().POST)     //默认pb模式
	rg.POST("/ssta", security.NewSecurityStatic().POST) //默认pb模式

	//历史K线
	rg.POST("/kline", kline.NewKline().POST) //默认pb模式

	//分笔成交 正序
	rg.POST("/tradeet", publish.NewTradeEveryTime().POST)
	//最近分笔成交 逆序
	rg.POST("/tradeetnow", publish.NewTradeEveryTimeNow().POST)

	// 信息栏 -zxw
	rg.POST("/infobar", publish.NewInfoBar().POST)
	// 证券集合(板块) -zxw
	rg.POST("/subset", publish.NewStockBlockSet().POST)
	// 公告信息集合 -zxw
	rg.POST("/hisevent", publish.NewNoticeInfo().POST)
	// 单条公告信息 -zxw
	rg.POST("/hiseventid", publish.NewHisEvent().POST)
	// 个股详情 移动端 -zxw
	rg.POST("/persdetail", publish.NewPerSDetail().POST)
	// 分价成交 -zxw
	rg.POST("/tradedp", publish.NewTradePriceRecordC().POST)

	// 除权除息
	rg.POST("/xrxd", publish.NewXRXD().POST)
	// 复权因子
	rg.POST("/factor", publish.NewFactor().POST)

	// 移动端首页
	rg.GET("/mindex", publish.NewMIndex().GET)
}
