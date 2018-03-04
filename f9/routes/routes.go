package routes

import (
	"haina.com/market/f9/controllers/basic"
	"haina.com/market/f9/controllers/summary"

	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	rg := engine.Group("/api/f9")

	rg.GET("summary/isNewShares", basic.GetNewData) //是否是新股

	rg.GET("summary/diagnose", summary.GetDiaData) //概述

	rg.GET("summary/bigDate", summary.GetBigData) //大数据透视

	rg.GET("valuePerspective", summary.GetValueData) //价值透视

	rg.GET("getChipList", summary.GetChipData) //筹码成本

	rg.GET("getTrendList", summary.GetTrendList) //多空趋势

	rg.GET("fundPerspective", summary.GetFundData) //资金透视
}
