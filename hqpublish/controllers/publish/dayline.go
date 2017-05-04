//日线
package publish

import (
	"ProtocolBuffer/format/kline"

	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type DayLine struct{}

func NewDayLine() *DayLine {
	return &DayLine{}
}

func (this *DayLine) POST(c *gin.Context) {
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
		//lib.WriteString(c, 40004, nil)
		return
	}
}

func (this *DayLine) PostJson(c *gin.Context) {
	buf, err := getRequestData(c)
	if err != nil && err != io.EOF {
		logging.Error("%v", err)
		return
	}

	var request kline.RequestHisK
	if err := json.Unmarshal(buf, &request); err != nil {
		logging.Error("Json Request Unmarshal: %v", err)
		return
	}
	logging.Info("Request Data: %+v", request)

	dlines, length, err := publish.NewDayLine().GetHisDayLine(request.SID, request.Begin)
	if err != nil {
		logging.Error("%v", err)
		return
	}

	reply := &kline.ReplyHisK{
		Data: &kline.HisK{},
	}
	reply.Code = 200
	reply.Data.Begin = request.Begin
	reply.Data.Type = request.Type
	reply.Data.Num = int32(len(*dlines))

	if request.Order == 0 { //升序
		GetASCStruct(dlines)
		reply.Data.List = *dlines
	} else if request.Order == 1 { //降序
		GetSECStruct(dlines)
		reply.Data.List = *dlines
	}

	reply.Data.Total = int32(length)

	c.JSON(http.StatusOK, reply)
}

func (this *DayLine) PostPB(c *gin.Context) {
	var (
		replypb []byte
		err     error
		request kline.RequestHisK
	)

	buf, err := getRequestData(c)
	if err != nil && err != io.EOF {
		logging.Error("%v", err)
		return
	}
	if err := proto.Unmarshal(buf, &request); err != nil {
		logging.Error("PB Request Unmarshal: %v", err)
		return
	}
	logging.Info("Request Data: %+v", request)

	dlines, length, err := publish.NewDayLine().GetHisDayLine(request.SID, request.Begin)
	if err != nil {
		reply := kline.ReplyHisK{
			Code: 40002,
		}
		replypb, err = proto.Marshal(&reply)
		if err != nil {
			logging.Error("pb marshal error: %v", err)
		}
	}

	reply := &kline.ReplyHisK{
		Data: &kline.HisK{},
	}
	reply.Code = 200
	reply.Data.Begin = request.Begin
	reply.Data.Type = request.Type
	reply.Data.Num = int32(len(*dlines))
	if request.Order == 0 { //升序
		GetASCStruct(dlines)
		reply.Data.List = *dlines
	} else if request.Order == 1 { //降序
		GetSECStruct(dlines)
		reply.Data.List = *dlines
	}
	reply.Data.Total = int32(length)

	//转PB
	replypb, err = proto.Marshal(reply)
	if err != nil {
		reply := kline.ReplyHisK{
			Code: 40002,
		}
		replypb, err = proto.Marshal(&reply)
		if err != nil {
			logging.Error("pb marshal error: %v", err)
		}

	}
	lib.WriteData(c, replypb)
}
