package publish

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	. "haina.com/market/hqpublish/models"
	. "haina.com/share/models"

	hsgrr "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
)

type BlockShotM struct {
	Model `db:"-"`
}

func NewBlockShotM() *BlockShotM {
	return &BlockShotM{
		Model: Model{
			CacheKey: REDISKEY_STOCK_BLOCK_SHOT,
		},
	}
}

// 查询板块快照
func (this *BlockShotM) GetBlockShotM(req *protocol.RequestBlockShot) (*protocol.PayloadBlockShot, error) {

	key := fmt.Sprintf(this.CacheKey, req.BoardCode, req.KeyCode)

	bin, err := RedisStore.GetBytes(key)
	if err != nil {
		if err == hsgrr.ErrNil {
			logging.Warning("redis not found key: %v", key)
			return nil, err
		}
		return nil, err
	}

	data := &protocol.StockBlockShotBase{}
	buffer := bytes.NewBuffer(bin)
	if err := binary.Read(buffer, binary.LittleEndian, data); err != nil && err != io.EOF {
		logging.Error("binary decode error: %v", err)
		return nil, err
	}
	var pbs protocol.PayloadBlockShot
	pbs.SKBShot = data
	return &pbs, err
}
