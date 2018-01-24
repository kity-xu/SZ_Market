package chip

import (
	"haina.com/market/f9/services/chipService"
	"haina.com/share/lib"
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
)

func GetChipData(c *gin.Context) {
	scode := c.Query("scode")
	data, err := chipService.GetChipData(scode)
	if err != nil {
		logging.Info(err.Error())
	}
	lib.WriteString(c, 200, data)
}

func GetTrendList(c *gin.Context) {
	scode := c.Query("scode")
	data, err := chipService.GetTrendData(scode)
	if err != nil {
		logging.Info(err.Error())
	}
	logging.Info("lo====%+v", data)
	lib.WriteString(c, 200, data)
}
