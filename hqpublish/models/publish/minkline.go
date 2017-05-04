package publish

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	. "haina.com/share/models"

	"ProtocolBuffer/format/kline"

	redigo "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

var _ = fmt.Println
var _ = redigo.Bytes
var _ = logging.Info
var _ = bytes.NewBuffer
var _ = binary.Read
var _ = io.ReadFull

type MinKLine struct {
	Model `db:"-"`
}

func NewMinKLine() *MinKLine {
	return &MinKLine{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_MIN,
		},
	}
}

// 获取分钟K线
func (this *MinKLine) GetMinKLine(request *kline.RequestMinK) (*kline.MinK, error) {
	key := fmt.Sprintf(this.CacheKey, request.SID)

	ls, err := redis.LRange(key, 0, -1)
	if err != nil {
		return nil, err
	}
	if len(ls) == 0 {
		return nil, ERROR_REDIS_LIST_NULL
	}

	if request.BeginTime > 150100 {
		return nil, ERROR_KLINE_BEGIN_TIME
	}

	kls := make([]*kline.KInfo, 0, 241)
	for _, v := range ls {
		k := &kline.KInfo{}
		buffer := bytes.NewBuffer([]byte(v))
		if err := binary.Read(buffer, binary.LittleEndian, k); err != nil && err != io.EOF {
			return nil, err
		}
		if k.NTime >= request.BeginTime {
			kls = append(kls, k)
		}
	}

	ret := &kline.MinK{
		SID:       request.SID,
		BeginTime: request.BeginTime,
		Num:       int32(len(kls)),
		List:      kls,
	}

	return ret, nil
}
