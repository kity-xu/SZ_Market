// 证券快照
package publish

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	. "haina.com/market/hqpublish/models"
	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	hsgrr "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
)

var _ = fmt.Println
var _ = hsgrr.Bytes
var _ = logging.Info
var _ = bytes.NewBuffer
var _ = binary.Read
var _ = io.ReadFull

type IndexSnapshot struct {
	Model `db:"-"`
}

/// 指数快照消息
/// 数据 Redis 里存的是二进制

type REDIS_BIN_INDEX_SNAPSHOT struct {
}

func NewIndexSnapshot() *IndexSnapshot {
	return &IndexSnapshot{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_SNAP,
		},
	}
}

// 获取指数快照
func (this *IndexSnapshot) GetIndexSnapshotObj(req *protocol.RequestSnapshot) (*protocol.IndexSnapshot, error) {
	key := fmt.Sprintf(this.CacheKey, req.SID)

	bin, err := RedisStore.GetBytes(key)
	if err != nil {
		if err == hsgrr.ErrNil {
			logging.Warning("redis not found key: %v", key)
			return nil, err
		}
		return nil, err
	}

	data := &protocol.IndexSnapshot{}
	buffer := bytes.NewBuffer(bin)
	if err := binary.Read(buffer, binary.LittleEndian, data); err != nil && err != io.EOF {
		logging.Error("binary decode error: %v", err)
		return nil, err
	}

	return data, nil
}
