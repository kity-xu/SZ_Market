package company

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type EquityInfo struct {
}

func NewEquityInfo() *EquityInfo {
	return &EquityInfo{}
}

/**
获取股东人数信息
*/
func (this *EquityInfo) GetShareholderJson(c *gin.Context) {

	enddate := c.Query(models.CONTEXT_END_DATE)
	count := c.Query(models.CONTEXT_COUNT)

	scodePrefix, market, err := ParseSCode(c.Query(models.CONTEXT_SCODE))

	if err != nil {
		lib.WriteString(c, 40004, "")
		return
	}
	var value_int = 0
	var erre error
	if count != "" {
		value_int, erre = strconv.Atoi(count)
		logging.Info("类型转换%v", erre)
	} else {
		value_int = 10
	}

	//根据条件查询股东信息
	var strDate = ""
	if enddate != "" {
		strDate = " and ENDDATE<'" + enddate + "'"
	}
	data, err := company.GetShareholderGroup(scodePrefix, value_int, strDate, market)
	lib.WriteString(c, 200, data)
}

/**
获取十大流通股东信息
*/
func (this *EquityInfo) GetTop10Json(c *gin.Context) {

	enddate := c.Query(models.CONTEXT_END_DATE)
	count := c.Query(models.CONTEXT_COUNT)
	scodePrefix, market, err := ParseSCode(c.Query(models.CONTEXT_SCODE))

	if err != nil {
		lib.WriteString(c, 40004, "")
		return
	}
	var value_int = 0
	var erre error
	if count == "" {
		value_int = 10
	} else {
		value_int, erre = strconv.Atoi(count)
	}

	if erre != nil {
		logging.Debug("%v", err)
		lib.WriteString(c, 88888, nil)
	}

	data, err := company.GetTop10Group(enddate, scodePrefix, value_int, market)

	lib.WriteString(c, 200, data)
}

/**
获取机构持股信息
*/
func (this *EquityInfo) GetOrganizationJson(c *gin.Context) {
	scodePrefix, market, err := ParseSCode(c.Query(models.CONTEXT_SCODE))

	if err != nil {
		lib.WriteString(c, 40004, "")
		return
	}
	data, err := company.GetCompGroup(scodePrefix, market)
	if err != nil {
		lib.WriteString(c, 300, err.Error())
		return
	}
	lib.WriteString(c, 200, data)
}
