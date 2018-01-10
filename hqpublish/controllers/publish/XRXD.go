package publish

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"

	. "haina.com/market/hqpublish/controllers"
	"haina.com/market/hqpublish/controllers/publish/kline"

	"haina.com/share/lib"
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
)

var (
	_ = fmt.Print
	_ = lib.WriteString
	_ = logging.Error
	_ = protocol.MarketStatus{}
	_ = publish.MarketStatus{}
	_ = strings.ToLower
	_ = proto.Marshal
	_ = json.Marshal
	_ = io.ReadFull
)

type XRXD struct{}

func NewXRXD() *XRXD {
	return &XRXD{}
}

type Null struct {
	SID   int32
	Type  int32
	Total int
	Num   int
}

func (this *XRXD) POST(c *gin.Context) {
	replayfmt := c.Query(models.CONTEXT_FORMAT)
	if len(replayfmt) == 0 {
		replayfmt = "pb" // 默认
	}

	switch replayfmt {
	case "json":
		this.PostJson(c)
	case "pb":
		this.PostPB(c)
	default:
		return
	}
}

func (this *XRXD) PostJson(c *gin.Context) {
	var req protocol.RequestXRXD
	code, err := RecvAndUnmarshalJson(c, 1024, &req)
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, code, nil)
		return
	}
	logging.Info("Request %+v", req)
	if req.SID == 0 || req.TimeIndex < 0 || req.Num < 0 {
		WriteJson(c, 40004, nil)
		return
	}

	if req.Type < 10 {
		klines, err := publish.NewXRXD().GetXRDAllKlines(&req)
		if err != nil {
			if err == models.DATA_ISNULL {
				res := &Null{
					SID:   req.SID,
					Type:  req.Type,
					Total: -1,
					Num:   -1,
				}
				WriteJson(c, 200, res)
				return
			}
			logging.Error("%v", err)
			WriteJson(c, 40002, nil)
			return
		}
		payload, err := CreateTypeKline(klines, &req)
		if err != nil {
			logging.Error("%v", err)
			WriteJson(c, 40002, nil)
			return
		}

		WriteJson(c, 200, payload)
		return
	} else {
		hmin := kline.NewKline()
		request := &protocol.RequestHisK{
			SID:       req.SID,
			Type:      req.Type,
			TimeIndex: req.TimeIndex,
			Num:       req.Num,
			Direct:    req.Direct,
		}

		switch protocol.HAINA_KLINE_TYPE(request.Type) {

		case protocol.HAINA_KLINE_TYPE_KMIN1:
			hmin.MinJson_01(c, request)

		case protocol.HAINA_KLINE_TYPE_KMIN5:
			hmin.MinJson_05(c, request)

		case protocol.HAINA_KLINE_TYPE_KMIN15:
			hmin.MinJson_15(c, request)

		case protocol.HAINA_KLINE_TYPE_KMIN30:
			hmin.MinJson_30(c, request)

		case protocol.HAINA_KLINE_TYPE_KMIN60:
			hmin.MinJson_60(c, request)

		default:
			logging.Error("Invalid parameter 'Type'...")
		}
	}

}

func (this *XRXD) PostPB(c *gin.Context) {
	var req protocol.RequestXRXD
	code, err := RecvAndUnmarshalPB(c, 1024, &req)
	if err != nil {
		logging.Error("post pb %v", err)
		WriteDataErrCode(c, code)
		return
	}
	logging.Info("Request %+v", req)
	if req.SID == 0 || req.TimeIndex < 0 || req.Num < 0 {
		WriteDataErrCode(c, 40004)
		return
	}

	if req.Type < 10 {
		klines, err := publish.NewXRXD().GetXRDAllKlines(&req)
		if err != nil {
			if err == models.DATA_ISNULL {
				res := &protocol.PayloadXRXD{
					SID:   req.SID,
					Type:  req.Type,
					Total: -1,
					Num:   -1,
				}
				WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_XRXD, res)
				return
			}
			logging.Error("%v", err)
			WriteDataErrCode(c, 40002)
			return
		}
		payload, err := CreateTypeKline(klines, &req)
		if err != nil {
			logging.Error("%v", err)
			WriteJson(c, 40002, nil)
			return
		}
		WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_XRXD, payload)
		return
	} else {
		hmin := kline.NewKline()
		request := &protocol.RequestHisK{
			SID:       req.SID,
			Type:      req.Type,
			TimeIndex: req.TimeIndex,
			Num:       req.Num,
			Direct:    req.Direct,
		}

		switch protocol.HAINA_KLINE_TYPE(request.Type) {

		case protocol.HAINA_KLINE_TYPE_KMIN1:
			hmin.MinPB_01(c, request)

		case protocol.HAINA_KLINE_TYPE_KMIN5:
			hmin.MinPB_05(c, request)

		case protocol.HAINA_KLINE_TYPE_KMIN15:
			hmin.MinPB_15(c, request)

		case protocol.HAINA_KLINE_TYPE_KMIN30:
			hmin.MinPB_30(c, request)

		case protocol.HAINA_KLINE_TYPE_KMIN60:
			hmin.MinPB_60(c, request)

		default:
			logging.Error("Invalid parameter 'Type'...")
		}
	}
}

func CreateTypeKline(dlines *[]*protocol.KInfo, request *protocol.RequestXRXD) (*protocol.PayloadXRXD, error) {
	kline.InitMarketTradeDate()
	var e error
	var ret *protocol.PayloadXRXD
	//models.GetASCStruct(dlines) //升序排序

	switch request.Type {
	case 1:
		e = kline.MaybeAddKline(dlines, request.SID, e)
		if e != nil {
			logging.Error(e.Error())
			return nil, e
		}

	case 2:
		e = kline.MaybeAddWeekLine(dlines, request.SID, e)
		if e != nil {
			logging.Error(e.Error())
			return nil, e
		}

	case 3:
		e = kline.MaybeAddMonthLine(dlines, request.SID, e)
		if e != nil {
			logging.Error(e.Error())
			return nil, e
		}

	case 4:
		e = kline.MaybeAddYearLine(dlines, request.SID, e)
		if e != nil {
			logging.Error(e.Error())
			return nil, e
		}

	}

	if len(*dlines) == 0 {
		return nil, e
	}

	total := int32(len(*dlines))

	if request.Num > total {
		request.Num = total
	}

	if request.Num == 0 { //num==0, 获取全部
		ret = &protocol.PayloadXRXD{
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
			ret = &protocol.PayloadXRXD{ //向前
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
				ret = &protocol.PayloadXRXD{
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

					ret = &protocol.PayloadXRXD{
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
			ret = &protocol.PayloadXRXD{
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
