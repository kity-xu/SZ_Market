// 证券快照
package publish

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	. "haina.com/share/models"

	"ProtocolBuffer/format/snap"

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

type StockSnapshot struct {
	Model `db:"-"`
}

/// 股票快照消息 (MSG_CALC_SNAPSHOT REDISKEY_SECURITY_SNAP )
type REDIS_BIN_STOCK_SNAPSHOT struct {
	NSID         int32           ///< 证券ID
	NTime        int32           ///< 时间 unix time
	NStatus      uint32          ///< 0:开市前  1:开盘集合竞价 2:连续竞价 3:临时停盘 4:收盘集合竞价 5:集中竞价闭市  6:协议转让结束  7:闭市
	NPreClosePx  uint32          ///< 昨收价 * 10000
	NOpenPx      uint32          ///< 开盘价 ..
	NHighPx      uint32          ///< 最高价 ..
	NLowPx       uint32          ///< 最低价 ..
	NLastPx      uint32          ///< 最新价 ..
	NHighLimitPx uint32          ///< 涨停价格(*10000)
	NLowLimitPx  uint32          ///< 跌停价格(*10000)
	LlTradeNum   int64           ///< 成交笔数
	LlVolume     int64           ///< 成交量
	LlValue      int64           ///< 成交额(*10000)
	NQuoteSize   int32           ///< 报价总档数
	NWABidPx     uint32          ///< 加权平均委买均价(*10000)
	NWAOfferPx   uint32          ///< 加权平均委卖均价(*10000)
	LlToBidVol   int64           ///< 总委买量
	LlToOfferVol int64           ///< 总委卖量
	NPxChg       int32           ///< 涨跌
	PxChgRatio   int32           ///< 涨跌幅*10000
	NPxAmplitude int32           ///< 振幅*10000
	NLiangbi     int32           ///< 量比*100
	NWeibi       int32           ///< 委比*10000
	NTurnOver    int32           ///< 换手率*10000
	NPE          int32           ///< 动态市盈率*10000
	NPB          int32           ///< 动态市净率*10000
	Bid          [5]QUOTE_RECORD ///< 买5档
	Offer        [5]QUOTE_RECORD ///< 卖5档
}
type QUOTE_RECORD struct {
	NPx      uint32 ///< 委托价格(*10000)
	LlVolume int64  ///< 委托数量
}

func NewStockSnapshot() *StockSnapshot {
	return &StockSnapshot{
		Model: Model{
			CacheKey: REDISKEY_SECURITY_SNAP,
		},
	}
}

// 获取证券快照
func (this *StockSnapshot) GetStockSnapshot(request *snap.RequestSnap) (*snap.Snapshot, error) {
	key := fmt.Sprintf(this.CacheKey, request.SID)

	bin, err := redigo.Bytes(redis.Get(key))
	if err != nil {
		if err == redigo.ErrNil {
			logging.Warning("redis not found key: %v", key)
			return nil, err
		}
		return nil, err
	}

	data := REDIS_BIN_STOCK_SNAPSHOT{}
	buffer := bytes.NewBuffer(bin)
	if err := binary.Read(buffer, binary.LittleEndian, &data); err != nil && err != io.EOF {
		logging.Fatal(err)
	}

	ret := &snap.Snapshot{
		SID: request.SID,
		Snap: &snap.StockSnapshot{
			NSID:         data.NSID,
			NTime:        data.NTime,
			NStatus:      data.NStatus,
			NPreClosePx:  data.NPreClosePx,
			NOpenPx:      data.NOpenPx,
			NHighPx:      data.NHighPx,
			NLowPx:       data.NLowPx,
			NLastPx:      data.NLastPx,
			NHighLimitPx: data.NHighLimitPx,
			NLowLimitPx:  data.NLowLimitPx,
			LlTradeNum:   data.LlTradeNum,
			LlVolume:     data.LlVolume,
			LlValue:      data.LlValue,
			NQuoteSize:   data.NQuoteSize,
			NWABidPx:     data.NWABidPx,
			NWAOfferPx:   data.NWAOfferPx,
			LlToBidVol:   data.LlToBidVol,
			LlToOfferVol: data.LlToOfferVol,
			NPxChg:       data.NPxChg,
			PxChgRatio:   data.PxChgRatio,
			NPxAmplitude: data.NPxAmplitude,
			NLiangbi:     data.NLiangbi,
			NWeibi:       data.NWeibi,
			NTurnOver:    data.NTurnOver,
			NPE:          data.NPE,
			NPB:          data.NPB,
			Bid:          make([]*snap.QuoteRecord, 0, 5),
			Offer:        make([]*snap.QuoteRecord, 0, 5),
		},
	}
	for _, v := range data.Bid {
		bid := &snap.QuoteRecord{
			NPx:      v.NPx,
			LlVolume: v.LlVolume,
		}
		ret.Snap.Bid = append(ret.Snap.Bid, bid)
	}
	for _, v := range data.Offer {
		offer := &snap.QuoteRecord{
			NPx:      v.NPx,
			LlVolume: v.LlVolume,
		}
		ret.Snap.Offer = append(ret.Snap.Offer, offer)
	}

	return ret, nil
}
