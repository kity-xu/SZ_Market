package diaData

import (
	"haina.com/market/f9/services/diaService"
	"haina.com/share/lib"
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
)

func GetDiaData(c *gin.Context) {
	scode := c.Query("scode")
	data, err := diaService.GetDiaData(scode)
	if err != nil {
		logging.Info(err.Error())
	}
	lib.WriteString(c, 200, data)
}
