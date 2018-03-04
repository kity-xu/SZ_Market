package fundService

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"sort"
	"strconv"
	"time"

	"haina.com/market/f9/models"
	"haina.com/market/f9/models/finchina"
	"haina.com/market/f9/models/fundModel"
	"haina.com/market/f9/services"
	"haina.com/share/logging"
)

type fundData struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
	Data    fund   `json:"data"`
}
type fund struct {
	NSID                  int32                   `json:"llHugeBuyValue"`
	LlHugeBuyValue        float64                 `json:"llHugeBuyValue"`
	LlBigBuyValue         float64                 `json:"llBigBuyValue"`
	LlMiddleBuyValue      float64                 `json:"llMiddleBuyValue"`
	LlSmallBuyValue       float64                 `json:"llSmallBuyValue"`
	LlHugeSellValue       float64                 `json:"llHugeSellValue"`
	LlBigSellValue        float64                 `json:"llBigSellValue"`
	LlMiddleSellValue     float64                 `json:"llMiddleSellValue"`
	LlSmallSellValue      float64                 `json:"llSmallSellValue"`
	LnflowValue           float64                 `json:"lnflowValue"`
	OutflowValue          float64                 `json:"outflowValue"`
	DayInflow             float64                 `json:"dayInflow"`
	Result                string                  `json:"result"`
	Fund_inflow           string                  `json:"fund_inflow"`
	Point_degree          float64                 `json:"point_degree"`
	Major_control         string                  `json:"major_control"`
	Maior_trend           string                  `json:"maior_trend"`
	Industry_inflow       string                  `json:"industry_inflow"`
	Industry_inflow_chart []industry_inflow_chart `json:"industry_inflow_chart"`
}
type industry_inflow_chart struct {
	Date           string  `json:"date"`
	Major_industry float64 `json:""major_industry"`
	Major_symbol   float64 `json:"major_symbol"`
}

