// 关键指标
package company

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type IndicatorsInfo struct {
}

func NewIndicatorsInfo() *IndicatorsInfo {
	return &IndicatorsInfo{}
}

func (this *IndicatorsInfo) getJson(c *gin.Context) (*company.ResponseInfo, error) {
	return company.NewIndicators().GetJson(c)
}

func (this *IndicatorsInfo) GET(c *gin.Context) {
	data, err := this.getJson(c)
	if err != nil {
		logging.Debug("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	lib.WriteString(c, 200, data)
}
