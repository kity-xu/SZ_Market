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

type StockBlockShotM struct {
	NBlockID        int32                   // 板块ID
	NTypeID         int32                   // 板块所属集合ID 行业 地区 概念
	NTime           int32                   // 时间 hhmmss
	NLastPx         uint32                  // 板块指数(*10000)
	NPreClosePx     uint32                  // 板块指数昨收(*10000)
	NOpenPx         uint32                  // 开盘价(*10000)
	NHighPx         uint32                  // 最高价(*10000)
	NLowPx          uint32                  // 最低价(*10000)
	NPxChg          int32                   // 涨跌(*10000)
	NPxAmplitude    int32                   // 振幅(*10000)
	NAveChgRate     int32                   // 平均涨跌幅(*10000)
	LlVolume        int64                   ///< 板块总成交量
	LlValue         int64                   ///< 板块总成交额(*10000)
	NStockID        int32                   // 领涨股
	NStockChgRate   int32                   // 领涨股涨跌幅 (*10000)
	NNum            int32                   // 成分股票个数
	NLongNum        int32                   // 上涨家数
	NShortNum       int32                   // 下跌家数
	NChgRatio       int32                   // 上涨比例 (*10000)
	LlValueOfInFlow int64                   ///< 资金净流入额(*10000)
	SzBlockName     [12]byte                //板块名称
	SzSName         [SECURITY_NAME_LEN]byte //领涨股名称
}

// 查询板块快照
func (this *BlockShotM) GetBlockShotM(req *protocol.RequestBlockShot) (*protocol.PayloadBlockShot, error) {

	key := fmt.Sprintf(this.CacheKey, req.BlockID)

	bin, err := RedisStore.GetBytes(key)
	if err != nil {
		if err == hsgrr.ErrNil {
			logging.Warning("redis not found key: %v", key)
			return nil, err
		}
		return nil, err
	}

	data := &protocol.StockBlockShotBase{}
	sbsm := &StockBlockShotM{}
	buffer := bytes.NewBuffer(bin)
	if err := binary.Read(buffer, binary.LittleEndian, sbsm); err != nil && err != io.EOF {
		logging.Error("binary decode error: %v", err)
		return nil, err
	}
	data.NBlockID = sbsm.NBlockID
	data.NTypeID = sbsm.NTypeID
	data.NLastPx = sbsm.NLastPx
	data.NPreClosePx = sbsm.NPreClosePx
	data.NAveChgRate = sbsm.NAveChgRate
	data.LlVolume = sbsm.LlVolume
	data.LlValue = sbsm.LlValue
	data.NStockID = sbsm.NStockID
	data.NStockChgRate = sbsm.NStockChgRate
	data.NNum = sbsm.NNum
	data.NLongNum = sbsm.NLongNum
	data.NShortNum = sbsm.NShortNum
	data.NChgRatio = sbsm.NChgRatio
	data.LlValueOfInFlow = sbsm.LlValueOfInFlow
	data.SzBlockName = byte12ToString(sbsm.SzBlockName)
	data.SzSName = byte40ToString(sbsm.SzSName)
	data.NTime = sbsm.NTime
	data.NOpenPx = sbsm.NOpenPx
	data.NHighPx = sbsm.NHighPx
	data.NLowPx = sbsm.NLowPx
	data.NPxChg = sbsm.NPxChg
	data.NPxAmplitude = sbsm.NPxAmplitude

	var pbs protocol.PayloadBlockShot
	pbs.BlockShort = data
	return &pbs, err
}
