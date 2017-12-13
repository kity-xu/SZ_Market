//资金流向
package publish2

import (
	"liveshow/share/lib"

	"haina.com/market/hqpublish/models/publish"
	"haina.com/market/hqpublish/models/publish2"

	. "haina.com/market/hqpublish/controllers"
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
)

type Capitalflow struct{}

func NewCapitalflow() *Capitalflow {
	return &Capitalflow{}
}

// 个股资金流向
func (this *Capitalflow) CapFlowSecuritySingle(c *gin.Context) {
	var _param struct {
		Sid int32 `json:"sid" binding:"required"`
	}

	if err := c.BindJSON(&_param); err != nil {
		lib.WriteString(c, 44001, nil)
		return
	}

	reply, err := publish2.NewFundflow(publish.REDISKEY_FUND_FLOW).GetFundflowReply(_param.Sid)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, reply)
}

// 市场分类的资金流向
func (this *Capitalflow) CapFlowMarket(c *gin.Context) {
	var _param struct {
		MarketID int32 `json:"marketID"`
		count    int32 `json:"count"`
	}
	if err := c.BindJSON(&_param); err != nil {
		lib.WriteString(c, 44001, nil)
		return
	}

	reply, err := publish2.NewMkCapflow().GetMkCapflow(_param.MarketID)
	if err != nil {
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, reply)
}
