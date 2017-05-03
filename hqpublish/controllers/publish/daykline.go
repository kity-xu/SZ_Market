package publish

import (
	"ProtocolBuffer/format/kline"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"haina.com/market/hqpublish/models"
	redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/lib"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

type DayKLine struct{}

func NewDayKLine() *DayKLine {
	return &DayKLine{}
}

// 获取历史日K线历史数据
func (this *DayKLine) GET(c *gin.Context) {
	snid := c.Query(models.CONTEXT_SNID)
	key := "hq:st:hday:" + snid
	var kli kline.ReplyKInfoTable
	klpbinfo, err := redigo.Bytes(redis.Get(key))
	if err != nil {
		// 没找到
		if err == redigo.ErrNil {
			logging.Info("没找到数据 %v", err)
			kli.Code = 40002
			v, err := proto.Marshal(&kli)
			if err != nil {
				logging.Info("%v", err)
			}
			lib.WriteData(c, v)
		}
	}

	lib.WriteData(c, klpbinfo)
}
