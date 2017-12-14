package publish2

import (
	"haina.com/share/lib"

	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/models/publish2"
)

type LDRQ struct {
}

func NewLDRQ() *LDRQ {
	return &LDRQ{}
}

func (*LDRQ) POST(c *gin.Context) {
	var _param struct {
		Sid int32 `json:"sid" binding:"required"`
	}

	if err := c.BindJSON(&_param); err != nil {
		lib.WriteString(c, 44001, nil)
		return
	}
	res, err := publish2.NewGJLDRQ().GetGJLDRQ(_param.Sid)
	if err != nil {
		lib.WriteString(c, 40002, nil)
		return
	}
	lib.WriteString(c, 200, res)
}
