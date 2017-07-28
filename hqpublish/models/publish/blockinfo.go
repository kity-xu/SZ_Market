// 证券集合（板块）
package publish

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"
	"fmt"

	. "haina.com/market/hqpublish/models"
	. "haina.com/share/models"

	"github.com/golang/protobuf/proto"
	//hsgrr "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
)

type BlockInfoM struct {
	Model `db:"-"`
}

func NewBlockInfoM() *BlockInfoM {
	return &BlockInfoM{
		Model: Model{
			CacheKey: REDISKEY_STOCK_BLOCK_SID,
		},
	}
}

// 根据成分股查询所属板块
func (this *BlockInfoM) GetBlockInfoBySID(req *protocol.Requeststockinfo) (*protocol.Payloadstockinfo, error) {

	key := fmt.Sprintf(this.CacheKey, req.SID)

	slist, err := RedisStore.GetBytes(key)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, err
	}
	var blkl = &protocol.BlockList{}
	if err = proto.Unmarshal(slist, blkl); err != nil {
		logging.Debug("%v", err.Error())
		return nil, err
	}
	var psk protocol.Payloadstockinfo

	for _, ite := range blkl.List {
		// 根据板块id查询板块快照
		var reqb protocol.RequestBlockShot
		reqb.BlockID = ite.SetID
		block, err := NewBlockShotM().GetBlockShotM(&reqb)
		if err != nil {
			logging.Error("sel SnapShot error val:%v", err)
			return nil, err
		}
		psk.SnapShoot = append(psk.SnapShoot, block.BlockShort)
	}
	psk.Num = int32(len(blkl.List))
	return &psk, err
}
