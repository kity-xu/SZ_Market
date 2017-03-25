package company

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models/company"
	. "haina.com/market/finance/models/finchina"
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

	enddate := c.Query(CONTEXT_END_DATE)
	count := c.Query(CONTEXT_COUNT)
	sCode := c.Query(CONTEXT_SECURITYCODE)

	value_int, err := strconv.Atoi(count)
	if err != nil {
		logging.Debug("%v", err)
		lib.WriteString(c, 88888, nil)
	}

	data, err := company.GetShareholderList(enddate, sCode, value_int)

	lib.WriteString(c, 200, data)
}

/**
获取十大流通股东信息
*/
func (this *EquityInfo) GetTop10Json(c *gin.Context) {

	enddate := c.Query(CONTEXT_END_DATE)
	count := c.Query(CONTEXT_COUNT)
	sCode := c.Query(CONTEXT_SECURITYCODE)

	value_int, err := strconv.Atoi(count)
	if err != nil {
		logging.Debug("%v", err)
		lib.WriteString(c, 88888, nil)
	}

	data, err := company.GetTop10List(enddate, sCode, value_int)

	lib.WriteString(c, 200, data)
}

/**
获取机构持股信息
*/
func (this *EquityInfo) GetOrganizationJson(c *gin.Context) {

	lib.WriteString(c, 200, "test")
}
