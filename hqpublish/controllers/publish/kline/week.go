//历史周K线
package kline

import (
	"fmt"

	"github.com/gin-gonic/gin"
	. "haina.com/market/hqpublish/controllers"
	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish/kline"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/logging"
)

func (this *Kline) WeekJson(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadKLineData(publish.REDISKEY_SECURITY_HWEEK, request)
	if err != nil {
		logging.Error("%v", err.Error())
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, reply)
}

func (this *Kline) WeekPB(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadKLineData(publish.REDISKEY_SECURITY_HWEEK, request)
	if err != nil {
		logging.Error("%v", err.Error())
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_HISKLINE, reply)
}

func maybeAddWeekLine(reply *[]*protocol.KInfo, Sid int32, e error) error {
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

			//判断是否同属一周
			b1, _ := models.DateAdd(kinfo.NTime) //找到该日期所在周日的那天
			b2, _ := models.DateAdd(today)

			if !b1.Equal(b2) { //是交易日但不同属一周（周一）新建
				kinfo.NTime = today
				kinfo.LlValue = 0
				kinfo.LlVolume = 0
				*reply = append(*reply, &kinfo)
			}
		}
	} else if kinfo.NSID/1000000 == 200 {
		if lday < Trade_200 {

			//判断是否同属一周
			b1, _ := models.DateAdd(kinfo.NTime) //找到该日期所在周日的那天
			b2, _ := models.DateAdd(today)

			if !b1.Equal(b2) { //是交易日但不同属一周（周一）新建
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
