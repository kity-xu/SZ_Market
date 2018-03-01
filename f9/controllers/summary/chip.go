package summary

import (
	"github.com/gin-gonic/gin"
	. "haina.com/market/f9/controllers"
	"haina.com/market/f9/services/chipService"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

func GetChipData(c *gin.Context) {
	sid, err := Param_Norm_Sid(c)
	if err != nil {
		lib.WriteString(c, 40002, nil)
		return
	}
	data, err := chipService.GetChipData(sid)
	if err != nil {
		logging.Info(err.Error())
	}
	lib.WriteString(c, 200, data)
}

func GetTrendList(c *gin.Context) {
	sid, err := Param_Norm_Sid(c)
	if err != nil {
		lib.WriteString(c, 40002, nil)
		return
	}
	data, err := chipService.GetTrendData(sid)
	if err != nil {
		logging.Info(err.Error())
	}
	lib.WriteString(c, 200, data)
}
