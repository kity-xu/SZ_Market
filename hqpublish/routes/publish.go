package routes

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/controllers/publish"
)

func RegPublish(rg *gin.RouterGroup) {

	// 股票代码表
	rg.GET("/securitytable", publish.NewSecurityTable().GET) //默认pb模式

	// 日K线历史数据
	rg.GET("/kline/dayk", publish.NewDayKLine().GET) //默认pb模式

}
