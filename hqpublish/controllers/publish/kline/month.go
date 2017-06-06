//月线
package kline

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	. "haina.com/market/hqpublish/controllers"
	"haina.com/market/hqpublish/models"

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
	WriteJson(c, 200, reply)
}

func (this *Kline) MonthPB(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadKLineData(publish.REDISKEY_SECURITY_HMONTH, request)
	if err != nil {
		logging.Error("%v", err.Error())
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_HISKLINE, reply)
}

func maybeAddMonthLine(reply *[]*protocol.KInfo) {
	if len(*reply) < 1 {
		logging.Error("PayloadHisK is null...")
		return
	}

	today := models.GetCurrentTime()
	if (*reply)[0].NSID/1000000 == 100 {
		if today == Trade_100 { //是交易日
			var kinfo = protocol.KInfo{}
			kinfo = *(*reply)[len(*reply)-1]

			if kinfo.NTime/100 != today/100 { //不同月
				kinfo.NTime = today
				kinfo.LlValue = 0
				kinfo.LlVolume = 0
				*reply = append(*reply, &kinfo)
			}
		}
	} else if (*reply)[0].NSID/1000000 == 200 {
		if today == Trade_200 {
			var kinfo = protocol.KInfo{}
			kinfo = *(*reply)[len(*reply)-1]

			if kinfo.NTime/100 != today/100 { //不同月
				kinfo.NTime = today
				kinfo.LlValue = 0
				kinfo.LlVolume = 0
				*reply = append(*reply, &kinfo)
			}
		}
	} else {
		logging.Error("Invalid NSID...")
	}
}
