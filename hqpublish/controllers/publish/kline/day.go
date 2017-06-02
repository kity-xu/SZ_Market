//日线
package kline

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"github.com/gin-gonic/gin"
	. "haina.com/market/hqpublish/controllers"
	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/logging"
)

func (this *Kline) DayJson(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadKLineData(publish.REDISKEY_SECURITY_HDAY, request)
	if err != nil {
		logging.Error("%v", err.Error())
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, reply)
}

func (this *Kline) DayPB(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadKLineData(publish.REDISKEY_SECURITY_HDAY, request)
	if err != nil {
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_HISKLINE, reply)
}

func (this *Kline) PayLoadKLineData(redisKey string, request *protocol.RequestHisK) (*protocol.PayloadHisK, error) {
	var ret *protocol.PayloadHisK
	dlines, err := publish.NewKLine(redisKey).GetHisKLineAll(request)
	if err != nil {
		return nil, err
	}

	//新增K线的判断   待改。。。。
	//if publish.IsHQpostRunOver() {
	if models.GetCurrentTimeHM()%10000 < 1530 {
		switch redisKey {
		case publish.REDISKEY_SECURITY_HDAY:
			maybeAddKline(dlines)
			break
		case publish.REDISKEY_SECURITY_HWEEK:
			maybeAddWeekLine(dlines)
			break
		case publish.REDISKEY_SECURITY_HMONTH:
			maybeAddMonthLine(dlines)
			break
		case publish.REDISKEY_SECURITY_HYEAR:
			maybeAddYearLine(dlines)
			break
		}
	}

	total := int32(len(*dlines))

	if request.Num > total {
		request.Num = total
	}
	models.GetASCStruct(dlines) //升序排序
	if request.Num == 0 {       //num==0, 获取全部
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
func maybeAddKline(reply *[]*protocol.KInfo) {
	if len(*reply) < 1 {
		logging.Error("PayloadHisK is null...")
		return
	}

	today := models.GetCurrentTime()
	if (*reply)[0].NSID/1000000 == 100 {
		if today == Trade_100 {
			var kinfo = protocol.KInfo{}

			kinfo = *(*reply)[len(*reply)-1]
			kinfo.NTime = today
			kinfo.LlValue = 0
			kinfo.LlVolume = 0
			*reply = append(*reply, &kinfo)
		}
	} else if (*reply)[0].NSID/1000000 == 200 {
		if today == Trade_200 {
			var kinfo = protocol.KInfo{}

			kinfo = *(*reply)[len(*reply)-1]
			kinfo.NTime = today
			kinfo.LlValue = 0
			kinfo.LlVolume = 0
			*reply = append(*reply, &kinfo)
		}
	} else {
		logging.Error("Invalid NSID...")
	}
}