func GetFundData(sid string) (interface{}, error) {
	var fu fundData
	basic, err := commonService.GetCommonData(sid)
	if err != nil {
		return nil, err
	}
	fundData, _ := getFundNow(sid) //当日的资金数据

	logging.Info("fundData====%+v", fundData)

	fu.Data.LlHugeBuyValue = fundData.LlHugeBuyValue
	fu.Data.LlBigBuyValue = fundData.LlBigBuyValue
	fu.Data.LlMiddleBuyValue = fundData.LlMiddleBuyValue
	fu.Data.LlSmallBuyValue = fundData.LlSmallBuyValue
	fu.Data.LlHugeSellValue = fundData.LlHugeSellValue
	fu.Data.LlBigSellValue = fundData.LlBigSellValue
	fu.Data.LlMiddleSellValue = fundData.LlMiddleSellValue
	fu.Data.LlSmallSellValue = fundData.LlSmallSellValue

	InflowValue := fu.Data.LlHugeBuyValue + fu.Data.LlBigBuyValue + fu.Data.LlMiddleBuyValue + fu.Data.LlSmallBuyValue      //流入
	OutflowValue := fu.Data.LlHugeSellValue + fu.Data.LlBigSellValue + fu.Data.LlMiddleSellValue + fu.Data.LlSmallSellValue //流出
	//今日主力净流入
	mainInflow := fu.Data.LlHugeBuyValue + fu.Data.LlBigBuyValue - fu.Data.LlHugeSellValue - fu.Data.LlBigSellValue
	//今日净值
	DayInflow := InflowValue - OutflowValue
	//主力资金占比%
	var mainProPortion float64
	if fu.Data.LlHugeBuyValue+fu.Data.LlBigBuyValue+fu.Data.LlMiddleBuyValue+fu.Data.LlSmallBuyValue != 0 {
		mainProPortion = mainInflow / (fu.Data.LlHugeBuyValue + fu.Data.LlBigBuyValue + fu.Data.LlMiddleBuyValue + fu.Data.LlSmallBuyValue)
	} else {
		mainProPortion = 0
	}
	hugeInflowValue := fu.Data.LlHugeBuyValue - fu.Data.LlHugeSellValue        //特大单净流入
	bigInflowValue := fu.Data.LlBigBuyValue - fu.Data.LlBigSellValue           //大单净流入
	middlesInflowValue := fu.Data.LlMiddleBuyValue - fu.Data.LlMiddleSellValue // 中单净流入
	smallInflowValue := fu.Data.LlSmallBuyValue - fu.Data.LlSmallSellValue     //小单净流入
	fu.Data.Fund_inflow = basic.Compname + "今日主力净流入" + strconv.FormatFloat(mainInflow/10000, 'f', -1, 64) + "万元，主力资金占比" + strconv.FormatFloat(mainProPortion, 'f', -1, 64) + "%。其中，特大单净流入" + strconv.FormatFloat(hugeInflowValue/10000, 'f', -1, 64) + "万元，大单净流入" + strconv.FormatFloat(bigInflowValue/10000, 'f', -1, 64) + "万元，中单净流入" + strconv.FormatFloat(middlesInflowValue/10000, 'f', -1, 64) + "万元，小单净流入" + strconv.FormatFloat(smallInflowValue/10000, 'f', -1, 64) + "万元。"
	fu.Data.DayInflow = DayInflow
	fu.Data.OutflowValue = OutflowValue
	fu.Data.LnflowValue = InflowValue

	logging.Info("fundData====%+v", fundData)
	fuData, _ := fundModel.NewFundFlow().GetFundData(basic.SID, 30)
	typecode, _ := fundModel.NewFundFlow().IndustryDateG(basic.Swlevelcode, 10)
	var fud []fund
	var chart []industry_inflow_chart
	for k, v := range fuData {
		var fun fund
		fun.LlHugeBuyValue = v.LlHugeBuyValue
		fun.LlBigBuyValue = v.LlBigBuyValue
		fun.LlMiddleBuyValue = v.LlMiddleBuyValue
		fun.LlSmallBuyValue = v.LlSmallBuyValue
		fun.LlHugeSellValue = v.LlHugeSellValue
		fun.LlBigSellValue = v.LlBigSellValue
		fun.LlMiddleSellValue = v.LlMiddleSellValue
		fun.LlSmallSellValue = v.LlSmallSellValue
		fud = append(fud, fun)
		if k < 10 {
			var cha industry_inflow_chart
			cha.Date = v.NTime
			cha.Major_industry = typecode[k].LlHugeBuyValue + typecode[k].LlBigBuyValue - typecode[k].LlHugeSellValue - typecode[k].LlBigSellValue
			cha.Major_symbol = v.LlHugeBuyValue + v.LlBigBuyValue - v.LlHugeSellValue - v.LlBigSellValue
			chart = append(chart, cha)
		}
	}
	fu.Data.Industry_inflow_chart = chart
	var funds []fund
	funds = append(funds, fundData)
	funds = append(funds, fud...)
	logging.Info("funds====%+v", funds)

	var cddd, lr, pDegree, U, V float64
	for i := 0; i < 5; i++ {
		cddd += Abs(funds[i].LlHugeBuyValue-funds[i].LlHugeSellValue) + Abs(funds[i].LlBigBuyValue-funds[i].LlBigSellValue)
		lr += (funds[i].LlHugeBuyValue + funds[i].LlBigBuyValue + funds[i].LlMiddleBuyValue + funds[i].LlSmallBuyValue)
	}
	if lr == 0 {
		pDegree = 0
	} else {
		pDegree = 200 * cddd / lr
	}
	var T string
	if pDegree >= 0 && pDegree < 20 {
		T = "力没有控盘，筹码分布非常分散"
	} else if pDegree >= 20 && pDegree < 60 {
		T = "主力轻度控盘，筹码分布较为分散"
	} else if pDegree >= 60 && pDegree < 120 {
		T = "主力中度控盘，且筹码比较集中"
	} else if pDegree >= 120 && pDegree < 200 {
		T = "主力高度控盘，且筹码非常集中"
	} else {
		T = "主力高度控盘，且筹码非常集中"
	}
	//【U】数值
	U = cddd / 10000
	//【V】数
	V = cddd / lr
	fu.Data.Major_control = "过去5天，该股票" + T + "，主力成交额" + strconv.FormatFloat(U/10000, 'f', -1, 64) + "亿元，占总成交额的" + strconv.FormatFloat(V, 'f', -1, 64) + "%。"

	//股票最近三条主力净流入数据，T1,T2,T3
	T1 := fuData[0].LlHugeBuyValue + fuData[0].LlBigBuyValue - fuData[0].LlHugeSellValue - fuData[0].LlBigSellValue
	T2 := fuData[1].LlHugeBuyValue + fuData[1].LlBigBuyValue - fuData[1].LlHugeSellValue - fuData[1].LlBigSellValue
	T3 := fuData[2].LlHugeBuyValue + fuData[2].LlBigBuyValue - fuData[2].LlHugeSellValue - fuData[2].LlBigSellValue
	var S, W string
	if T1 > 0 && T2 > 0 && T3 > 0 {
		S = "连续3日被主力资金增仓"
	} else if T1 > 0 && T2 > 0 && T3 <= 0 {
		S = "连续2日被主力资金增仓"
	} else if T1 < 0 && T2 < 0 && T3 < 0 {
		S = "连续3日被主力资金减仓"
	} else if T1 < 0 && T2 < 0 && T3 > 0 {
		S = "连续2日被主力资金减仓"
	} else {
		S = "无连续增减仓现象，主力趋势不明显"
	}
	fu.Data.Maior_trend = "该股当前" + S + "。"

	ddd, err := models.RedisCache.GetBytes("AAA")

	var allCompany []finchina.CompanyDetail
	if err == nil {
		json.Unmarshal(ddd, &allCompany)
	} else {
		allCompany, _ := commonService.IndustryOfAllCompany(basic.Swlevelcode)
		jsonData, _ := json.Marshal(allCompany)
		models.RedisCache.Setex("AAA", 10000, jsonData)
	}

	var fuud []*mainInflowData
	for _, v := range allCompany {
		var keysid string
		if v.EXCHANGE == "001002" {
			keysid = "100" + v.SYMBOL
		} else {
			keysid = "200" + v.SYMBOL
		}
		funData, _ := getFundNow(keysid)

		var fuu mainInflowData
		fuu.SYMBOL = funData.NSID
		fuu.MainInflow = funData.LlHugeBuyValue + funData.LlBigBuyValue - funData.LlHugeSellValue - funData.LlBigSellValue
		fuud = append(fuud, &fuu)
	}
	sort.Sort(sortText(fuud))
	var E int = 0
	for k, v := range fuud {
		a_int64, _ := strconv.ParseInt(basic.SID, 10, 64)
		if int64(v.SYMBOL) == a_int64 {
			E = k + 1
		}
	}

	//【W】判断条件，行业最近三条主力净流入数据，T1,T2,T3
	T1 = typecode[0].LlHugeBuyValue + typecode[0].LlBigBuyValue - typecode[0].LlHugeSellValue - typecode[0].LlBigSellValue
	T2 = typecode[1].LlHugeBuyValue + typecode[1].LlBigBuyValue - typecode[1].LlHugeSellValue - typecode[1].LlBigSellValue
	T3 = typecode[2].LlHugeBuyValue + typecode[2].LlBigBuyValue - typecode[2].LlHugeSellValue - typecode[2].LlBigSellValue
	if T1 > 0 && T2 > 0 && T3 > 0 {
		W = "该行业连续3日被主力资金增仓"
	} else if T1 > 0 && T2 > 0 && T3 <= 0 {
		W = "该行业连续2日被主力资金增仓"
	} else if T1 < 0 && T2 < 0 && T3 < 0 {
		W = "该行业连续3日被主力资金减仓"
	} else if T1 < 0 && T2 < 0 && T3 > 0 {
		W = "该行业连续2日被主力资金减仓"
	} else {
		W = "该行业当前无连续增减仓现象，主力趋势不明显"
	}
	tNow := time.Now()
	timeNow := tNow.Format("2006年01月02日")
	public_D := len(allCompany)
	fu.Data.Point_degree = pDegree
	fu.Data.Result = "该股当前" + S + "，过去5个交易日该股" + T + "，主力成交额" + strconv.FormatFloat(U/10000, 'f', -1, 64) + "亿元,占总成交额的" + strconv.FormatFloat(V, 'f', -1, 64) + "%。" + timeNow + "，" + basic.Compname + "主力净流入" + strconv.FormatFloat(mainInflow/10000, 'f', -1, 64) + "万元,主力资金占比" + strconv.FormatFloat(mainProPortion, 'f', -1, 64) + "%。"
	fu.Data.Industry_inflow = basic.Swlevelname + "" + W + "。" + basic.Compname + "主力资金占比在该行业中排名" + strconv.Itoa(E) + "/" + strconv.Itoa(public_D) + "。"

	return fu, nil

}

