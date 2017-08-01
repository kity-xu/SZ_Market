//排序
package publish

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	//"haina.com/share/logging"

	. "haina.com/market/hqpublish/models"
)

type UserDefine struct {
	Model `db:"-"`
}

func NewUserDefine(redis_key string) *UserDefine {
	return &UserDefine{
		Model: Model{
			CacheKey: redis_key,
		},
	}
}

//ret := uint32(bs[0]) | uint32(bs[1])<<8 | uint32(bs[2])<<16 | uint32(bs[3])<<24 //小端
// binary.LittleEndian.Uint32(bs)
func (this *UserDefine) GetSortByFieldID(req *protocol.RequestUserdef) (*protocol.RedisSortTable, error) {
	key := fmt.Sprintf(this.CacheKey, 2, absInt32(req.FieldID)) //个股

	bdata, err := RedisStore.GetBytes(key)
	if err != nil {
		return nil, err
	}
	if len(bdata) == 0 {
		return nil, ERROR_REDIS_DATE_NULL
	}

	//	ikey := fmt.Sprintf(this.CacheKey, 1, absInt32(req.FieldID)) //指数
	//	idata, err := RedisStore.GetBytes(ikey)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if len(idata) == 0 {
	//		return nil, ERROR_REDIS_DATE_NULL
	//	}
	//	bdata = append(bdata, idata...) // 个股+指数

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
			NTurnOver:         v.NTurnOver,
		}
		table.List = append(table.List, info)
	}
	return table, nil
}

func (this *UserDefine) GetSecurityUserdefine(req *protocol.RequestUserdef) (*protocol.PayloadUserdef, error) {
	table, err := this.GetSortByFieldID(req)
	if err != nil {
		return nil, err
	}

	userdef, err := sidSearch(table, req.Sids)
	if err != nil {
		return nil, err
	}
	if req.FieldID < 0 {
		reverseSort(userdef)
	}

	payload := &protocol.PayloadUserdef{
		FieldID: req.FieldID,
		Num:     int32(len(userdef.List)),
		STList:  userdef.List,
	}
	return payload, nil
}

func sidSearch(table *protocol.RedisSortTable, src []int32) (*protocol.RedisSortTable, error) {
	var lengh int = len(table.List)
	if lengh == 0 {
		return nil, fmt.Errorf("src data is null...")
	}

	var result = &protocol.RedisSortTable{}

	for _, tb := range table.List {
		for _, v := range src {
			if tb.NSID == v {
				result.List = append(result.List, tb)
				break
			}
		}
	}
	return result, nil
}
