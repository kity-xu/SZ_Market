package company

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models/company"
	. "haina.com/market/finance/models/finchina"
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
	sCode := c.Query(CONTEXT_SECURITYCODE)

	data, err := company.GetStructure(sCode)
	if err != nil {
		fmt.Println(err)
	}
	lib.WriteString(c, 200, data)
}

/**
获取股本变动信息
*/
func (this *CapitalizationInfo) GetChangesJson(c *gin.Context) {

	enddate := c.Query(CONTEXT_END_DATE)
	count := c.Query(CONTEXT_COUNT)
	sCode := c.Query(CONTEXT_SECURITYCODE)
	value_int, err := strconv.Atoi(count)
	if err != nil {
		logging.Debug("%v", err)
		lib.WriteString(c, 88888, nil)
	}
	data, err := company.GetChangesStrInfo(enddate, sCode, value_int)
	if err != nil {
		fmt.Println(err)
	}
	lib.WriteString(c, 200, data)
}
