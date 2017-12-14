package routes

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/controllers/publish2"
)

func RegFinance(engine *gin.Engine) {

	rg := engine.Group("/api/finance")

	// 财务-图表
	rg.POST("/chart", publish2.NewFinanceChart().POST)

	// 财务-报表
	rg.POST("/report", publish2.NewFinanceReport().POST)

	// 研报-评级统计
	rg.POST("/report/statistics", publish2.NewReportStatistics().POST)

	// 研报-盈利预测
	rg.POST("/report/forecast", publish2.NewReportForecast().POST)
}
