package chipService

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"
	"time"

	"haina.com/market/f9/models"
	"haina.com/market/f9/services"
	"haina.com/share/logging"
)

type Chip struct {
	Chip_price  float64 `json:"chip_price"`
	Chip_text   string  `json:"chip_text"`
	Chip_charts []Klist `json:"Chip_charts"`
}

type KlistNow struct {
	NSID          int32 //int< 证券ID  4个字节
	NTime         int32 //int< 时间 unix time
	NTradingPhase int32 //unsigned int< 详细见产品交易阶段(TradingPhase) 4个字节
	NPreClosePx   int32 //unsigned int< 昨收价 * 10000
	NOpenPx       int32 //unsigned int< 开盘价 ..
	NHighPx       int32 //unsigned int< 最高价 ..
	NLowPx        int32 //unsigned int< 最低价 ..
	NLastPx       int32 //unsigned int< 最新价 ..
	NHighLimitPx  int32 //unsigned int< 涨停价格(*10000)
	NLowLimitPx   int32 //unsigned int< 跌停价格(*10000)
	LlTradeNum    int64 //long long< 成交笔数   8个字节
	LlVolume      int64 //long long< 成交量
	LlValue       int64 //long long< 成交额(*10000)
	NQuoteSize    int32 //int	< 报价总档数
	NWABidPx      int32 //unsigned int < 加权平均委买均价(*10000)
	NWAOfferPx    int32 //unsigned int< 加权平均委卖均价(*10000)
	LlToBidVol    int64 //long long< 总委买量
	LlToOfferVol  int64 //long long< 总委卖量
	LlInnerVolume int64 //long long< 内盘成交量
	LlOuterVolume int64 //long long	< 外盘成交量
	LlInnerValue  int64 //long long< 内盘成交额
	LlOuterValue  int64 //long long< 外盘成交额
	NPxChg        int32 //int	< 涨跌
	NPxChgRatio   int32 //int	< 涨跌幅*10000
	NPxAmplitude  int32 //int	 < 振幅*10000
	NLiangbi      int32 //int	< 量比*100
	NWeibi        int32 //int	< 委比*10000
	NTurnOver     int32 //int	< 换手率*10000
	NPE           int32 //int	< 动态市盈率*10000
	NPB           int32 //int	< 动态市净率*10000
}

type Klist struct {
	Date     int64   `json:"date"`
	Value    int64   `json:"value,omitempty"`
	NSid     int64   `json:"nSID,omitempty"`
	NTime    int64   `json:"nTime,omitempty"`
	NPreCPx  float64 `json:"nPreCPx"`
	NOpenPx  float64 `json:"nOpenPx"`
	NHighPx  float64 `json:"nHighPx"`
	NLastPx  float64 `json:"nLastPx"`
	LlVolume int64   `json:"llVolume,omitempty"`
	LlValue  int64   `json:"llValue,omitempty"`
	NAvgPx   int64   `json:"nAvgPx,omitempty"`
}

var public_N float64
var public_O string

