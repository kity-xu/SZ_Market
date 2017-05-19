package publish

import (
	"fmt"

	. "haina.com/share/models"

	"ProtocolBuffer/format/kline"

	"github.com/golang/protobuf/proto"

	//redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

type HMinKLine struct {
	Model `db:"-"`
}

func NewHMinKLine(redis_key string) *HMinKLine {
	return &HMinKLine{
		Model: Model{
			CacheKey: redis_key,
		},
	}
}

// 获取历史分钟线
func (this *HMinKLine) GetHMinKLine(sid int32) (*kline.HMinTable, error) {
	var lines kline.HMinTable

	key := fmt.Sprintf(this.CacheKey, sid)

	dls, err := redis.Get(key)
	if err != nil {
		return nil, err
	}
	if dls == "" {
		logging.Info("Couldn't find the historical data from redis...May be generated for the first time")
		return nil, ERROR_REDIS_LIST_NULL
	}

	data := []byte(dls)

	if err := proto.Unmarshal(data, &lines); err != nil {
		logging.Error("%v", err.Error())
		return nil, err
	}

	return &lines, nil
}
