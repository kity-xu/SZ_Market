package publish

import (
	"github.com/gin-gonic/gin"
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

// 获取日K先历史数据
func (this *DayKLine) GET(c *gin.Context) {
	snid := c.Query(models.CONTEXT_SNID)
	key := "hgs:hq:" + snid + ":hkday"
	klpbinfo, err := redigo.Bytes(redis.Get(key))
	if err != nil {
		// 没找到
		if err == redigo.ErrNil {
			logging.Info("没找到数据 %v", err)
			return
		}
		// 其他错误
		return
	}

	lib.WriteData(c, klpbinfo)
}
