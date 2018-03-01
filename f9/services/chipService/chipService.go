package chipService

import (
	"ProtocolBuffer/projects/f9/go/protocol"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"

	"haina.com/market/f9/models"
	"haina.com/market/f9/services"
	"haina.com/share/logging"
)

type Chip struct {
	Chip_price  float64          `json:"chip_price"`
	Chip_text   string           `json:"chip_text"`
	Chip_charts []protocol.KInfo `json:"Chip_charts"`
}

//
//type KlistNow struct {
//	NSID          int32  //int< 证券ID  4个字节
//	NTime         int32  //int< 时间 unix time
//	NTradingPhase uint32 //unsigned int< 详细见产品交易阶段(TradingPhase) 4个字节
//	NPreClosePx   uint32 //unsigned int< 昨收价 * 10000
//	NOpenPx       uint32 //unsigned int< 开盘价 ..
//	NHighPx       uint32 //unsigned int< 最高价 ..
//	NLowPx        uint32 //unsigned int< 最低价 ..
//	NLastPx       uint32 //unsigned int< 最新价 ..
//	NHighLimitPx  uint32 //unsigned int< 涨停价格(*10000)
//	NLowLimitPx   uint32 //unsigned int< 跌停价格(*10000)
//	LlTradeNum    int64  //long long< 成交笔数   8个字节
//	LlVolume      int64  //long long< 成交量
//	LlValue       int64  //long long< 成交额(*10000)
//	NQuoteSize    int32  //int	< 报价总档数
//	NWABidPx      uint32 //unsigned int < 加权平均委买均价(*10000)
//	NWAOfferPx    uint32 //unsigned int< 加权平均委卖均价(*10000)
//	LlToBidVol    int64  //long long< 总委买量
//	LlToOfferVol  int64  //long long< 总委卖量
//	LlInnerVolume int64  //long long< 内盘成交量
//	LlOuterVolume int64  //long long	< 外盘成交量
//	LlInnerValue  int64  //long long< 内盘成交额
//	LlOuterValue  int64  //long long< 外盘成交额
//	NPxChg        int32  //int	< 涨跌
//	NPxChgRatio   int32  //int	< 涨跌幅*10000
//	NPxAmplitude  int32  //int	 < 振幅*10000
//	NLiangbi      int32  //int	< 量比*100
//	NWeibi        int32  //int	< 委比*10000
//	NTurnOver     int32  //int	< 换手率*10000
//	NPE           int32  //int	< 动态市盈率*10000
//	NPB           int32  //int	< 动态市净率*10000
//}
//
//type Klist struct {
//	//Date     int64   `json:"date"`
//	//Value    int64   `json:"value,omitempty"`
//	NSid     int32  `json:"nSID,omitempty"`
//	NTime    int32  `json:"nTime,omitempty"`
//	NPreCPx  int32  `json:"nPreCPx"`
//	NOpenPx  int32  `json:"nOpenPx"`
//	NHighPx  int32  `json:"nHighPx"`
//	NLastPx  int32  `json:"nLastPx"`
//	LlVolume int64  `json:"llVolume,omitempty"`
//	LlValue  int64  `json:"llValue,omitempty"`
//	NAvgPx   uint32 `json:"nAvgPx,omitempty"`
//}

var public_N float64
var public_O string

func GetChipData(sid string) (*Chip, error) {
	SID, _ := strconv.Atoi(sid)
	klistNow, err := getKlineNow(SID) //当日K线
	if err != nil {
		logging.Error("The dayLine is null")
		return nil, err
	}
	klistPast, err := getKlinePast(SID) //历史K线
	if len(*klistPast) == 0 || err != nil {
		logging.Error("The historical dayLine is null")
		return nil, err
	}

	var vwap5_T, vwap5_T1, vwap5_T3, T_tmp float64

	if klistNow.LlVolume != 0 {
		vwap5_T = (float64(klistNow.NLastPx/10000)*float64(klistNow.LlVolume) + sumOfday(klistPast, 0, 4)) / (float64(klistNow.LlVolume) + sumOfFieldByNum(klistPast, 0, 4, "llVolume"))
	} else {
		vwap5_T = sumOfday(klistPast, 0, 5) / sumOfFieldByNum(klistPast, 0, 5, "llVolume")
	}
	logging.Info("vwap5_T===%v", vwap5_T)
	var T9_sum float64 = 0
	for i := 1; i <= 9; i++ {
		if sumOfday(klistPast, i, i+5) != 0 {
			T9_sum = T9_sum + sumOfday(klistPast, i, i+5)/sumOfFieldByNum(klistPast, i, i+5, "llVolume")
		}
	}
	logging.Info("T9_sumT9_sum===%v", T9_sum)

	public_N = (vwap5_T + T9_sum) / 10
	vwap5_T1 = sumOfday(klistPast, 1, 6) / sumOfFieldByNum(klistPast, 1, 6, "llVolume")
	vwap5_T3 = sumOfday(klistPast, 3, 8) / sumOfFieldByNum(klistPast, 3, 8, "llVolume")
	T_tmp = (vwap5_T1 / vwap5_T3) - 1

	if vwap5_T >= T9_sum && T_tmp > 0 && T_tmp <= 0.03 {
		public_O = "该股获筹码青睐，且集中度渐增"
	}
	if vwap5_T >= T9_sum && T_tmp >= 0.03 {
		public_O = "近日该股快速猛烈吸筹,短线操作建议强烈关注"
	}
	if vwap5_T >= T9_sum && T_tmp < 0 {
		public_O = "该股有吸筹现象，但吸筹力度不强"
	}
	if vwap5_T < T9_sum && T_tmp < -0.03 {
		public_O = "近日筹码快速出逃，建议调仓换股"
	}
	if vwap5_T < T9_sum && T_tmp >= -0.03 && T_tmp < 0 {
		public_O = "筹码减仓，但减仓程度减缓"
	}
	if vwap5_T < T9_sum && T_tmp >= 0 {
		public_O = "筹码关注程度减弱"
	}

	klistNow.NPreCPx = klistNow.NPreCPx / 10000
	klistNow.NOpenPx = klistNow.NOpenPx / 10000
	klistNow.NHighPx = klistNow.NHighPx / 10000
	klistNow.NLastPx = klistNow.NLastPx / 10000

	market, err := GetMarketStatus(100000000)
	if err != nil {
		return nil, err
	}

	l := len(*klistPast)
	var ktable []protocol.KInfo
	if l > 1 && (*klistPast)[l-1].LlVolume == 0 {
		ktable = append(ktable, (*klistPast)[0:l-2]...)
	}

	lt := len(ktable)
	if lt == 0 {
		ktable = append(ktable, *klistNow)
	} else {
		if market.NTradeDate > ktable[lt-1].NTime {
			ktable = append(ktable, *klistNow) //合并后的总K线
		}
	}

	var ch Chip
	ch.Chip_price = public_N
	ch.Chip_text = "近期该股筹码平均交易成本为" + strconv.FormatFloat(public_N, 'f', -1, 64) + "元，" + public_O + "。"
	ch.Chip_charts = ktable

	return &ch, nil
}

