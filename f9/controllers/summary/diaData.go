package summary

import (
	. "haina.com/market/f9/controllers"
	"haina.com/market/f9/services/diaService"
	"haina.com/share/lib"
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
)

func GetDiaData(c *gin.Context) {
	sid, err := Param_Norm_Sid(c)
	if err != nil {
		lib.WriteString(c, 40002, nil)
		return
	}
	data, err := diaService.GetDiaData(sid)
	if err != nil {
		logging.Info(err.Error())
	}
	lib.WriteString(c, 200, data)
}
