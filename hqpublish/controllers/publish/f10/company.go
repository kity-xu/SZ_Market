package f10

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/models/publish/f10"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type Company struct {
}

func NewCompany() *Company {
	return &Company{}
}

type Share struct {
	Scode     string      `json:"sid"`
	ComDetail interface{} `json:"comDetail"`
	Leader    interface{} `json:"leader"`
}

// 获取公司详细信息
func (this *Company) GetF10_ComInfo(c *gin.Context) {

	var _param struct {
		Scode int `json:"sid" binding:"required"`
	}

	if err := c.BindJSON(&_param); err != nil {
		logging.Debug("Bind Json | %v", err)
		lib.WriteString(c, 40004, nil)
		return
	}

	date, err := f10.GetF10Company(_param.Scode)
	if err != nil {
		logging.Error("%v", err)
	}
	lib.WriteString(c, 200, date)
}
