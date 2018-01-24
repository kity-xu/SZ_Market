package value

import (
	"haina.com/market/f9/services/valueService"
	"haina.com/share/lib"
	"haina.com/share/logging"

	"runtime"

	"github.com/gin-gonic/gin"
)

func GetValueData(c *gin.Context) {
	logging.Info("44444444444444444444==============", runtime.NumGoroutine())
	scode := c.Query("scode")
	//_, err := valueService.GetYingliData(scode)
	grow, err := valueService.GetValueData(scode)

	if err != nil {
		logging.Info(err.Error())
	}
	lib.WriteString(c, 200, grow)
}
