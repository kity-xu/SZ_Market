package kline

import (
	"ProtocolBuffer/format/kline"
	"encoding/json"
	"io"

	"haina.com/market/hqpublish/models"
	"haina.com/share/logging"

	"github.com/golang/protobuf/proto"

	"github.com/gin-gonic/gin"
)

type Kline struct {
}

func NewKline() *Kline {
	return &Kline{}
}

func (this *Kline) POST(c *gin.Context) {
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
	var request kline.RequestHisK

	buf, err := models.GetRequestData(c)
	if err != nil && err != io.EOF {
		logging.Error("%v", err)
		return
	}

	if err := json.Unmarshal(buf, &request); err != nil {
		logging.Error("Json Request Unmarshal: %v", err)
		return
	}
	logging.Info("Request Data: %+v", request)

	switch kline.KLINE_TYPE(request.Type) {
	case kline.KLINE_TYPE_KLINE_TYPE_DAY:
		this.DayJson(c, &request)
		break
	case kline.KLINE_TYPE_KLINE_TYPE_WEEK:
		this.WeekJson(c, &request)
		break
	case kline.KLINE_TYPE_KLINE_TYPE_MONTH:
		this.MonthJson(c, &request)
		break
	case kline.KLINE_TYPE_KLINE_TYPE_YEAR:
		this.YearJson(c, &request)
		break
	default:
		logging.Error("Invalid parameter 'Type'...")
	}

}

func (this *Kline) PostPB(c *gin.Context) {
	var request kline.RequestHisK

	buf, err := models.GetRequestData(c)
	if err != nil && err != io.EOF {
		logging.Error("%v", err)
		return
	}
	if err := proto.Unmarshal(buf, &request); err != nil {
		logging.Error("PB Request Unmarshal: %v", err)
		return
	}
	logging.Info("Request Data: %+v", request)

	switch kline.KLINE_TYPE(request.Type) {
	case kline.KLINE_TYPE_KLINE_TYPE_DAY:
		this.DayPB(c, &request)
		break
	case kline.KLINE_TYPE_KLINE_TYPE_WEEK:
		this.WeekPB(c, &request)
		break
	case kline.KLINE_TYPE_KLINE_TYPE_MONTH:
		this.MonthPB(c, &request)
		break
	case kline.KLINE_TYPE_KLINE_TYPE_YEAR:
		this.YearPB(c, &request)
		break
	default:
		logging.Error("Invalid parameter 'Type'...")

	}

}
