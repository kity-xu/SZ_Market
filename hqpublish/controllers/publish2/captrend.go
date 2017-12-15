// 资金趋势
package publish2

import (
	"haina.com/market/hqpublish/models/publish2"
	"haina.com/share/lib"

	"github.com/gin-gonic/gin"
)

type CapTendency struct {
}

func NewCapTendency() *CapTendency {
	return &CapTendency{}
}

func (*CapTendency) POST(c *gin.Context) {
	var _param struct {
		Sid       int32 `json:"sid"`
		NType     int32 `json:"type"`
		TimeIndex int32 `json:"timeIndex"`
		Num       int32 `json:"num"`
		Direct    int32 `json:"direct"`
	}
	if err := c.BindJSON(&_param); err != nil {
		lib.WriteString(c, 44001, nil)
		return
	}

	ld, err := publish2.NewPeriodCapFlow().GetPeriodCapFlowList(_param.Sid, _param.NType, _param.TimeIndex, _param.Num, _param.Direct)
	if err != nil {
		lib.WriteString(c, 40002, nil)
		return
	}
	lib.WriteString(c, 200, ld)
}
