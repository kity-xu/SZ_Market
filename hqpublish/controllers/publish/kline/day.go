//日线
package kline

import (
	"ProtocolBuffer/format/kline"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"haina.com/market/hqpublish/models"

	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

func (this *Kline) DayJson(c *gin.Context, request *kline.RequestHisK) {
	reply := this.ReplyKLine(c, publish.REDISKEY_SECURITY_HDAY, request)

	c.JSON(http.StatusOK, reply)
}

func (this *Kline) DayPB(c *gin.Context, request *kline.RequestHisK) {
	reply := this.ReplyKLine(c, publish.REDISKEY_SECURITY_HDAY, request)

	//转PB
	replypb, err := proto.Marshal(reply)
	if err != nil {
		reply := kline.ReplyHisK{
			Code: 40002,
		}
		replypb, err = proto.Marshal(&reply)
		if err != nil {
			logging.Error("pb marshal error: %v", err)
		}
		lib.WriteData(c, replypb)
		return

	}
	lib.WriteData(c, replypb)
}

func (this *Kline) ReplyKLine(c *gin.Context, redisKey string, request *kline.RequestHisK) *kline.ReplyHisK {
	reply := &kline.ReplyHisK{}

	dlines, _, err := publish.NewKLine(redisKey).GetHisKLine(request.SID)
	if err != nil {
		logging.Error("%v", err)
		return &kline.ReplyHisK{
			Code: 4002,
			Data: &kline.HisK{},
		}
	}

	if int(request.Num) > len(dlines.List) {
		request.Num = int32(len(dlines.List))
		//		logging.Error("Invalid request parameters 'Num'...")
		//		return &kline.ReplyHisK{
		//			Code: 4002,
		//			Data: &kline.HisK{},
		//		}
	}

	models.GetASCStruct(&dlines.List) //升序排序

	if request.Num == 0 { //num==0, 获取全部
		reply = &kline.ReplyHisK{
			Code: 200,
			Data: &kline.HisK{
				SID:   request.SID,
				Type:  request.Type,
				Total: int32(len(dlines.List)),
				Begin: dlines.List[0].NTime,
				Num:   int32(len(dlines.List)),
				List:  dlines.List,
			},
		}

	} else { //根据num, 获取部分
		if request.TimeIndex == 0 { //起始日期最新
			var table []*kline.KInfo

			lindex := len(dlines.List)
			for i := lindex - int(request.Num); i < lindex; i++ {
				table = append(table, dlines.List[i])
			}
			reply = &kline.ReplyHisK{
				Code: 200,
				Data: &kline.HisK{
					SID:   request.SID,
					Type:  request.Type,
					Total: int32(len(dlines.List)),
					Begin: table[0].NTime,
					Num:   int32(len(table)),
					List:  table,
				},
			}

			if request.Direct == 1 {
				var sig []*kline.KInfo
				sig = append(sig, table[len(table)-1])
				return &kline.ReplyHisK{
					Code: 200,
					Data: &kline.HisK{
						SID:   request.SID,
						Type:  request.Type,
						Total: int32(len(dlines.List)),
						Begin: table[len(table)-1].NTime,
						Num:   1,
						List:  sig,
					},
				}
			}

		} else { //TimeIndex作为起始日期

			var frontedSwap, palinalSwap []*kline.KInfo
			var databuf []*kline.KInfo

			for _, v := range dlines.List {
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
					var table []*kline.KInfo

					lindex := len(dlines.List)
					for i := lindex - int(request.Num); i < lindex; i++ {
						table = append(table, dlines.List[i])
					}

					var sig []*kline.KInfo
					sig = append(sig, table[len(table)-1])
					return &kline.ReplyHisK{
						Code: 200,
						Data: &kline.HisK{
							SID:   request.SID,
							Type:  request.Type,
							Total: int32(len(dlines.List)),
							Begin: table[len(table)-1].NTime,
							Num:   1,
							List:  sig,
						},
					}
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
				reply = &kline.ReplyHisK{
					Code: 4002,
					Data: &kline.HisK{},
				}

				logging.Error("Invalid request parameters 'Direct'...")
				return &kline.ReplyHisK{
					Code: 4002,
					Data: &kline.HisK{},
				}
			}

			reply = &kline.ReplyHisK{
				Code: 200,
				Data: &kline.HisK{
					SID:   request.SID,
					Type:  request.Type,
					Total: int32(len(dlines.List)),
					Begin: databuf[0].NTime,
					Num:   int32(len(databuf)),
					List:  databuf,
				},
			}

		}
	}
	return reply
}