type sortText []*mainInflowData

func (I sortText) Len() int {
	return len(I)
}
func (I sortText) Less(i, j int) bool {
	return I[i].MainInflow > I[j].MainInflow
}
func (I sortText) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

type mainInflowData struct {
	SYMBOL     int32
	MainInflow float64
}

//绝对值
func Abs(f float64) float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

type FundNow struct {
	NSID               int32 //int< 股票代码
	NTime              int32 //int< 当前时间
	LlHugeBuyValue     int64 //long long< 特大买单成交额*10000
	LlBigBuyValue      int64 //long long< 大买单成交额*10000
	LlMiddleBuyValue   int64 //long long< 中买单成交额*10000
	LlSmallBuyValue    int64 //long long< 小买单成交额*10000
	LlHugeSellValue    int64 //long long< 特大卖单成交额*10000
	LlBigSellValue     int64 //long long< 大卖单成交额*10000
	LlMiddleSellValue  int64 //long long< 中卖单成交额*10000
	LlSmallSellValue   int64 //long long< 小卖单成交额*10000
	LlHugeBuyVolume    int64 //long long< 特大买单成交量
	LlBigBuyVolume     int64 //long long< 大买单成交量
	LlMiddleBuyVolume  int64 //long long< 中买单成交量
	LlSmallBuyVolume   int64 //long long< 小买单成交量
	LlHugeSellVolume   int64 //long long< 特大卖单成交量
	LlBigSellVolume    int64 //long long< 大卖单成交量
	LlMiddleSellVolume int64 //long long< 中卖单成交量
	LlSmallSellVolume  int64 //long long< 小卖单成交量
	LlValueOfInFlow    int64 //long long<资金净流入额(*10000)
}

