//成份证券
package publish

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/share/logging"

	"github.com/golang/protobuf/proto"
	. "haina.com/market/hqpublish/models"
)

type Element struct {
	Model `db:"-"`
}

func NewElement(redis_key string) *Element {
	return &Element{
		Model: Model{
			CacheKey: redis_key,
		},
	}
}

func (this *Element) GetBlockElementReply(req *protocol.RequestElement) (*protocol.PayloadElement, error) {
	if req.Begin < 0 || req.Num < 0 {
		return nil, INVALID_REQUEST_PARA
	}

	//------------------------------------板块下的成份股------------------------------------------------------//
	key := fmt.Sprintf(this.CacheKey, req.Classify, req.SetID)

	slist, err := RedisStore.GetBytes(key)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, err
	}
	var elms = &protocol.ElementList{}
	if err = proto.Unmarshal(slist, elms); err != nil {
		logging.Debug("%v", err.Error())
		return nil, err
	}
	logging.Debug("--------%v", len(elms.List))

	//-------------------------------------排序结果数据-------------------------------------------------------//
	dkey := fmt.Sprintf("hq:sort:2:%d", absInt32(req.FieldID))
	data, err := RedisStore.GetBytes(dkey)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, err
	}

	var stk = &TagStockSortInfo{}
	size := binary.Size(stk)

	var stocks []*protocol.TagStockSortInfo
	for i := 0; i < len(data); i += size {
		stock := &TagStockSortInfo{}
		if err = binary.Read(bytes.NewBuffer(data[i:i+size]), binary.LittleEndian, stock); err != nil && err != io.EOF {
			logging.Debug("%v", err.Error())
			return nil, err
		}

		for _, elm := range elms.List {
			if stock.NSID == elm.NSid {
				pstock := &protocol.TagStockSortInfo{
					NSID:              stock.NSID,
					NLastPx:           stock.NLastPx,
					NOpenPx:           stock.NOpenPx,
					NHighPx:           stock.NHighPx,
					NLowPx:            stock.NLowPx,
					NPreClosePx:       stock.NPreClosePx,
					LlVolume:          stock.LlVolume,
					LlValue:           stock.LlValue,
					NPxChgRatio:       stock.NPxChgRatio,
					NPxAmplitude:      stock.NPxAmplitude,
					NPxChg:            stock.NPxChg,
					NPE:               stock.NPE,
					NPB:               stock.NPB,
					NLiangbi:          stock.NLiangbi,
					NWeibi:            stock.NWeibi,
					LlMarketVal:       stock.LlMarketVal,
					LlFlowVal:         stock.LlFlowVal,
					NAveBidPx:         stock.NAveBidPx,
					NAveOfferPx:       stock.NAveOfferPx,
					LlBidVol:          stock.LlBidVol,
					LlOfferVol:        stock.LlOfferVol,
					NBid1Px:           stock.NBid1Px,
					NOffer1Px:         stock.NOffer1Px,
					LlBid1Vol:         stock.LlBid1Vol,
					LlOffer1Vol:       stock.LlOffer1Vol,
					LlValueOfInFlow:   stock.LlValueOfInFlow,
					SzSName:           byte40ToString(stock.SzSName),
					SzIndusCode:       byte8ToString(stock.SzIndusCode),
					NPxChgRatioIn5Min: stock.NPxChgRatioIn5Min,
					NTurnOver:         stock.NTurnOver,
				}
				stocks = append(stocks, pstock)
				break
			}
		}
	}
	logging.Debug("----%v", len(stocks))

	if len(stocks)-1 < int(req.Begin) {
		return nil, INVALID_REQUEST_PARA
	}

	if req.FieldID < 0 {
		swapElement(&stocks)
	}

	var board []*protocol.TagStockSortInfo

	if req.Num == 0 {
		board = stocks[0:]
	} else {
		if req.Num+req.Begin < int32(len(stocks)) {
			board = stocks[req.Begin : req.Num+req.Begin]
		} else {
			board = stocks[req.Begin:]
		}
	}

	payload := &protocol.PayloadElement{
		Classify: req.Classify,
		SetID:    req.SetID,
		FieldID:  req.FieldID,
		Total:    int32(len(stocks)),
		Begin:    req.Begin,
		Num:      int32(len(board)),
		List:     board,
	}

	return payload, nil
}

func swapElement(table *[]*protocol.TagStockSortInfo) {
	lengh := len(*table)

	for i := 0; i < lengh; i++ {
		(*table)[i], (*table)[lengh-i-1] = (*table)[lengh-i-1], (*table)[i]
		if i == lengh-i-2 || i == lengh-i-3 {
			break
		}
	}
}
