package kline

import (
	"github.com/gin-gonic/gin"
	. "haina.com/market/hqpublish/controllers"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/logging"
)

func (this *Kline) MinJson_15(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadMinLineData(publish.REDISKEY_SECURITY_HMIN15, request)
	if err != nil {
		logging.Error("%v", err.Error())
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, reply)
}

func (this *Kline) MinPB_15(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadMinLineData(publish.REDISKEY_SECURITY_HMIN15, request)
	if err != nil {
		logging.Error("%v", err.Error())
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_HISKLINE, reply)
}
