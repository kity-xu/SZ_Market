package kline

import (
	. "haina.com/market/hqpublish/controllers"
	"haina.com/market/hqpublish/models"
	"haina.com/share/logging"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"github.com/gin-gonic/gin"
)

type Kline struct {
}

func NewKline() *Kline {
	return &Kline{}
}

func (this *Kline) POST(c *gin.Context) {
	initMarketTradeDate()
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

func (this *Kline) PostJson(c *gin.Context) {
	var request = &protocol.RequestHisK{}

	code, err := RecvAndUnmarshalJson(c, 1024, request)
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, code, nil)
		return
	}
	logging.Info("Request Data: %+v", request)

	switch protocol.HAINA_KLINE_TYPE(request.Type) {
	case protocol.HAINA_KLINE_TYPE_KDAY:
		this.DayJson(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KWEEK:
		this.WeekJson(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KMONTH:
		this.MonthJson(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KYEAR:
		this.YearJson(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KMIN1:
		this.MinJson_01(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KMIN5:
		this.MinJson_05(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KMIN15:
		this.MinJson_15(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KMIN30:
		this.MinJson_30(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KMIN60:
		this.MinJson_60(c, request)
		break
	default:
		logging.Error("Invalid parameter 'Type'...")
	}

}

func (this *Kline) PostPB(c *gin.Context) {
	var request = &protocol.RequestHisK{}

	code, err := RecvAndUnmarshalPB(c, 1024, request)
	if err != nil {
		logging.Error("post pb %v", err)
		WriteJson(c, code, nil)
		return
	}
	logging.Info("Request Data: %+v", request)

	switch protocol.HAINA_KLINE_TYPE(request.Type) {
	case protocol.HAINA_KLINE_TYPE_KDAY:
		this.DayPB(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KWEEK:
		this.WeekPB(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KMONTH:
		this.MonthPB(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KYEAR:
		this.YearPB(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KMIN1:
		this.MinPB_01(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KMIN5:
		this.MinPB_05(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KMIN15:
		this.MinPB_15(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KMIN30:
		this.MinPB_30(c, request)
		break
	case protocol.HAINA_KLINE_TYPE_KMIN60:
		this.MinPB_60(c, request)
		break
	default:
		logging.Error("Invalid parameter 'Type'...")

	}

}
