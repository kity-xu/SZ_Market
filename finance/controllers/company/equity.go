package company

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models/company"
	"haina.com/market/finance/models/finchina"
	"haina.com/share/lib"
	"haina.com/share/logging"
	//"haina.com/share/store/redis"
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

	enddate := c.Query(finchina.CONTEXT_END_DATE)
	count := c.Query(finchina.CONTEXT_COUNT)
	scode := strings.Split(c.Query(finchina.CONTEXT_SECURITYCODE), ".")[0]
	value_int, err := strconv.Atoi(count)
	if err != nil {
		logging.Debug("%v", err)
		lib.WriteString(c, 88888, nil)
	}

	//根据条件查询股东信息
	data, err := company.GetShareholderList(enddate, scode, value_int)

	lib.WriteString(c, 200, data)
}

/**
获取十大流通股东信息
*/
func (this *EquityInfo) GetTop10Json(c *gin.Context) {

	enddate := c.Query(finchina.CONTEXT_END_DATE)
	count := c.Query(finchina.CONTEXT_COUNT)
	scode := strings.Split(c.Query(finchina.CONTEXT_SECURITYCODE), ".")[0]
	value_int, err := strconv.Atoi(count)
	if err != nil {
		logging.Debug("%v", err)
		lib.WriteString(c, 88888, nil)
	}

	data, err := company.GetTop10List(enddate, scode, value_int)

	lib.WriteString(c, 200, data)
}

/**
获取机构持股信息
*/
func (this *EquityInfo) GetOrganizationJson(c *gin.Context) {

	lib.WriteString(c, 200, "test")
}
