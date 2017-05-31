package redistore

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	//redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/market/hqpost/models"
	. "haina.com/share/models"

	"ProtocolBuffer/format/kline"

	"haina.com/share/store/redis"
)

type MinKLine struct {
	Model `db:"-"`
}

func NewMinKLine(key string) *MinKLine {
	return &MinKLine{
		Model: Model{
			CacheKey: key,
		},
	}
}

// 获取分钟K线
func (this *MinKLine) GetMinKLineToday(sid int32) (*[]*kline.KInfo, error) {
	key := fmt.Sprintf(this.CacheKey, sid)

	ls, err := redis.LRange(key, 0, -1)
	if err != nil {
		return nil, err
	}
	if len(ls) == 0 {
		return nil, models.ERROR_REDIS_LIST_NULL

	}

	kls := make([]*kline.KInfo, 0, 241)
	for _, v := range ls {
		k := &kline.KInfo{}
		buffer := bytes.NewBuffer([]byte(v))
		if err := binary.Read(buffer, binary.LittleEndian, k); err != nil && err != io.EOF {
			return nil, err
		}
		kls = append(kls, k)
	}

	return &kls, nil
}