//获取当日资金数据（c++的redis数据）
func getFundNow(sid string) (fund, error) {
	tNow := time.Now()
	timeNow := tNow.Format("20060102")
	t, _ := strconv.ParseInt(timeNow, 10, 64)

	redisData, err := models.RedisStore.Get("hq:trade:day:" + sid)

	var kNow fund
	if err != nil || len(redisData) == 0 {
		kNow.LlHugeBuyValue = 0
		kNow.LlBigBuyValue = 0
		kNow.LlMiddleBuyValue = 0
		kNow.LlSmallBuyValue = 0
		kNow.LlHugeSellValue = 0
		kNow.LlBigSellValue = 0
		kNow.LlMiddleSellValue = 0
		kNow.LlSmallSellValue = 0
		logging.Info("==", err.Error(), t)
		return kNow, err
	} else {
		var kn FundNow
		size := binary.Size(&kn)

		k := &FundNow{}
		buffer := bytes.NewBuffer([]byte(redisData[0:size]))
		if err := binary.Read(buffer, binary.LittleEndian, k); err != nil && err != io.EOF {
			logging.Info("error===%v", err)
		}
		//kno = append(kno, k)
		logging.Info("kkk====%+v", k)

		//
		kNow.NSID = k.NSID
		kNow.LlHugeBuyValue = float64(k.LlHugeBuyValue) / 10000
		kNow.LlBigBuyValue = float64(k.LlBigBuyValue) / 10000
		kNow.LlMiddleBuyValue = float64(k.LlMiddleBuyValue) / 10000
		kNow.LlSmallBuyValue = float64(k.LlSmallBuyValue) / 10000
		kNow.LlHugeSellValue = float64(k.LlHugeSellValue) / 10000
		kNow.LlBigSellValue = float64(k.LlBigSellValue) / 10000
		kNow.LlMiddleSellValue = float64(k.LlMiddleSellValue) / 10000
		kNow.LlSmallSellValue = float64(k.LlSmallSellValue) / 10000
		return kNow, err
	}
}

