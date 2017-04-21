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
	ntype := c.Query(models.CONTEXT_NTYPE) // 查询类型 d0 全部 d1 一季度  d2 半年报 d3 三季度  d4 年报

	logging.Info("参数 %v", ntype)

	data, err := company.GetStructure(scode, ntype)
	if err != nil {
		lib.WriteString(c, 300, err.Error())
		return
	}
	lib.WriteString(c, 200, data)
}

/**
获取股本变动信息
*/
func (this *CapitalizationInfo) GetChangesJson(c *gin.Context) {

	enddate := c.Query(models.CONTEXT_END_DATE)
	count := c.Query(models.CONTEXT_COUNT)
	scode := strings.Split(c.Query(models.CONTEXT_SCODE), ".")[0]
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
	data, err := company.GetChangesStrInfo(enddate, scode, value_int)
	if err != nil {
		lib.WriteString(c, 300, err.Error())
		return
	}
	lib.WriteString(c, 200, data)
}