func GetChipData(scode string) (Chip, error) {

	//getKlineNow()
	//RedisCache()
	//	RedisCache().

	//	redis.Set("name", []byte("123456"))
	//	data, _ := redis.Get("hq:st:snap:100600000")
	//	logging.Info("data===%v", data)
	//	var line KlistNow
	//	size := binary.Size(&line)
	//	logging.Info("size===%v", size)
	//	for i := 0; i < size; i += size {
	//		l := &KlistNow{}
	//		buffer := bytes.NewBuffer([]byte(data[i:size]))
	//		if err := binary.Read(buffer, binary.LittleEndian, l); err != nil && err != io.EOF {
	//			logging.Error("error===%v", err)
	//		}
	//		//break
	//		logging.Info("----------------line===%v", l)
	//	}

	klistNow, _ := getKlineNow() //当日K线
	logging.Info("knnow===%+v", klistNow)
	klistPastAsc, _ := getKlinePast() //历史K线

	logging.Info("历史K线=klistPastAsc===%+v", klistPastAsc)

	var klistPast []Klist
	for i := len(klistPastAsc) - 1; i >= 0; i-- {
		var kp Klist
		kp.Date = klistPastAsc[i].NTime
		kp.NPreCPx = klistPastAsc[i].NPreCPx
		kp.NOpenPx = klistPastAsc[i].NOpenPx
		kp.NHighPx = klistPastAsc[i].NHighPx
		kp.NLastPx = klistPastAsc[i].NLastPx
		kp.LlVolume = klistPastAsc[i].LlVolume
		kp.LlValue = klistPastAsc[i].LlValue
		klistPast = append(klistPast, kp)
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

	var klists []Klist

	for i := 0; i < len(klistPastAsc); i++ {
		var kp Klist
		kp.Date = klistPastAsc[i].NTime
		kp.NPreCPx = klistPastAsc[i].NPreCPx / 10000
		kp.NOpenPx = klistPastAsc[i].NOpenPx / 10000
		kp.NHighPx = klistPastAsc[i].NHighPx / 10000
		kp.NLastPx = klistPastAsc[i].NLastPx / 10000
		kp.LlVolume = klistPastAsc[i].LlVolume
		kp.LlValue = klistPastAsc[i].LlValue
		klists = append(klists, kp)
	}

	klistNow.NPreCPx = klistNow.NPreCPx / 10000
	klistNow.NOpenPx = klistNow.NOpenPx / 10000
	klistNow.NHighPx = klistNow.NHighPx / 10000
	klistNow.NLastPx = klistNow.NLastPx / 10000

	klists = append(klists, klistNow) //合并后的总K线

	var ch Chip
	ch.Chip_price = public_N
	ch.Chip_text = "近期该股筹码平均交易成本为" + strconv.FormatFloat(public_N, 'f', -1, 64) + "元，" + public_O + "。"
	ch.Chip_charts = klists

	//var chip_chart []*chip_charts

	//logging.Info("vwap5_T====", vwap5_T)
	logging.Info("klistPast====%+v", klists)
	return ch, nil
}

//计算成交量乘以收盘价的之和
//从一个数组中计算指定开始日期向后推n天的数据总和，如果不够n天，有几天算几天 指定开始日期从0开始,
func sumOfday(Klist []Klist, start int, n int) float64 {
	if start > len(Klist) {
		return float64(0)
	}
	var num float64 = 0
	for i := start; i < len(Klist); i++ {
		if i < n {
			if Klist[0].LlVolume == 0 {
				n++
			} else {
				num = num + float64(Klist[i].NLastPx/10000)*float64(Klist[i].LlVolume)
			}
		} else {
			break
		}
	}
	return num
}

//计算成交量之和
func sumOfFieldByNum(Klist []Klist, start int, n int, feild string) float64 {
	if start > len(Klist) {
		return float64(0)
	}
	var num float64 = 0
	for i := start; i < len(Klist); i++ {
		if i < n {
			if feild == "nLastPx" {
				num = num + float64(Klist[i].NLastPx/10000)
			} else {
				num = num + float64(Klist[i].LlVolume)
			}

		} else {
			break
		}
	}
	return num
}

func getKlinePast() ([]Klist, error) {
	//url := "https://hbmk.0606.com.cn/api/hq/kline?format=json"
	url := "http://47.94.109.175:7280/api/hq/kline?format=json"
	//keysid, _ := strconv.Atoi(KeyNsid)
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	tm1 := tm.Format("20060102")
	tm2, _ := strconv.Atoi(tm1)
	posts := map[string]interface{}{
		"SID":       100600000,
		"Type":      1,
		"TimeIndex": tm2,
		"Num":       100,
		"Direct":    0,
	}
	var _param struct {
		Code int64 `json:"code"`
		Data struct {
			SID   int64   `json:"sid"`
			Type  int64   `json:"type"`
			Total int64   `json:"total"`
			Begin int64   `json:"begin"`
			Num   int64   `json:"num"`
			KList []Klist `json:"klist"`
		} `json:"data"`
	}
	err := commonService.HttpPostJson(url, posts, &_param)
	if err != nil {
		logging.Info(err.Error())
	}
	return _param.Data.KList, err
}

func getKlineNow() (Klist, error) {
	tNow := time.Now()
	timeNow := tNow.Format("20060102")
	t, _ := strconv.ParseInt(timeNow, 10, 64)

	redisData, err := models.RedisStore.Get("hq:st:snap:100600000")

	var kNow Klist
	if err != nil || len(redisData) == 0 {
		kNow.Date = t
		kNow.NPreCPx = 0
		kNow.NHighPx = 0
		kNow.NLastPx = 0
		kNow.LlVolume = 0
		kNow.LlValue = 0
		logging.Info("==", err.Error())
		return kNow, err
	} else {
		//var kno []KlistNow
		var kn KlistNow
		size := binary.Size(&kn)

		k := &KlistNow{}
		buffer := bytes.NewBuffer([]byte(redisData[0:size]))
		if err := binary.Read(buffer, binary.LittleEndian, k); err != nil && err != io.EOF {
			logging.Info("error===%v", err)
		}
		//kno = append(kno, k)
		logging.Info("kkk====%+v", k)

		//
		kNow.Date = t
		kNow.NPreCPx = float64(k.NPreClosePx)
		kNow.NOpenPx = float64(k.NOpenPx)
		kNow.NHighPx = float64(k.NHighPx)
		kNow.NLastPx = float64(k.NLastPx)
		kNow.LlVolume = k.LlVolume
		kNow.LlValue = k.LlValue
		return kNow, err
	}

}
