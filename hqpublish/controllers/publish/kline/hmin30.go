package kline

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"github.com/gin-gonic/gin"

	. "haina.com/market/hqpublish/controllers"
	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/logging"
)

func (this *Kline) MinJson_30(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadMinLineData(publish.REDISKEY_SECURITY_HMIN30, request)
	if err != nil {
		logging.Error("%v", err.Error())
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, reply)

}

func (this *Kline) MinPB_30(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadMinLineData(publish.REDISKEY_SECURITY_HMIN30, request)
	if err != nil {
		logging.Error("%v", err.Error())
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_HISKLINE, reply)

}
