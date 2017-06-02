package redistore

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"haina.com/market/hqpost/models"
	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpost/go/protocol"

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
func (this *MinKLine) GetMinKLineToday(sid int32) (*[]*protocol.KInfo, error) {
	key := fmt.Sprintf(this.CacheKey, sid)

	ls, err := redis.LRange(key, 0, -1)
	if err != nil {
		return nil, err
	}
	if len(ls) == 0 {
		return nil, models.ERROR_REDIS_LIST_NULL

	}

	kls := make([]*protocol.KInfo, 0, 241)
	for _, v := range ls {
		k := &protocol.KInfo{}
		buffer := bytes.NewBuffer([]byte(v))
		if err := binary.Read(buffer, binary.LittleEndian, k); err != nil && err != io.EOF {
			return nil, err
		}
		kls = append(kls, k)
	}

	return &kls, nil
}
