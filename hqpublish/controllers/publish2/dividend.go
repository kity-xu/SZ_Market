package publish2

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/models/publish2"
	"haina.com/share/lib"
)

type KlineXDXR struct{}

func NewKlineXDXR() *KlineXDXR {
	return &KlineXDXR{}
}

func (*KlineXDXR) POST(c *gin.Context) {
	var _param struct {
		Sid int32 `json:"sid"`
	}
	if err := c.BindJSON(&_param); err != nil {
		lib.WriteString(c, 44001, nil)
		return
	}
	res, err := publish2.NewDividendJson().GetDividendJson(_param.Sid)
	if err != nil {
		lib.WriteString(c, 40002, nil)
		return
	}
	lib.WriteString(c, 200, res)
}
