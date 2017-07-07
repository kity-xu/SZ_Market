//排序
package publish

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	//	"haina.com/share/logging"

	"github.com/golang/protobuf/proto"
	. "haina.com/market/hqpublish/models"
)

// 字符串型数据长度定义
const (
	SECURITY_CODE_LEN = 24 ///< 证券代码长度
	SECURITY_NAME_LEN = 40 ///< 证券名称长度
	SECURITY_DESC_LEN = 8  ///< 英文简称
	INDUSTRY_CODE_LEN = 8  ///< 行业代码
	SECURITY_ISIN_LEN = 16 ///< 证券国际代码信息
)

type TagStockSortInfo struct {
	NSID              int32                   //SID
	NLastPx           int32                   //最新价(*10000)
	NOpenPx           int32                   //开盘价(*10000)
	NHighPx           int32                   //最高价(*10000)
	NLowPx            int32                   //最低价(*10000)
	NPreClosePx       int32                   //昨收价
	LlVolume          int64                   ///< 成交量
	LlValue           int64                   ///< 总成交额(*10000)
	NPxChgRatio       int32                   ///< 涨幅(*10000)
	NPxAmplitude      int32                   //振幅(*10000)
	NPxChg            int32                   //涨跌(*10000)
	NPE               int32                   //市盈（动）(*10000)
	NPB               int32                   //市净（动）(*10000)
	NLiangbi          int32                   //量比(*10000)
	NWeibi            int32                   //委比(*10000)
	LlMarketVal       int64                   //总市值
	LlFlowVal         int64                   //流通市值
	NAveBidPx         int32                   //委买均价(*10000)
	NAveOfferPx       int32                   //委卖均价(*10000)
	LlBidVol          int64                   //委买总量
	LlOfferVol        int64                   //委卖总量
	NBid1Px           int32                   //买一价(*10000)
	NOffer1Px         int32                   //卖一价(*10000)
	LlBid1Vol         int64                   //买一量
	LlOffer1Vol       int64                   //卖一量
	LlValueOfInFlow   int64                   //资金净流入额(*10000)
	SzSName           [SECURITY_NAME_LEN]byte //证券代码名称
	SzIndusCode       [INDUSTRY_CODE_LEN]byte ///< 行业代码
	NPxChgRatioIn5Min int32                   ///5分钟涨跌幅(*10000)
}

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

//ret := uint32(bs[0]) | uint32(bs[1])<<8 | uint32(bs[2])<<16 | uint32(bs[3])<<24 //小端
// binary.LittleEndian.Uint32(bs)
func (this *Sort) GetSortByFieldID(req *protocol.RequestSort) (*protocol.RedisSortTable, error) {
	key := fmt.Sprintf(this.CacheKey, req.SetID, absInt32(req.FieldID))

	bdata, err := RedisStore.GetBytes(key)
	if err != nil {
		return nil, err
	}
	if len(bdata) == 0 {
		return nil, ERROR_REDIS_DATE_NULL
	}

	var table = &protocol.RedisSortTable{}
	var sortSize = &TagStockSortInfo{}
	size := binary.Size(sortSize)

	for i := 0; i < len(bdata); i += size {
		sort := bdata[i : i+size]

		buffer := bytes.NewBuffer(sort)
		v := &TagStockSortInfo{}
		if err := binary.Read(buffer, binary.LittleEndian, v); err != nil && err != io.EOF {
			return nil, err
		}
		info := &protocol.TagStockSortInfo{
			NSID:              v.NSID,
			NLastPx:           v.NLastPx,
			NOpenPx:           v.NOpenPx,
			NHighPx:           v.NHighPx,
			NLowPx:            v.NLowPx,
			NPreClosePx:       v.NPreClosePx,
			LlVolume:          v.LlVolume,
			LlValue:           v.LlValue,
			NPxChgRatio:       v.NPxChgRatio,
			NPxAmplitude:      v.NPxAmplitude,
			NPxChg:            v.NPxChg,
			NPE:               v.NPE,
			NPB:               v.NPB,
			NLiangbi:          v.NLiangbi,
			NWeibi:            v.NWeibi,
			LlMarketVal:       v.LlMarketVal,
			LlFlowVal:         v.LlFlowVal,
			NAveBidPx:         v.NAveBidPx,
			NAveOfferPx:       v.NAveOfferPx,
			LlBidVol:          v.LlBidVol,
			LlOfferVol:        v.LlOfferVol,
			NBid1Px:           v.NBid1Px,
			NOffer1Px:         v.NOffer1Px,
			LlBid1Vol:         v.LlBid1Vol,
			LlOffer1Vol:       v.LlOffer1Vol,
			LlValueOfInFlow:   v.LlValueOfInFlow,
			SzSName:           byte40ToString(v.SzSName),
			SzIndusCode:       byte8ToString(v.SzIndusCode),
			NPxChgRatioIn5Min: v.NPxChgRatioIn5Min,
		}
		table.List = append(table.List, info)
	}
	return table, nil
}

