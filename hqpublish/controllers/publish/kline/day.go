//日线
package kline

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"
	"fmt"

	"github.com/gin-gonic/gin"
	. "haina.com/market/hqpublish/controllers"
	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"
	"haina.com/market/hqpublish/models/publish/kline"
	"haina.com/share/logging"
)

func (this *Kline) DayJson(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadKLineData(publish.REDISKEY_SECURITY_HDAY, request)
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

func (this *Kline) DayPB(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadKLineData(publish.REDISKEY_SECURITY_HDAY, request)
	if err != nil {
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

func (this *Kline) PayLoadKLineData(redisKey string, request *protocol.RequestHisK) (*protocol.PayloadHisK, error) {
	var ret *protocol.PayloadHisK
	dlines, e := kline.NewKLine(redisKey).GetHisKLineAll(request)
	if e == publish.INVALID_FILE_PATH || e == publish.FILE_HMINDATA_NULL || e == publish.ERROR_KLINE_DATA_NULL {
		return nil, e
	}
	if e == nil {
		models.GetASCStruct(dlines) //升序排序
	}

	logging.Info("Create new kline...")
	switch redisKey {
	case publish.REDISKEY_SECURITY_HDAY:
		e = maybeAddKline(dlines, request.SID, e)
		if e != nil {
			logging.Error(e.Error())
			return nil, e
		}
		break
	case publish.REDISKEY_SECURITY_HWEEK:
		e = maybeAddWeekLine(dlines, request.SID, e)
		if e != nil {
			logging.Error(e.Error())
			return nil, e
		}
		break
	case publish.REDISKEY_SECURITY_HMONTH:
		e = maybeAddMonthLine(dlines, request.SID, e)
		if e != nil {
			logging.Error(e.Error())
			return nil, e
		}
		break
	case publish.REDISKEY_SECURITY_HYEAR:
		e = maybeAddYearLine(dlines, request.SID, e)
		if e != nil {
			logging.Error(e.Error())
			return nil, e
		}
		break
	}

	if len(*dlines) == 0 {
		return nil, e
	}

	total := int32(len(*dlines))

	if request.Num > total {
		request.Num = total
	}

	if request.Num == 0 { //num==0, 获取全部
		ret = &protocol.PayloadHisK{
			SID:   request.SID,
			Type:  request.Type,
			Total: total,
			Begin: request.TimeIndex,
			Num:   total,
			KList: *dlines,
		}
		return ret, nil
	} else { //根据num, 获取部分
		if request.TimeIndex == 0 { //起始日期最新
			var table []*protocol.KInfo

			lindex := total
			for i := lindex - request.Num; i < lindex; i++ {
				table = append(table, (*dlines)[i])
			}
			ret = &protocol.PayloadHisK{ //向前
				SID:   request.SID,
				Type:  request.Type,
				Total: total,
				Begin: table[0].NTime,
				Num:   int32(len(table)),
				KList: table,
			}

			if request.Direct == 1 { //向后
				var sig []*protocol.KInfo
				sig = append(sig, table[len(table)-1])
				ret = &protocol.PayloadHisK{
					SID:   request.SID,
					Type:  request.Type,
					Total: total,
					Begin: table[len(table)-1].NTime,
					Num:   1,
					KList: sig,
				}
			}
			return ret, nil

		} else { //TimeIndex作为起始日期
			var frontedSwap, palinalSwap []*protocol.KInfo
			var databuf []*protocol.KInfo

			for _, v := range *dlines {
				if v.NTime <= request.TimeIndex {
					frontedSwap = append(frontedSwap, v)
				}
				if v.NTime >= request.TimeIndex {
					palinalSwap = append(palinalSwap, v)
				}
			}

			if request.Direct == 0 { //向前 frontedSwap
				size := len(frontedSwap)
				if size < int(request.Num) {
					databuf = frontedSwap
				} else {
					for i := size - int(request.Num); i < size; i++ {
						databuf = append(databuf, frontedSwap[i])
					}
				}
			} else if request.Direct == 1 { //向后 palinalSwap
				if len(palinalSwap) == 0 { //不加此判断 最新日期向后取，会越界panic
					var table []*protocol.KInfo

					lindex := total
					for i := lindex - request.Num; i < lindex; i++ {
						table = append(table, (*dlines)[i])
					}

					var sig []*protocol.KInfo
					sig = append(sig, table[len(table)-1])

					ret = &protocol.PayloadHisK{
						SID:   request.SID,
						Type:  request.Type,
						Total: total,
						Begin: table[len(table)-1].NTime,
						Num:   1,
						KList: sig,
					}
					return ret, nil
				}

				if int(request.Num) > len(palinalSwap) {
					for i := 0; i < len(palinalSwap); i++ {
						databuf = append(databuf, palinalSwap[i])
					}
				} else {
					for i := 0; i < int(request.Num); i++ {
						databuf = append(databuf, palinalSwap[i])
					}
				}
			} else {
				return nil, ERROR_REQUEST_PARAM
			}
			ret = &protocol.PayloadHisK{
				SID:   request.SID,
				Type:  request.Type,
				Total: total,
				Begin: databuf[0].NTime,
				Num:   int32(len(databuf)),
				KList: databuf,
			}
			return ret, nil
		}
	}
}

//新增K线
func maybeAddKline(reply *[]*protocol.KInfo, Sid int32, e error) error {
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
	logging.Debug("Trade_100:%v------Trade_200:%v", Trade_100, Trade_200)
	var kinfo = protocol.KInfo{}
	kinfo = *(*reply)[len(*reply)-1]

	lday := kinfo.NTime
	today := models.GetCurrentTime()

	if kinfo.NSID/1000000 == 100 {
		if lday < Trade_100 { //如果K线最后一天的日期小于交易日  则新增
			kinfo.NTime = today
			kinfo.NPreCPx = kinfo.NLastPx
			kinfo.NOpenPx = kinfo.NLastPx
			kinfo.NHighPx = kinfo.NLastPx
			kinfo.NLowPx = kinfo.NLastPx
			kinfo.NLastPx = kinfo.NLastPx
			kinfo.LlValue = 0
			kinfo.LlVolume = 0
			kinfo.NAvgPx = 1
			logging.Info("新增一根K线-----Trade_100：%v", Trade_100)
			*reply = append(*reply, &kinfo)
		}
	} else if kinfo.NSID/1000000 == 200 {
		if lday < Trade_200 {
			kinfo.NTime = today
			kinfo.NPreCPx = kinfo.NLastPx
			kinfo.NOpenPx = kinfo.NLastPx
			kinfo.NHighPx = kinfo.NLastPx
			kinfo.NLowPx = kinfo.NLastPx
			kinfo.NLastPx = kinfo.NLastPx
			kinfo.LlValue = 0
			kinfo.LlVolume = 0
			kinfo.NAvgPx = 1
			logging.Info("新增一根K线-----Trade_100：%v", Trade_200)
			*reply = append(*reply, &kinfo)
		}
	} else {
		return fmt.Errorf("Invalid NSID...")
	}
	return nil
}

func MaybeAddKline(reply *[]*protocol.KInfo, Sid int32, e error) error {
	return maybeAddKline(reply, Sid, e)
}