//计算成交量乘以收盘价的之和
//从一个数组中计算指定开始日期向后推n天的数据总和，如果不够n天，有几天算几天 指定开始日期从0开始,
func sumOfday(Klist *[]protocol.KInfo, start int, n int) float64 {
	if start > len(*Klist) {
		return float64(0)
	}
	var num float64 = 0
	for i := start; i < len(*Klist); i++ {
		if i < n {
			if (*Klist)[0].LlVolume == 0 {
				n++
			} else {
				num = num + float64((*Klist)[i].NLastPx/10000)*float64((*Klist)[i].LlVolume)
			}
		} else {
			break
		}
	}
	return num
}

//计算成交量之和
func sumOfFieldByNum(Klist *[]protocol.KInfo, start int, n int, feild string) float64 {
	if start > len(*Klist) {
		return float64(0)
	}
	var num float64 = 0
	for i := start; i < len(*Klist); i++ {
		if i < n {
			if feild == "nLastPx" {
				num = num + float64((*Klist)[i].NLastPx/10000)
			} else {
				num = num + float64((*Klist)[i].LlVolume)
			}

		} else {
			break
		}
	}
	return num
}

func getKlinePast(sid int) (*[]protocol.KInfo, error) {
	url := "http://47.94.109.175:7280/api/hq/kline?format=json"

	posts := map[string]interface{}{
		"SID":       sid,
		"Type":      1,
		"TimeIndex": 0,
		"Num":       100,
		"Direct":    0,
	}
	var _param struct {
		Code int64 `json:"code"`
		Data struct {
			SID   int64            `json:"sid"`
			Type  int64            `json:"type"`
			Total int64            `json:"total"`
			Begin int64            `json:"begin"`
			Num   int64            `json:"num"`
			KList []protocol.KInfo `json:"klist"`
		} `json:"data"`
	}
	err := commonService.HttpPostJson(url, posts, &_param)
	if err != nil {
		logging.Info(err.Error())
	}
	return &_param.Data.KList, nil
}

func getKlineNow(sid int) (*protocol.KInfo, error) {
	//当前交易日
	market, err := GetMarketStatus(100000000)
	if err != nil {
		return nil, err
	}

	snapKey := fmt.Sprintf(models.REDISKEY_SNAP, sid)
	redisData, err := models.RedisStore.Get(snapKey)
	if err != nil || len(redisData) == 0 {
		logging.Error("get redis error:kline | %v", err.Error())
		return nil, err
	}
	snap := &protocol.StockSnapshot{}
	buffer := bytes.NewBuffer([]byte(redisData))
	if err := binary.Read(buffer, binary.LittleEndian, snap); err != nil && err != io.EOF {
		logging.Error("binary read error |%v", err)
		return nil, err
	}

	kNow := &protocol.KInfo{
		NSID:     int32(sid),
		NTime:    market.NTradeDate,
		NPreCPx:  int32(snap.NPreClosePx),
		NOpenPx:  int32(snap.NOpenPx),
		NHighPx:  int32(snap.NHighPx),
		NLowPx:   int32(snap.NLowPx),
		NLastPx:  int32(snap.NLastPx),
		LlVolume: snap.LlVolume,
		LlValue:  snap.LlValue,
		NAvgPx:   uint32(snap.LlVolume / snap.LlValue),
	}
	return kNow, err
}

//市场状态获取当前交易日
func GetMarketStatus(mid int) (*protocol.MarketStatus, error) {
	var obj protocol.MarketStatus

	key := fmt.Sprintf(models.REDISKEY_MARKET_STATUS, mid)
	data, err := models.RedisCache.GetBytes(key)
	if err != nil {
		data, err = models.RedisStore.GetBytes(key)
	}
	if err != nil || len(data) == 0 {
		logging.Error("Failed to obtain market status error.")
		return nil, err
	}

	buffer := bytes.NewBuffer([]byte(data))
	if err := binary.Read(buffer, binary.LittleEndian, &obj); err != nil && err != io.EOF {
		return nil, err
	}
	return &obj, nil
}
