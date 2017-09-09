//月线
package kline

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"
	"fmt"

	"haina.com/market/hqpublish/models/publish/kline"

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
		if err == publish.INVALID_FILE_PATH || err == publish.FILE_HMINDATA_NULL || err == publish.ERROR_KLINE_DATA_NULL {
			ret := &protocol.PayloadHisK{
				SID:   request.SID,
				Type:  request.Type,
				Total: -1,
				Num:   -1,
			}
			WriteJson(c, 200, ret)
			return
		}
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, reply)
}

func (this *Kline) MonthPB(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadKLineData(publish.REDISKEY_SECURITY_HMONTH, request)
	if err != nil {
		logging.Error("%v", err.Error())
		if err == publish.INVALID_FILE_PATH || err == publish.FILE_HMINDATA_NULL || err == publish.ERROR_KLINE_DATA_NULL {
			ret := &protocol.PayloadHisK{
				SID:   request.SID,
				Type:  request.Type,
				Total: -1,
				Num:   -1,
			}
			WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_HISKLINE, ret)
			return
		}
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_HISKLINE, reply)
}

func maybeAddMonthLine(reply *[]*protocol.KInfo, Sid int32, e error) error {
	if e == publish.INVALID_FILE_PATH { //可能是今天上市的新股
		key := fmt.Sprintf(publish.REDISKEY_SECURITY_NAME_ID, Sid) //去股票代码表查是否有此ID
		if !kline.IsExistInRedis(key) {
			return e
		}
		kinfo := &protocol.KInfo{
			NTime: models.GetCurrentTime(),
		}
		*reply = append(*reply, kinfo)
		return nil
	}

	if len(*reply) < 1 {
		return fmt.Errorf("PayloadHisK is null...")
	}

	var kinfo = protocol.KInfo{}
	kinfo = *(*reply)[len(*reply)-1]

	lday := kinfo.NTime
	today := models.GetCurrentTime()
	if kinfo.NSID/1000000 == 100 {
		if lday < Trade_100 { //是交易日
			if kinfo.NTime/100 != today/100 { //不同月
				kinfo.NTime = today
				kinfo.LlValue = 0
				kinfo.LlVolume = 0
				*reply = append(*reply, &kinfo)
			}
		}
	} else if kinfo.NSID/1000000 == 200 {
		if lday < Trade_200 {
			if kinfo.NTime/100 != today/100 { //不同月
				kinfo.NTime = today
				kinfo.LlValue = 0
				kinfo.LlVolume = 0
				*reply = append(*reply, &kinfo)
			}
		}
	} else {
		return fmt.Errorf("Invalid NSID...")
	}
	return nil
}