func (this *Sort) GetPayloadSort(req *protocol.RequestSort) (*protocol.PayloadSort, error) {
	fkey := fmt.Sprintf("hq:sort:%d:%d:%s", req.SetID, req.FieldID, "f") //正序
	bkey := fmt.Sprintf("hq:sort:%d:%d:%s", req.SetID, req.FieldID, "b") //逆序

	var original = &protocol.RedisSortTable{}
	var sortTable = &protocol.RedisSortTable{}

	if req.FieldID >= 0 { //正序
		bs, err := GetCache(fkey)
		if err != nil {
			original, err = this.GetSortByFieldID(req)
			if err != nil {
				return nil, err
			}

			//入 redis cache
			data, err := proto.Marshal(original)
			if err != nil {
				return nil, err
			}

			if err = SetCache(fkey, TTL.Sort, data); err != nil {
				return nil, err
			}
		} else {
			if err = proto.Unmarshal(bs, original); err != nil {
				return nil, err
			}
		}

		//按所传参数做数据解析
		if len(original.List)-1 < int(req.Begin) || req.Begin < 0 || req.Num < 0 {
			return nil, INVALID_REQUEST_PARA
		}

		if int(req.Begin+req.Num) < len(original.List) {
			sortTable.List = original.List[req.Begin : req.Begin+req.Num]
		} else {
			sortTable.List = original.List[req.Begin:]
		}

	} else { //逆序
		bs, err := GetCache(bkey)
		if err != nil {
			original, err = this.GetSortByFieldID(req)
			if err != nil {
				return nil, err
			}

			swapSort(original) //逆序排序
			//入 redis cache
			data, err := proto.Marshal(original)
			if err != nil {
				return nil, err
			}

			if err = SetCache(bkey, TTL.Sort, data); err != nil {
				return nil, err
			}
		} else {
			if err = proto.Unmarshal(bs, original); err != nil {
				return nil, err
			}
		}

		if len(original.List)-1 < int(req.Begin) || req.Begin < 0 || req.Num < 0 {
			return nil, INVALID_REQUEST_PARA
		}
		if int(req.Begin+req.Num) < len(original.List) {
			sortTable.List = original.List[req.Begin : req.Begin+req.Num]
		} else {
			sortTable.List = original.List[req.Begin:]
		}
	}

	payload := &protocol.PayloadSort{
		SetID:   req.SetID,
		FieldID: req.FieldID,
		Total:   int32(len(original.List)),
		Begin:   req.Begin,
		Num:     int32(len(sortTable.List)),
		List:    sortTable.List,
	}
	return payload, nil
}

func byte40ToString(src [40]byte) string {
	var des []byte

	for _, v := range src {
		if v == '\u0000' || v == '0' {
			break
		}
		des = append(des, v)
	}
	return string(des)
}

func byte8ToString(src [8]byte) string {
	var des []byte

	for _, v := range src {
		if v == '\u0000' || v == '0' {
			break
		}
		des = append(des, v)
	}
	return string(des)
}

func byte12ToString(src [12]byte) string {
	var ss []byte
	for _, v := range src {
		if v == '0' || v == '\u0000' {
			break
		}
		ss = append(ss, v)
	}
	return string(ss)
}

//首尾交换未知
func swapSort(table *protocol.RedisSortTable) {
	lengh := len(table.List)

	for i := 0; i < lengh; i++ {
		table.List[i], table.List[lengh-i-1] = table.List[lengh-i-1], table.List[i]
		if i == lengh-i-2 || i == lengh-i-3 {
			break
		}
	}
}

func absInt32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
