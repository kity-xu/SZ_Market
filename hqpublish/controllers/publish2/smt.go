// 融资融券
package publish2

import (
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/models/publish2"
	"haina.com/share/lib"
)

type SMT struct {
}

func NewSmt() *SMT {
	return &SMT{}
}

func (SMT) POST(c *gin.Context) {
	var _param struct {
		Which int32 `json:"marketId" binding:"required"`
		Count int32 `json:"count"`
	}

	if err := c.BindJSON(&_param); err != nil {
		lib.WriteString(c, 44001, nil)
		return
	}

	marketID := getExchageByReq(_param.Which)
	if marketID == "" {
		logging.Error("SMTS:Invalid request param 'Which'")
		lib.WriteString(c, 40002, nil)
		return
	}
	result := publish2.GetSMTbyMarket(_param.Which, marketID)
	if result == nil {
		lib.WriteString(c, 40002, nil)
		return
	}
	lib.WriteString(c, 200, result)
}

func getExchageByReq(which int32) string {
	switch which {
	case 100:
		return "001002"
	case 200:
		return "001003"
	case 300:
		return "001000"
	}
	return ""
}
