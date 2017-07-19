//移动端f10(指数K线图-f10 行情2)
package company

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type HN_F10_Mobile struct {
}

func NewHN_F10_Mobile() *HN_F10_Mobile {
	return &HN_F10_Mobile{}
}

type F10 struct {
	Scode  string      `json:"scode"`
	Name   *string     `json:"name"`
	Mobile interface{} `json:"f10"`
}

func (this *HN_F10_Mobile) GetF10_Mobile(c *gin.Context) {
	scode := c.Query(models.CONTEXT_SCODE)

	scodePrefix, market, err := ParseSCode(scode)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40004, nil)
		return
	}

	f10, name, err := company.F10Mobile(scodePrefix, market)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	result := &F10{
		Scode:  scode,
		Name:   name,
		Mobile: f10,
	}
	lib.WriteString(c, 200, result)
}
