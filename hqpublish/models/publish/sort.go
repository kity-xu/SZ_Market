//排序
package publish

import (
	"fmt"

	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/share/logging"

	//"github.com/golang/protobuf/proto"
	. "haina.com/market/hqpublish/models"
)

type Sort struct {
	Model `db:"-"`
}

func NewSort(redis_key string) *Sort {
	return &Sort{
		Model: Model{
			CacheKey: redis_key,
		},
	}
}

func (this *Sort) GetSortByFieldID(req *protocol.RequestSort) (*protocol.PayloadSort, error) {
	key := fmt.Sprintf(this.CacheKey, req.SetID, req.FieldID)

	bdata, err := RedisStore.GetBytes(key)
	if err != nil {
		return nil, err
	}

	step := 4

	for i := 0; i < len(bdata); i += step {
		bs := bdata[i : i+step]
		ret := uint32(bs[0]) | uint32(bs[1])<<8 | uint32(bs[2])<<16 | uint32(bs[3])<<24 //小端
		// binary.LittleEndian.Uint32(bs)
		logging.Debug("data:%v", ret)
	}
	return nil, nil
}
