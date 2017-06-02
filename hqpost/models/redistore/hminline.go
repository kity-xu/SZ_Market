package redistore

import (
	"fmt"

	"haina.com/market/hqpost/models"
	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpost/go/protocol"

	"github.com/golang/protobuf/proto"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

type HMinKLine struct {
	Ktype string
	Model `db:"-"`
}

func NewHMinKLine(key string) *HMinKLine {
	return &HMinKLine{
		Ktype: key,
		Model: Model{
			CacheKey: key,
		},
	}
}

func (this *HMinKLine) LPushHMinKLine(sid int32, kinfo *protocol.HMinLineDay) error {
	key := fmt.Sprintf(this.CacheKey, sid)
	data, err := proto.Marshal(kinfo)
	if err != nil {
		logging.Error("%v", err.Error())
		return err
	}
	if err = redis.Lpush(key, data); err != nil {
		logging.Error("%v", err.Error())
		return err
	}
	return nil
}

// 获取历史分钟线
func (this *HMinKLine) GetHMinKLine(sid int32) (*protocol.HMinTable, error) {
	var lines protocol.HMinTable

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
