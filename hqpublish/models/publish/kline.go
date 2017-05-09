//K线
package publish

import (
	"fmt"

	. "haina.com/share/models"

	"ProtocolBuffer/format/kline"

	"github.com/golang/protobuf/proto"

	redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

var _ = fmt.Println
var _ = redigo.Bytes
var _ = logging.Info

type KLine struct {
	Model `db:"-"`
}

func NewKLine(rediskey string) *KLine {
	return &KLine{
		Model: Model{
			CacheKey: rediskey,
		},
	}
}

// 获取日K线
func (this *KLine) GetHisKLine(sid int32) (*kline.KInfoTable, int, error) {
	var lines kline.KInfoTable

	key := fmt.Sprintf(this.CacheKey, sid)

	dls, err := redis.Get(key)
	if err != nil {
		return nil, 0, err
	}
	if dls == "" {
		return nil, 0, ERROR_REDIS_LIST_NULL
	}

	data := []byte(dls)

	if err := proto.Unmarshal(data, &lines); err != nil {
		logging.Error("%v", err.Error())
		return nil, 0, err
	}

	return &lines, len(lines.List), nil
}
