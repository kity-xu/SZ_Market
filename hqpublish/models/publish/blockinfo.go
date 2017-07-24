// 证券集合（板块）
package publish

import (
	"strconv"

	"ProtocolBuffer/projects/hqpublish/go/protocol"
	"fmt"

	. "haina.com/market/hqpublish/models"
	. "haina.com/share/models"

	"github.com/golang/protobuf/proto"
	hsgrr "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
)

type BlockInfoM struct {
	Model `db:"-"`
}

func NewBlockInfoM() *BlockInfoM {
	return &BlockInfoM{
		Model: Model{
			CacheKey: REDISKEY_STOCK_BLOCK_INFO,
		},
	}
}

// 根据成分股查询所属板块
func (this *BlockInfoM) GetBlockInfoBySID(req *protocol.Requeststockinfo) (*protocol.Payloadstockinfo, error) {

	key1 := fmt.Sprintf(this.CacheKey, 1102)
	ski1, err := this.GetPayloadStockInfo(key1, req)
	if err != nil {
		logging.Error("block 1102 cast error:%v", err)
	}
	key2 := fmt.Sprintf(this.CacheKey, 1105)
	ski2, err := this.GetPayloadStockInfo(key2, req)
	if err != nil {
		logging.Error("block 1105 cast error:%v", err)
	}
	ski1.SKIList = append(ski1.SKIList, ski2.SKIList...)
	key3 := fmt.Sprintf(this.CacheKey, 1109)
	ski3, err := this.GetPayloadStockInfo(key3, req)
	if err != nil {
		logging.Error("block 1109 cast error:%v", err)
	}
	ski1.SKIList = append(ski1.SKIList, ski3.SKIList...)

	return ski1, err
}

func (this *BlockInfoM) GetPayloadStockInfo(key string, req *protocol.Requeststockinfo) (*protocol.Payloadstockinfo, error) {
	ls, err := hsgrr.Strings(RedisStore.Do("keys", key))
	if err != nil {
		logging.Info("error:%v", err)
	}
	var psk protocol.Payloadstockinfo
	for _, v := range ls {
		bin, err := RedisStore.GetBytes(v)
		if err != nil {
			if err == hsgrr.ErrNil {
				logging.Warning("redis not found key: %v", key)
				return nil, err
			}
			return nil, err
		}

		var elms = protocol.ElementList{}
		if err = proto.Unmarshal(bin, &elms); err != nil {
			logging.Debug("%v", err.Error())
			return nil, err
		}

		var sbp protocol.StockBlockinfoP
		for _, ite := range elms.List {
			//判断是属于此板块
			if req.SID == ite.NSid {
				sid, err := strconv.Atoi(v[16:])
				bocod, err := strconv.Atoi(v[11:15])
				if err != nil {
					logging.Error("v type cast error:%v", err)
					return nil, err
				}
				sbp.BoardCode = int32(bocod)
				sbp.SetID = int32(sid)
				psk.SKIList = append(psk.SKIList, &sbp)
				break
			}
		}
	}
	return &psk, err
}
