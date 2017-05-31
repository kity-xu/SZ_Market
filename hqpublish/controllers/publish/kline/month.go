//月线
package kline

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	. "haina.com/market/hqpublish/controllers"

	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/logging"
)

func (this *Kline) MonthJson(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadKLineData(publish.REDISKEY_SECURITY_HMONTH, request)
	if err != nil {
		logging.Error("%v", err.Error())
		WriteJson(c, 40002, nil)
		return
	}

	maybeAddKline(reply)
	WriteJson(c, 200, reply)
}

func (this *Kline) MonthPB(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadKLineData(publish.REDISKEY_SECURITY_HMONTH, request)
	if err != nil {
		logging.Error("%v", err.Error())
		WriteDataErrCode(c, 40002)
		return
	}

	maybeAddKline(reply)
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_HISKLINE, reply)
}