//============================================
func Statement4(KeyNsid string) string {
	fundData, _ := getFundNow(KeyNsid) //当日的资金数据

	//今日主力净流入
	mainInflow := fundData.LlHugeBuyValue + fundData.LlBigBuyValue - fundData.LlHugeSellValue - fundData.LlBigSellValue
	//主力资金占比%
	var mainProPortion float64
	if fundData.LlHugeBuyValue+fundData.LlBigBuyValue+fundData.LlMiddleBuyValue+fundData.LlSmallBuyValue != 0 {
		mainProPortion = mainInflow / (fundData.LlHugeBuyValue + fundData.LlBigBuyValue + fundData.LlMiddleBuyValue + fundData.LlSmallBuyValue)
	} else {
		mainProPortion = 0
	}

	fuData, _ := fundModel.NewFundFlow().GetFundData(KeyNsid, 30)

	var fud []fund
	for _, v := range fuData {
		var fun fund
		fun.LlHugeBuyValue = v.LlHugeBuyValue
		fun.LlBigBuyValue = v.LlBigBuyValue
		fun.LlMiddleBuyValue = v.LlMiddleBuyValue
		fun.LlSmallBuyValue = v.LlSmallBuyValue
		fun.LlHugeSellValue = v.LlHugeSellValue
		fun.LlBigSellValue = v.LlBigSellValue
		fun.LlMiddleSellValue = v.LlMiddleSellValue
		fun.LlSmallSellValue = v.LlSmallSellValue
		fud = append(fud, fun)
	}

	var funds []fund
	funds = append(funds, fundData)
	funds = append(funds, fud...)

	var cddd, lr, pDegree float64
	for i := 0; i < 5; i++ {
		cddd += Abs(funds[i].LlHugeBuyValue-funds[i].LlHugeSellValue) + Abs(funds[i].LlBigBuyValue-funds[i].LlBigSellValue)
		lr += (funds[i].LlHugeBuyValue + funds[i].LlBigBuyValue + funds[i].LlMiddleBuyValue + funds[i].LlSmallBuyValue)
	}
	if lr == 0 {
		pDegree = 0
	} else {
		pDegree = 200 * cddd / lr
	}
	var T string
	if pDegree >= 0 && pDegree < 20 {
		T = "力没有控盘，筹码分布非常分散"
	} else if pDegree >= 20 && pDegree < 60 {
		T = "主力轻度控盘，筹码分布较为分散"
	} else if pDegree >= 60 && pDegree < 120 {
		T = "主力中度控盘，且筹码比较集中"
	} else if pDegree >= 120 && pDegree < 200 {
		T = "主力高度控盘，且筹码非常集中"
	} else {
		T = "主力高度控盘，且筹码非常集中"
	}

	//股票最近三条主力净流入数据，T1,T2,T3
	T1 := fuData[0].LlHugeBuyValue + fuData[0].LlBigBuyValue - fuData[0].LlHugeSellValue - fuData[0].LlBigSellValue
	T2 := fuData[1].LlHugeBuyValue + fuData[1].LlBigBuyValue - fuData[1].LlHugeSellValue - fuData[1].LlBigSellValue
	T3 := fuData[2].LlHugeBuyValue + fuData[2].LlBigBuyValue - fuData[2].LlHugeSellValue - fuData[2].LlBigSellValue
	var S string
	if T1 > 0 && T2 > 0 && T3 > 0 {
		S = "连续3日被主力资金增仓"
	} else if T1 > 0 && T2 > 0 && T3 <= 0 {
		S = "连续2日被主力资金增仓"
	} else if T1 < 0 && T2 < 0 && T3 < 0 {
		S = "连续3日被主力资金减仓"
	} else if T1 < 0 && T2 < 0 && T3 > 0 {
		S = "连续2日被主力资金减仓"
	} else {
		S = "无连续增减仓现象，主力趋势不明显"
	}
	text_s := "该股当前" + S + "。"

	var statement4_str, ty, typ string
	if mainInflow > 0 {
		ty = "up"
	} else {
		ty = "down"
	}
	if mainProPortion > 0 {
		typ = "up"
	} else {
		typ = "down"
	}
	statement4_str = "该股今日主力净流入<span class ='" + ty + "'>" + strconv.FormatFloat(mainInflow/10000, 'f', -1, 64) + "</span>万元，主力资金占比<span class ='" + typ + "'>" + strconv.FormatFloat(mainProPortion, 'f', -1, 64) + "</span>%，" + text_s + "" + T + "。"
	return statement4_str
}
