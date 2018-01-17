package redistore

import (
	"ProtocolBuffer/projects/hqpost/go/protocol"
	"bytes"
	"encoding/binary"
	"io"

	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

// TradeDate 交易日
var (
	SH_TradeDate int32 //上海交易所
	SZ_TradeDate int32 //深证交易所
)

type QUOTE_RECORD struct {
	NPx      uint32 ///< 委托价格(*10000)
	LlVolume int64  ///< 委托数量
}

func init() {
	SH_TradeDate = TradeDateByMarketStatus(100000000)
	SZ_TradeDate = TradeDateByMarketStatus(200000000)
}

/// 股票快照消息 (MSG_CALC_SNAPSHOT REDISKEY_SECURITY_SNAP )
/// 数据 Redis 里存的是二进制
/// 由于 protobuf 数组模型无法定义确定大小，特在此写出 Redis 里对应的计算端存的数据定义
type REDIS_BIN_STOCK_SNAPSHOT struct {
	NSID          int32           ///< 证券ID
	NTime         int32           ///< 时间 unix time
	NTradingPhase uint32          ///< 0:开市前  1:开盘集合竞价 2:连续竞价 3:临时停盘 4:收盘集合竞价 5:集中竞价闭市  6:协议转让结束  7:闭市
	NPreClosePx   uint32          ///< 昨收价 * 10000
	NOpenPx       uint32          ///< 开盘价 ..
	NHighPx       uint32          ///< 最高价 ..
	NLowPx        uint32          ///< 最低价 ..
	NLastPx       uint32          ///< 最新价 ..
	NHighLimitPx  uint32          ///< 涨停价格(*10000)
	NLowLimitPx   uint32          ///< 跌停价格(*10000)
	LlTradeNum    int64           ///< 成交笔数
	LlVolume      int64           ///< 成交量
	LlValue       int64           ///< 成交额(*10000)
	NQuoteSize    int32           ///< 报价总档数
	NWABidPx      uint32          ///< 加权平均委买均价(*10000)
	NWAOfferPx    uint32          ///< 加权平均委卖均价(*10000)
	LlToBidVol    int64           ///< 总委买量
	LlToOfferVol  int64           ///< 总委卖量
	LlInnerVolume int64           ///< 内盘成交量
	LlOuterVolume int64           ///< 外盘成交量
	LlInnerValue  int64           ///< 内盘成交额
	LlOuterValue  int64           ///< 外盘成交额
	NPxChg        int32           ///< 涨跌
	PxChgRatio    int32           ///< 涨跌幅*10000
	NPxAmplitude  int32           ///< 振幅*10000
	NLiangbi      int32           ///< 量比*100
	NWeibi        int32           ///< 委比*10000
	NTurnOver     int32           ///< 换手率*10000
	NPE           int32           ///< 动态市盈率*10000
	NPB           int32           ///< 动态市净率*10000
	Bid           [5]QUOTE_RECORD ///< 买5档
	Offer         [5]QUOTE_RECORD ///< 卖5档
}

// 指数快照
type REDIS_BIN_INDEX_SNAPSHOT struct {
	NSID          int32
	NTime         int32
	NTradingPhase uint32
	NPreClosePx   uint32
	NOpenPx       uint32
	NHighPx       uint32
	NLowPx        uint32
	NLastPx       uint32
	NHighLimitPx  uint32
	NLowLimitPx   uint32
	LlTradeNum    int64
	LlVolume      int64
	LlValue       int64
	NPxChg        int32
	NPxChgRatio   int32
	NPxAmplitude  int32
	NLiangbi      int32
	NWeibi        int32
	NTurnOver     int32
	NPE           int32
	NPB           int32
	NLongNum      int32
	NShortNum     int32
	NBandNum      int32
}

// GetStockSnapshotObj 获取证券快照
func GetStockSnapshotObj(key string, sid int32) (*protocol.KInfo, error) {
	var ntime int32
	bin, err := redis.Get(key)
	if err != nil {
		logging.Debug("redis not found key: %v", key)
		return nil, nil
	}

	data := REDIS_BIN_STOCK_SNAPSHOT{}
	buffer := bytes.NewBuffer([]byte(bin))
	if err := binary.Read(buffer, binary.LittleEndian, &data); err != nil && err != io.EOF { //个股
		buf := bytes.NewBuffer([]byte(bin))
		index := &REDIS_BIN_INDEX_SNAPSHOT{}
		if err = binary.Read(buf, binary.LittleEndian, index); err != nil && err != io.EOF { //指数
			logging.Error("binary decode error: %v", err)
			return nil, err
		}

		if index.LlVolume == 0 { // 停盘
			logging.Debug("index sid:%v delist----TradingPhase:%c", sid, index.NTradingPhase)
			return nil, nil
		}

		var avgpx uint32
		if index.LlVolume == 0 {
			avgpx = 0
		} else {
			avgpx = uint32(index.LlValue / index.LlVolume)
		}

		if index.NSID/1000000 == 100 {
			ntime = SH_TradeDate
		} else {
			ntime = SZ_TradeDate
		}
		ret := &protocol.KInfo{
			NSID:     index.NSID,
			NTime:    ntime,
			NPreCPx:  int32(index.NPreClosePx),
			NOpenPx:  int32(index.NOpenPx),
			NHighPx:  int32(index.NHighPx),
			NLowPx:   int32(index.NLowPx),
			NLastPx:  int32(index.NLastPx),
			LlVolume: index.LlVolume,
			LlValue:  index.LlValue,
			NAvgPx:   avgpx,
		}
		return ret, nil
	}
	if data.LlVolume == 0 { //停盘
		logging.Debug("stock sid:%v delist----TradingPhase:%c", sid, data.NTradingPhase)
		return nil, nil
	}

	var avgpx uint32
	if data.LlVolume == 0 {
		avgpx = 0
	} else {
		avgpx = uint32(data.LlValue / data.LlVolume)
	}
	if data.NSID/1000000 == 100 {
		ntime = SH_TradeDate
	} else {
		ntime = SZ_TradeDate
	}
	ret := &protocol.KInfo{
		NSID:     data.NSID,
		NTime:    ntime,
		NPreCPx:  int32(data.NPreClosePx),
		NOpenPx:  int32(data.NOpenPx),
		NHighPx:  int32(data.NHighPx),
		NLowPx:   int32(data.NLowPx),
		NLastPx:  int32(data.NLastPx),
		LlVolume: data.LlVolume,
		LlValue:  data.LlValue,
		NAvgPx:   avgpx,
	}
	return ret, nil
}
