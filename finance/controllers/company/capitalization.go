package company

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type CapitalizationInfo struct {
}

func NewCapitalizationInfo() *CapitalizationInfo {
	return &CapitalizationInfo{}
}

/**
获取股本结构信息
*/
func (this *CapitalizationInfo) GetStructureJson(c *gin.Context) {
	scode := strings.Split(c.Query(models.CONTEXT_SCODE), ".")[0]
	count := c.Query(models.CONTEXT_COUNT)
	scodePrefix, market, err := ParseSCode(c.Query(models.CONTEXT_SCODE))

	if err != nil {
		lib.WriteString(c, 40004, "")
		return
	}
	var value_int = 0
	var err error
	if count != "" {
		value_int, err = strconv.Atoi(count)
	} else {
		value_int = 8
	}

	if err != nil {
		//logging.Debug("%v", err)
		return
	}
	ntype := c.Query(models.CONTEXT_NTYPE) // 查询类型 d0 全部 d1 一季度  d2 半年报 d3 三季度  d4 年报
	var selwhe = ""
	if ntype == "d0" || ntype == "" { // 查询全部
		selwhe = ""
	}
	if ntype == "d1" { // 一季度
		selwhe = " and ENDDATE LIKE '%0330' "
	}
	if ntype == "d2" { // 半年
		selwhe = " and ENDDATE LIKE '%0629' "
	}
	if ntype == "d3" { // 三季度
		selwhe = " and ENDDATE LIKE '%0929' "
	}
	if ntype == "d4" { // 年度
		selwhe = " and ENDDATE LIKE '%1230' "
	}

	data, err := company.GetStructure(scode, selwhe, value_int, market)
	if err != nil {
		lib.WriteString(c, 300, err.Error())
		return
	}

	var rtj company.RetTrucsInfoJson
	rtj.SCode = scode
	rtj.TrucsList = data

	lib.WriteString(c, 200, &rtj)
}

/**
获取股本变动信息
*/
func (this *CapitalizationInfo) GetChangesJson(c *gin.Context) {

	enddate := c.Query(models.CONTEXT_END_DATE)
	count := c.Query(models.CONTEXT_COUNT)
	scode := strings.Split(c.Query(models.CONTEXT_SCODE), ".")[0]
	scodePrefix, market, err := ParseSCode(c.Query(models.CONTEXT_SCODE))

	if err != nil {
		lib.WriteString(c, 40004, "")
		return
	}
	var value_int = 0
	var err error
	if count == "" {
		value_int = 10
	} else {
		value_int, err = strconv.Atoi(count)
	}
	if err != nil {
		logging.Debug("%v", err)
		lib.WriteString(c, 88888, nil)
	}
	data, err := company.GetChangesStrInfo(enddate, scode, value_int, market)
	if err != nil {
		lib.WriteString(c, 300, err.Error())
		return
	}
	lib.WriteString(c, 200, data)
}

// 判断市场类型
func ReturnExchang(np string) string {
	ind := strings.Index(np, ".")
	if ind < 1 {
		return "error"
	}
	ntype := strings.Split(np, ".")[1]
	var exchange = ""
	if len(ntype) > 1 {
		switch strings.ToUpper(ntype) {
		case "SH":
			exchange = "SH"
		case "SZ":
			exchange = "SZ"
		default:
			logging.Info("其他市场股票，或SCODE有误！")
		}
	}
	return exchange
}
