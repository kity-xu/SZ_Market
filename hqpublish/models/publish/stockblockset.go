// 证券集合（板块）
package publish

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	. "haina.com/market/hqpublish/models"
	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	hsgrr "haina.com/share/garyburd/redigo/redis"
	"haina.com/share/logging"
)

type StockBlockBase struct {
	Model `db:"-"`
}

/// 板块定义
/// 数据 Redis 里存的是二进制
type StockBlockInfo struct {
	NBlockID        int32                   // 板块ID
	SzBlockName     [12]byte                // 板块名称
	NAveChgRate     int32                   // 平均涨跌幅*10000
	LlVolume        int64                   ///< 板块总成交量
	LlValue         int64                   ///< 板块总成交额(*10000)
	NStockID        int32                   // 领涨股ID
	NStockChgRate   int32                   // 领涨股涨跌幅
	SzSName         [SECURITY_NAME_LEN]byte // 领涨股名称
	NNum            int32                   // 成分股票个数
	NChgRatio       int32                   //上涨比例
	NLongNum        int32                   // 上涨家数
	NShortNum       int32                   // 下跌家数
	LlValueOfInFlow int64                   ///< 资金流入额(*10000)
	NLastPx         int32                   // 板块指数(*10000)
}

func NewStockBlockBase() *StockBlockBase {
	return &StockBlockBase{
		Model: Model{
			CacheKey: REDISKEY_STOCK_BLOCK,
		},
	}
}

// 获取板块列表
func (this *StockBlockBase) GetStockBlockBase(req *protocol.RequestStockBlockBase) (*protocol.PayloadStockBlockSet, error) {

	var psb protocol.PayloadStockBlockBase

	switch req.SetID {
	case 1102:
		key := fmt.Sprintf(this.CacheKey, 1102)
		gpsb := GetPayloadStockBlock(key)
		psb.List = append(psb.List, gpsb.List...)
	case 1105:
		key := fmt.Sprintf(this.CacheKey, 1105)
		gpsb := GetPayloadStockBlock(key)
		psb.List = append(psb.List, gpsb.List...)
	case 1109:
		key := fmt.Sprintf(this.CacheKey, 1109)
		gpsb := GetPayloadStockBlock(key)
		psb.List = append(psb.List, gpsb.List...)
	case 1111:
		key2 := fmt.Sprintf(this.CacheKey, 1102)
		gpsb2 := GetPayloadStockBlock(key2)
		psb.List = append(psb.List, gpsb2.List...)
		key5 := fmt.Sprintf(this.CacheKey, 1105)
		gpsb5 := GetPayloadStockBlock(key5)
		psb.List = append(psb.List, gpsb5.List...)
		key9 := fmt.Sprintf(this.CacheKey, 1109)
		gpsb9 := GetPayloadStockBlock(key9)
		psb.List = append(psb.List, gpsb9.List...)
	}

	// 对比排序后返回结构体
	var psbs protocol.PayloadStockBlockSet

	// 排序key
	var sortk int32
	if req.SortID > 4008 || req.SortID < 4001 {
		sortk = 4006
	} else {
		sortk = req.SortID
	}
	Sortkey := fmt.Sprintf(REDISKEY_STOCK_BLOCK_BASE, sortk)
	// 排序规则 默认正序 101 正序 102 倒序
	var orderrule int32

	if orderrule == 102 {
		orderrule = int32(102)
	} else {
		if req.OrderRule == int32(101) {
			orderrule = int32(101)
		} else {
			orderrule = int32(101)
		}
	}
	sbinl, err := GetSortStockBlockL(Sortkey, orderrule)
	if err != nil {
		logging.Info("获取排序集合error %v", err)
	}
	// 板块排序集合
	for _, ite := range sbinl {
		// 板块集合
		for _, stb := range psb.List {
			if stb.NBlockID == ite.NBlockID {
				var tbhi protocol.TagBlockHSortInfo
				tbhi.NBlockID = stb.NBlockID
				tbhi.SzBlockName = strings.Replace(string(ite.SzBlockName[:]), "\u0000", "", -1)
				tbhi.NAveChgRate = ite.NAveChgRate
				tbhi.LlVolume = ite.LlVolume
				tbhi.LlValue = ite.LlValue
				tbhi.NStockID = ite.NStockID
				tbhi.NStockChgRate = ite.NStockChgRate
				tbhi.SzSName = strings.Replace(string(ite.SzSName[:]), "\u0000", "", -1)
				tbhi.NNum = ite.NNum
				tbhi.NChgRatio = ite.NChgRatio
				tbhi.NLongNum = ite.NLongNum
				tbhi.NShortNum = ite.NShortNum
				tbhi.LlValueOfInFlow = ite.LlValueOfInFlow
				tbhi.NBaseBlockID = stb.NBaseBlockID
				tbhi.NStockID = stb.NStockID
				tbhi.NLastPx = stb.NLastPx
				psbs.List = append(psbs.List, &tbhi)
			}
		}
	}

	ret := &protocol.PayloadStockBlockSet{
		SetID: req.SetID,
		Num:   int32(len(psb.List)),
		List:  psbs.List,
	}

	return ret, nil
}

// 返回板块定义集合
func GetPayloadStockBlock(key string) protocol.PayloadStockBlockBase {
	ls, err := hsgrr.Strings(RedisStore.Do("keys", key))
	if err != nil {
		logging.Info("error:%v", err)
	}

	//var psb []protocol.TagStockBlockBase
	var psb protocol.PayloadStockBlockBase
	for _, v := range ls {
		bin, err := RedisStore.GetBytes(v)
		if err != nil {
			if err == hsgrr.ErrNil {
				logging.Warning("redis not found key: %v", key)
				return psb
			}
			return psb
		}
		data := protocol.TagStockBlockBase{}
		buffer := bytes.NewBuffer(bin)
		if err := binary.Read(buffer, binary.LittleEndian, &data); err != nil && err != io.EOF {
			logging.Error("binary decode error: %v", err)
			return psb
		}
		psb.List = append(psb.List, &data)
	}
	return psb
}

// 返回排序集合
func GetSortStockBlockL(Sortkey string, rule int32) ([]*StockBlockInfo, error) {
	bin, err := RedisStore.GetBytes(Sortkey)
	if err != nil {
		if err == hsgrr.ErrNil {
			logging.Warning("redis not found key: %v", Sortkey)
			return nil, err
		}
		return nil, err
	}
	// 对应redis结构体
	var sbinl []*StockBlockInfo
	var sbin StockBlockInfo
	size := binary.Size(&sbin)
	// 循环 取每股结构体 100字节
	if len(bin) <= 0 {
		return nil, err
	}
	for i := 0; i < len(bin); i += size {

		v := bin[i : i+size]
		bfi := StockBlockInfo{}
		buffer := bytes.NewBuffer(v)
		if err = binary.Read(buffer, binary.LittleEndian, &bfi); err != nil && err != io.EOF {
			return nil, err
		}
		sbinl = append(sbinl, &bfi)
	}

	var sl []*StockBlockInfo
	// 返回正序集合
	if rule == 101 {
		sl = sbinl
	}
	// 返回倒序集合
	if rule == 102 {
		for i := len(sbinl) - 1; i >= 0; i-- {
			sl = append(sl, sbinl[i])
		}
	}
	return sl, err
}
