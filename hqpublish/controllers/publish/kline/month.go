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

	//	NSID     int32  `protobuf:"varint,1,opt,name=nSID" json:"nSID,omitempty"`
	//	NTime    int32  `protobuf:"varint,2,opt,name=nTime" json:"nTime,omitempty"`
	//	NPreCPx  int32  `protobuf:"varint,3,opt,name=nPreCPx" json:"nPreCPx,omitempty"`
	//	NOpenPx  int32  `protobuf:"varint,4,opt,name=nOpenPx" json:"nOpenPx,omitempty"`
	//	NHighPx  int32  `protobuf:"varint,5,opt,name=nHighPx" json:"nHighPx,omitempty"`
	//	NLowPx   int32  `protobuf:"varint,6,opt,name=nLowPx" json:"nLowPx,omitempty"`
	//	NLastPx  int32  `protobuf:"varint,7,opt,name=nLastPx" json:"nLastPx,omitempty"`
	//	LlVolume int64  `protobuf:"varint,8,opt,name=llVolume" json:"llVolume,omitempty"`
	//	LlValue  int64  `protobuf:"varint,9,opt,name=llValue" json:"llValue,omitempty"`
	//	NAvgPx   uint32 `protobuf:"varint,10,opt,name=nAvgPx" json:"nAvgPx,omitempty"`

	var kinfo = protocol.KInfo{}
	kinfo = *(*reply)[len(*reply)-1]

	lday := kinfo.NTime
	today := models.GetCurrentTime()
	if kinfo.NSID/1000000 == 100 {
		if lday < Trade_100 { //是交易日
			if kinfo.NTime/100 != today/100 { //不同月
				kinfo.NTime = today
				kinfo.NPreCPx = kinfo.NPreCPx
				kinfo.NOpenPx = kinfo.NPreCPx
				kinfo.NHighPx = kinfo.NPreCPx
				kinfo.NLowPx = kinfo.NPreCPx
				kinfo.NLastPx = kinfo.NPreCPx
				kinfo.LlValue = 0
				kinfo.LlVolume = 0
				kinfo.NAvgPx = 1
				*reply = append(*reply, &kinfo)
			}
		}
	} else if kinfo.NSID/1000000 == 200 {
		if lday < Trade_200 {
			if kinfo.NTime/100 != today/100 { //不同月
				kinfo.NTime = today
				kinfo.NPreCPx = kinfo.NPreCPx
				kinfo.NOpenPx = kinfo.NPreCPx
				kinfo.NHighPx = kinfo.NPreCPx
				kinfo.NLowPx = kinfo.NPreCPx
				kinfo.NLastPx = kinfo.NPreCPx
				kinfo.LlValue = 0
				kinfo.LlVolume = 0
				kinfo.NAvgPx = 1
				*reply = append(*reply, &kinfo)
			}
		}
	} else {
		return fmt.Errorf("Invalid NSID...")
	}
	return nil
}
