//年线
package kline

import (
	"fmt"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/market/hqpublish/models/publish/kline"

	"github.com/gin-gonic/gin"
	. "haina.com/market/hqpublish/controllers"
	"haina.com/market/hqpublish/models"

	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/logging"
)

func (this *Kline) YearJson(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadKLineData(publish.REDISKEY_SECURITY_HYEAR, request)
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

func (this *Kline) YearPB(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadKLineData(publish.REDISKEY_SECURITY_HYEAR, request)
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

func MaybeAddYearLine(reply *[]*protocol.KInfo, Sid int32, e error) error {
	return maybeAddYearLine(reply, Sid, e)
}

func maybeAddYearLine(reply *[]*protocol.KInfo, Sid int32, e error) error {
	is, err := kline.IsIndex(Sid)
	if err != nil {
		logging.Debug("%v", e.Error())
	}
	if err == nil && !is {
		if kline.IsDelist(Sid) { // 停盘
			return nil
		}
	}

	if e == publish.INVALID_FILE_PATH { //可能是今天上市的新股
		key := fmt.Sprintf(publish.REDISKEY_SECURITY_NAME_ID, Sid) //去股票代码表查是否有此ID
		if !kline.IsExistInRedis(key) {
			return e
		}
		kinfo := &protocol.KInfo{
			NTime:  models.GetCurrentTime(),
			NAvgPx: 1,
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
	//today := models.GetCurrentTime()
	if lday < Trade_100 { //是交易日
		if kinfo.NTime/10000 != Trade_100/10000 { //不同年
			kinfo.NTime = Trade_100
			kinfo.NPreCPx = kinfo.NLastPx
			kinfo.NOpenPx = kinfo.NLastPx
			kinfo.NHighPx = kinfo.NLastPx
			kinfo.NLowPx = kinfo.NLastPx
			kinfo.NLastPx = kinfo.NLastPx
			kinfo.LlValue = 0
			kinfo.LlVolume = 0
			kinfo.NAvgPx = 1
			*reply = append(*reply, &kinfo)
		} else {
			(*reply)[len(*reply)-1].NTime = Trade_100
		}
	}
	return nil
}
