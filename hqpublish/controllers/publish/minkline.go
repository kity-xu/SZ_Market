package publish

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	//"haina.com/share/lib"

	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"

	"haina.com/share/logging"

	"ProtocolBuffer/format/redis/pbdef/kline"

	"github.com/gin-gonic/gin"
	//"github.com/golang/protobuf/proto"
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

func (this *MinKLine) PostJson(c *gin.Context) {
	temp := make([]byte, 1024)
	n, err := c.Request.Body.Read(temp)
	if err != nil && err != io.EOF {
		logging.Error("Body Read: %v", err)
		return
	}
	buf := temp[:n]
	logging.Info("%d %s", n, buf)

	var request kline.RequestMinK
	err = json.Unmarshal(buf, &request)
	if err != nil {
		logging.Error("Json Request Unmarshal: %v", err)
		return
	}

	fmt.Printf("%+v\n", request)

	//lib.WriteString(c, 200, "ok")
	reply, err := publish.NewMinKLine().GetMinKLine(&request)
	if err != nil {
		logging.Error("%v", err)
		reply := kline.ReplyMinK{
			Code: 40002,
		}
		//lib.WriteString(c, 40002, nil)
		c.JSON(http.StatusOK, reply)
		return
	}
	fmt.Printf("%+v\n", reply)

	c.JSON(http.StatusOK, reply)
}

func (this *MinKLine) PostPB(c *gin.Context) {
	var (
		//replypb []byte
		err error
	)

	temp := make([]byte, 1024)
	n, err := c.Request.Body.Read(temp)
	if err != nil && err != io.EOF {
		logging.Error("%v", err)
		return
	}
	buf := temp[:n]
	logging.Info("%d %s", n, buf)

	return

	//	replypb, err = publish.NewMinKLine().GetMinKLineReplyBytes()
	//	if err != nil {
	//		reply := securitytable.ReplySecurityCodeTable{
	//			Code: 40002,
	//		}
	//		replypb, err = proto.Marshal(&reply)
	//		if err != nil {
	//			logging.Error("pb marshal error: %v", err)
	//		}
	//	}
	//	lib.WriteData(c, replypb)
}
