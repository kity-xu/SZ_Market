package redis_minline

import (
	"fmt"

	"haina.com/market/hqpost/models"
	. "haina.com/share/models"

	"ProtocolBuffer/format/kline"

	"github.com/golang/protobuf/proto"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

type HMinKLine struct {
	Model `db:"-"`
}

func NewHMinKLine(key string) *HMinKLine {
	return &HMinKLine{
		Model: Model{
			CacheKey: key,
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
		return nil, models.ERROR_REDIS_LIST_NULL
	}

	data := []byte(dls)

	if err := proto.Unmarshal(data, &lines); err != nil {
		logging.Error("%v", err.Error())
		return nil, err
	}

	return &lines, nil
}

func WriteHMinLine(key string, data []byte) error {
	if err := redis.Set(key, data); err != nil {
		logging.Error("历史分钟线入redis出错...%v", err)
		return err
	}
	return nil
}
