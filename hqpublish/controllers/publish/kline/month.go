//月线
package kline

import (
	"ProtocolBuffer/format/kline"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"

	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

func (this *Kline) MonthJson(c *gin.Context, request *kline.RequestHisK) {
	reply := this.ReplyKLine(c, publish.REDISKEY_SECURITY_HMONTH, request)

	c.JSON(http.StatusOK, reply)

}

func (this *Kline) MonthPB(c *gin.Context, request *kline.RequestHisK) {
	reply := this.ReplyKLine(c, publish.REDISKEY_SECURITY_HMONTH, request)

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
