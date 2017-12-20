package f10

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/models/publish/f10"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type ShareholderslTop10 struct {
}

func NewShareholderslTop10() *ShareholderslTop10 {
	return &ShareholderslTop10{}
}

// 历史股本变动
func (this *ShareholderslTop10) GetShareholdersTop10(c *gin.Context) {

	var _param struct {
		Scode   int `json:"sid" binding:"required"`
		Count   int `json:"count"`
		HType   int `json:"htype" binding:"required"`
		EndDate int `json:"enddate"`
	}

	if err := c.BindJSON(&_param); err != nil {
		logging.Debug("Bind Json | %v", err)
		lib.WriteString(c, 40004, nil)
		return
	}

	date, err := f10.GetHN_F10_ShareholdersTop10(_param.Scode, _param.HType, _param.EndDate)
	if err != nil {
		lib.WriteString(c, 40002, nil)
		return
	}

	lib.WriteString(c, 200, date)
}
