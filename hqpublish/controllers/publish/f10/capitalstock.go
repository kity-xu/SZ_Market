// 历史股本变动
package f10

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/models/publish/f10"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type CapitalStock struct {
}

func NewCapitalStock() *CapitalStock {
	return &CapitalStock{}
}

func (this *CapitalStock) GetF10_CapitalStock(c *gin.Context) {

	var _param struct {
		Scode int `json:"sid" binding:"required"`
		Count int `json:"count"`
	}

	if err := c.BindJSON(&_param); err != nil {
		logging.Debug("Bind Json | %v", err)
		lib.WriteString(c, 40004, nil)
		return
	}

	csdate, err := f10.GetF10CapitalStock(_param.Scode, 10)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40002, err)
		return
	}

	lib.WriteString(c, 200, csdate)
}
