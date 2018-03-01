package chipService

import (
	"math"
	"strconv"

	"ProtocolBuffer/projects/f9/go/protocol"

	"haina.com/share/logging"
)

type Trend struct {
	Trend_price    float64          `json:"chip_price"`
	Tendency_text  string           `json:"tendency_text"`
	Tendency_chart []protocol.KInfo `json:"tendency_chart"`
	Zcw            float64          `json:"zcw"` //支撑位
	Zlm            float64          `json:"zlm"` //压力位 阻力位
	Point          float64          `json:"point"`
	Point_text     string           `json:"point_text"`
}

func GetTrendData(sid string) (*Trend, error) {
	SID, _ := strconv.Atoi(sid)
	klistNow, err := getKlineNow(SID) //当日K线
	if err != nil {
		logging.Error("The dayLine is null")
		return nil, err
	}

	klistPast, err := getKlinePast(SID) //历史K线
	if err != nil {
		logging.Error("The history dayLine is null")
		return nil, err
	}

	l := len(*klistPast)
	var ktable []protocol.KInfo //合并后的总K线
	if l > 1 && (*klistPast)[l-1].LlVolume == 0 {
		ktable = append(ktable, (*klistPast)[0:l-2]...)
	}

	market, err := GetMarketStatus(100000000)
	if err != nil {
		return nil, err
	}
	lt := len(ktable)
	if lt == 0 {
		ktable = append(ktable, *klistNow)
	} else {
		if market.NTradeDate > ktable[lt-1].NTime {
			ktable = append(ktable, *klistNow) //合并后的总K线
		}
	}

	var public_L, public_M string
	var n int64
	if ktable[0].NLastPx > 0 && float64(ktable[0].NLastPx/10000) > sumOfFieldByNum(&ktable, 0, 10, "nLastPx") {
		public_L = "多头"
	} else {
		public_L = "空头"
	}

	var MA5, MA10, MA3, MA3_3 float64
	MA5 = sumOfFieldByNum(&ktable, 0, 5, "nLastPx") / 5
	MA10 = sumOfFieldByNum(&ktable, 0, 10, "nLastPx") / 10
	MA3 = sumOfFieldByNum(&ktable, 0, 3, "nLastPx") / 3
	MA3_3 = sumOfFieldByNum(&ktable, 3, 6, "nLastPx") / 3
	if n == 1 && MA5 >= MA10 && MA3 >= MA3_3 {
		public_M = "且上涨趋势有所加快，建议强烈关注。"
	}
	if n == 1 && MA5 >= MA10 && MA3 < MA3_3 {
		public_M = "但上涨趋势有所减缓，可考虑波段操作。"
	}
	if n == 1 && MA5 < MA10 && MA3 < MA3_3 {
		public_M = "但小趋势回调，且回调有所加快，建议暂时观望。"
	}
	if n == 1 && MA5 < MA10 && MA3 >= MA3_3 {
		public_M = "小趋势回调，但回调有所减缓，可考虑波段操作。"
	}
	if n == 0 && MA5 < MA10 && MA3 < MA3_3 {
		public_M = "且下行趋势有所加快，建议观望或减仓。"
	}
	if n == 0 && MA5 < MA10 && MA3 >= MA3_3 {
		public_M = "但下行趋势有所减缓，建议持币观望。"
	}
	if n == 0 && MA5 >= MA10 && MA3 < MA3_3 {
		public_M = "小趋势反弹，但反弹强度有所降低，建议谨慎操作。"
	}
	if n == 0 && MA5 >= MA10 && MA3 >= MA3_3 {
		public_M = "小趋势反弹，且反弹程度有所加强，建议抄底关注。"
	}

	var tr Trend

	tr.Tendency_text = "该股" + public_L + "，" + public_M
	tr.Tendency_chart = ktable

	var boll_11_agv float64 = 0 //计算boll从1到11的之和

	for i := 1; i <= 11; i++ {
		boll_11_agv = boll_11_agv + boll(ktable, i)
	}
	boll_11_agv = boll_11_agv / 11

	var boll_11_nn float64 = 0 //计算平方之和

	for i := 1; i <= 11; i++ {
		boll_11_nn = boll_11_nn + (boll(ktable, i)-boll_11_agv)*(boll(ktable, i)-boll_11_agv)
	}
	var standard, pressure_level, support_level, degrees float64
	var public_R string
	standard = float64(InvSqrt(float32(boll_11_agv) / 11)) //标准差

	pressure_level = boll(ktable, 0) + 2.1*standard //压力位
	support_level = boll(ktable, 0) - 2.1*standard  //支撑位

	if ktable[0].NLastPx == 0 {
		degrees = 30 + 180*((0-support_level)/(pressure_level-support_level)) //指针度数
	} else {
		degrees = 30 + 180*((float64(ktable[0].NLastPx/10000)-support_level)/(pressure_level-support_level)) //指针度数
	}

	if degrees <= 30 {
		public_R = "目前该股已经跌破支撑位，注意风险。"
	}
	if degrees > 30 && degrees <= 64 {
		public_R = "目前该股遇到强支撑位，短期破位下行风险。"
	}
	if degrees > 64 && degrees <= 98 {
		public_R = "目前该股遇到弱支撑位,注意下行风险。"
	}
	if degrees > 98 && degrees <= 167.52 {
		public_R = "目前该股处于震荡区间内。"
	}
	if degrees > 167.52 && degrees <= 210 {
		public_R = "目前该股遇到弱阻力位,短期注意风险。"
	}
	if degrees > 210 && degrees <= 240 {
		public_R = "目前该股遇到强阻力位,短期注意获利回吐风险。"
	}
	if degrees > 240 {
		public_R = "目前该股已经突破强阻力位,短期注意冲高回落。"
	}

	logging.Info("public_R====%v", public_R)

	tr.Zcw = support_level
	tr.Zlm = pressure_level
	tr.Point = degrees
	tr.Point_text = public_R

	return &tr, nil

}

func InvSqrt(x float32) float32 {
	var xhalf float32 = 0.5 * x // get bits for floating VALUE
	i := math.Float32bits(x)    // gives initial guess y0
	i = 0x5f375a86 - (i >> 1)   // convert bits BACK to float
	x = math.Float32frombits(i) // Newton step, repeating increases accuracy
	x = x * (1.5 - xhalf*x*x)
	x = x * (1.5 - xhalf*x*x)
	x = x * (1.5 - xhalf*x*x)
	return 1 / x
}

func boll(Klist []protocol.KInfo, n int) float64 {
	var MA_CLOSE_n_3, MA_CLOSE_n_6, MA_CLOSE_n_12, MA_CLOSE_n_24 float64
	MA_CLOSE_n_3 = sumOfFieldByNumBoll(Klist, n, n+3) / 3
	MA_CLOSE_n_6 = sumOfFieldByNumBoll(Klist, n, n+6) / 6
	MA_CLOSE_n_12 = sumOfFieldByNumBoll(Klist, n, n+12) / 12
	MA_CLOSE_n_24 = sumOfFieldByNumBoll(Klist, n, n+24) / 24
	return (MA_CLOSE_n_3 + MA_CLOSE_n_6 + MA_CLOSE_n_12 + MA_CLOSE_n_24) / 4
}
func sumOfFieldByNumBoll(Klist []protocol.KInfo, start int, n int) float64 {
	if start >= len(Klist) {
		return 0
	}
	var sum float64 = 0
	for i := start; i < len(Klist); i++ {
		if i < n {
			sum = sum + float64(Klist[i].NLastPx/10000)
		} else {
			break
		}
	}
	return sum
}
