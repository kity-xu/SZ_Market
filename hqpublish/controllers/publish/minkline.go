package publish

import (
	"encoding/json"
	"io"
	"net/http"

	"haina.com/share/lib"

	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"

	"haina.com/share/logging"

	"ProtocolBuffer/format/kline"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
)

type MinKLine struct{}

func NewMinKLine() *MinKLine {
	return &MinKLine{}
}

func (this *MinKLine) POST(c *gin.Context) {
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

func getRequestData(c *gin.Context) ([]byte, error) {
	temp := make([]byte, 1024)
	n, err := c.Request.Body.Read(temp)
	if err != nil && err != io.EOF {
		logging.Error("Body Read: %v", err)
		return nil, err
	}
	//logging.Info("\nBody len %d\n%s", n, temp[:n])
	logging.Info("Body len %d", n)
	return temp[:n], nil

}

func (this *MinKLine) PostJson(c *gin.Context) {
	buf, err := getRequestData(c)
	if err != nil && err != io.EOF {
		logging.Error("%v", err)
		return
	}

	var request kline.RequestMinK
	if err := json.Unmarshal(buf, &request); err != nil {
		logging.Error("Json Request Unmarshal: %v", err)
		return
	}
	logging.Info("Request Data: %+v", request)
	reply, err := publish.NewMinKLine().GetMinKLine(&request)
	if err != nil {
		logging.Error("%v", err)
		reply := kline.ReplyMinK{
			Code: 40002,
		}
		c.JSON(http.StatusOK, reply)
		return
	}

	c.JSON(http.StatusOK, reply)
}

func (this *MinKLine) PostPB(c *gin.Context) {
	var (
		replypb []byte
		err     error
		request kline.RequestMinK
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
	reply, err := publish.NewMinKLine().GetMinKLine(&request)
	if err != nil {
		reply := kline.ReplyMinK{
			Code: 40002,
		}
		replypb, err = proto.Marshal(&reply)
		if err != nil {
			logging.Error("pb marshal error: %v", err)
		}
	}
	replypb, err = proto.Marshal(reply)
	if err != nil {
		reply := kline.ReplyMinK{
			Code: 40002,
		}
		replypb, err = proto.Marshal(&reply)
		if err != nil {
			logging.Error("pb marshal error: %v", err)
		}

	}
	lib.WriteData(c, replypb)
}
