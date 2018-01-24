package routes

import (
	"haina.com/market/f9/controllers/bigData"
	"haina.com/market/f9/controllers/chip"
	"haina.com/market/f9/controllers/diaData"
	"haina.com/market/f9/controllers/fund"
	"haina.com/market/f9/controllers/isNew"
	"haina.com/market/f9/controllers/value"

	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	engine.GET("summary/isNewShares", isNew.GetNewData) //是否是新股

	engine.GET("summary/diagnose", diaData.GetDiaData) //概述

	engine.GET("summary/bigDate", bigData.GetBigData) //大数据透视

	engine.GET("valuePerspective", value.GetValueData) //价值透视

	engine.GET("getChipList", chip.GetChipData) //筹码成本

	engine.GET("getTrendList", chip.GetTrendList) //多空趋势

	engine.GET("fundPerspective", fund.GetFundData) //资金透视

	niu := engine.Group("/api")
	//注册公司信息获取路径
	RegNiuniu(niu)
}
