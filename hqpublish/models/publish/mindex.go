package publish

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	ctrl "haina.com/market/hqpublish/controllers"
	. "haina.com/market/hqpublish/models"
	. "haina.com/share/models"

	pro "ProtocolBuffer/projects/hqpublish/go/protocol"

	hsgrr "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

var (
	_ = redis.Init
	_ = GetCache
	_ = ctrl.MakeRespDataByBytes
	_ = errors.New
	_ = fmt.Println
	_ = hsgrr.Bytes
	_ = logging.Info
	_ = bytes.NewBuffer
	_ = binary.Read
	_ = io.ReadFull
)

type MIndex struct {
	Model `db:"-"`
}

func NewMIndex() *MIndex {
	return &MIndex{
		Model: Model{
		//CacheKey: REDISKEY_SECURITY_TRADE,
		},
	}
}

func (this MIndex) GetMIndexObj() (*pro.PayloadMIndex, error) {
	return nil, nil
}
