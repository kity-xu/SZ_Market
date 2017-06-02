package kline

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"github.com/gin-gonic/gin"
	. "haina.com/market/hqpublish/controllers"
	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/logging"
)

func (this *Kline) MinJson_01(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadMinLineData(publish.REDISKEY_SECURITY_HMIN, request)
	if err != nil {
		logging.Error("%v", err.Error())
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, reply)
}

func (this *Kline) MinPB_01(c *gin.Context, request *protocol.RequestHisK) {
	reply, err := this.PayLoadMinLineData(publish.REDISKEY_SECURITY_HMIN, request)
	if err != nil {
		logging.Error("%v", err.Error())
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_HISKLINE, reply)
}

func (this *Kline) PayLoadMinLineData(redisKey string, request *protocol.RequestHisK) (*protocol.PayloadHisK, error) {
	var ret *protocol.PayloadHisK
	mlines, err := publish.NewHMinKLine(redisKey).GetHMinKLineAll(request)
	if err != nil {
		return nil, err
	}
	total := int32(len(*mlines))

	if request.Num > total {
		request.Num = total
	}
	models.GetASCStruct(mlines) //升序排序
	if request.Num == 0 {       //num==0, 获取全部
		ret = &protocol.PayloadHisK{
			SID:   request.SID,
			Type:  request.Type,
			Total: total,
			Begin: request.TimeIndex,
			Num:   total,
			KList: *mlines,
		}
		return ret, nil
	} else { //根据num, 获取部分
		if request.TimeIndex == 0 { //起始日期最新
			var table []*protocol.KInfo

			lindex := total
			for i := lindex - request.Num; i < lindex; i++ {
				table = append(table, (*mlines)[i])
			}

			if len(table) < 1 {
				return nil, ERROR_INDEX_MAYBE_OUTOF_RANGE
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

			for _, v := range *mlines {
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
						table = append(table, (*mlines)[i])
					}

					var sig []*protocol.KInfo
					sig = append(sig, table[len(table)-1])

					if len(table) < 1 {
						return nil, ERROR_INDEX_MAYBE_OUTOF_RANGE
					}
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

			if len(databuf) < 1 {
				return nil, ERROR_INDEX_MAYBE_OUTOF_RANGE
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
