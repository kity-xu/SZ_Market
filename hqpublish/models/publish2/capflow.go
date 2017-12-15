//个股资金流向
package publish2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"haina.com/market/hqpublish/models/szdb"
	. "haina.com/share/models"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	. "haina.com/market/hqpublish/models"
	"haina.com/share/logging"
)

/// 资金统计
type TagTradeScaleStat struct {
	NSID               int32 ///< 股票代码
	NTime              int32 ///< 当前时间
	LlHugeBuyValue     int64 ///< 特大买单成交额*10000
	LlBigBuyValue      int64 ///< 大买单成交额*10000
	LlMiddleBuyValue   int64 ///< 中买单成交额*10000
	LlSmallBuyValue    int64 ///< 小买单成交额*10000
	LlHugeSellValue    int64 ///< 特大卖单成交额*10000
	LlBigSellValue     int64 ///< 大卖单成交额*10000
	LlMiddleSellValue  int64 ///< 中卖单成交额*10000
	LlSmallSellValue   int64 ///< 小卖单成交额*10000
	LlHugeBuyVolume    int64 ///< 特大买单成交量
	LlBigBuyVolume     int64 ///< 大买单成交量
	LlMiddleBuyVolume  int64 ///< 中买单成交量
	LlSmallBuyVolume   int64 ///< 小买单成交量
	LlHugeSellVolume   int64 ///< 特大卖单成交量
	LlBigSellVolume    int64 ///< 大卖单成交量
	LlMiddleSellVolume int64 ///< 中卖单成交量
	LlSmallSellVolume  int64 ///< 小卖单成交量
	LlValueOfInFlow    int64 ///<资金净流入额(*10000)
}

// 个股资金流向拼接总结构
type PayloadFundflow struct {
	SID     int32                       `json:"sid"`
	Num     int32                       `json:"num"`
	Last    *protocol.TagTradeScaleStat `json:"last"`
	Funds   []*protocol.Fund            `json:"flows"`
	CapDays []*protocol.Fund            `json:"capDays"`
}

type Fundflow struct {
	Model `db:"-"`
}

func NewFundflow(redis_key string) *Fundflow {
	return &Fundflow{
		Model: Model{
			CacheKey: redis_key,
		},
	}
}

func (this *Fundflow) GetFundflowReply(sid int32) (*PayloadFundflow, error) {
	key := fmt.Sprintf(this.CacheKey, sid)

	str, err := RedisStore.LRange(key, 0, -1)
	if err != nil {
		logging.Error("%v", err.Error())
		return nil, err
	}

	var funds []*protocol.Fund
	var ele = TagTradeScaleStat{}
	for _, data := range str {
		if err = binary.Read(bytes.NewBuffer([]byte(data)), binary.LittleEndian, &ele); err != nil && err != io.EOF {
			logging.Error("%v", err.Error())
			return nil, err
		}
		fund := &protocol.Fund{
			NTime: ele.NTime,
			Flow:  ele.LlValueOfInFlow,
		}
		funds = append(funds, fund)
	}

	keyd := fmt.Sprintf("hq:trade:day:%d", sid)
	data, err := RedisStore.GetBytes(keyd)
	if err != nil {
		logging.Error("%v", err.Error())
		return nil, err
	}
	if err = binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &ele); err != nil && err != io.EOF {
		logging.Error("%v", err.Error())
		return nil, err
	}

	var payload = &PayloadFundflow{
		SID: sid,
		Num: int32(len(funds)),
		Last: &protocol.TagTradeScaleStat{
			NTime:             ele.NTime,
			LlHugeBuyValue:    ele.LlHugeBuyValue,
			LlBigBuyValue:     ele.LlBigBuyValue,
			LlMiddleBuyValue:  ele.LlMiddleBuyValue,
			LlSmallBuyValue:   ele.LlSmallBuyValue,
			LlHugeSellValue:   ele.LlHugeSellValue,
			LlBigSellValue:    ele.LlBigSellValue,
			LlMiddleSellValue: ele.LlMiddleSellValue,
			LlSmallSellValue:  ele.LlSmallSellValue,
		},
		CapDays: this.getCapflowDays(sid),
		Funds:   funds,
	}
	return payload, nil
}

func (f *Fundflow) getCapflowDays(sid int32) []*protocol.Fund {
	capdays, err := szdb.NewSZ_HQ_SECURITYFUNDFLOW().GetHisSecurityFlow(4, sid)
	if len(capdays) == 0 || err != nil {
		return nil
	}

	var funds []*protocol.Fund
	for _, v := range capdays {
		flow := &protocol.Fund{
			NTime: v.TRADEDATE,
			Flow:  int64(v.HUGEBUYVALUE.Float64 + v.BIGBUYVALUE.Float64 - v.HUGESELLVALUE.Float64 - v.BIGSELLVALUE.Float64),
		}
		funds = append(funds, flow)
	}
	return funds
}
