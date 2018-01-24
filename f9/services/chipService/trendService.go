package chipService

import (
	"math"

	"haina.com/share/logging"
)

type Trend struct {
	Trend_price    float64 `json:"chip_price"`
	Tendency_text  string  `json:"tendency_text"`
	Tendency_chart []Klist `json:"tendency_chart"`
	Zcw            float64 `json:"zcw"` //支撑位
	Zlm            float64 `json:"zlm"` //压力位 阻力位
	Point          float64 `json:"point"`
	Point_text     string  `json:"point_text"`
}

func GetTrendData(scode string) (Trend, error) {
	//	var klistNow Klist //当日K线
	//	klistNow.Date = 20171222
	//	klistNow.NPreCPx = 0
	//	klistNow.NHighPx = 0
	//	klistNow.NLastPx = 0
	//	klistNow.LlVolume = 0
	//	klistNow.LlValue = 0

	klistNow, _ := getKlineNow() //当日K线

	klistPastAsc, _ := getKlinePast() //历史K线

	var klistPastDesc []Klist
	for i := len(klistPastAsc) - 1; i >= 0; i-- {
		logging.Info("k====", i)
		var kp Klist
		kp.Date = klistPastAsc[i].NTime
		kp.NPreCPx = klistPastAsc[i].NPreCPx
		kp.NOpenPx = klistPastAsc[i].NOpenPx
		kp.NHighPx = klistPastAsc[i].NHighPx
		kp.NLastPx = klistPastAsc[i].NLastPx
		kp.LlVolume = klistPastAsc[i].LlVolume
		klistPastDesc = append(klistPastDesc, kp)
	}

	var klists []Klist //合并后的总K线

	klists = append(klists, klistNow)
	klists = append(klists, klistPastDesc...)

	logging.Info("klistPastDesc%+v", klists)
	var newKlist []Klist
	for i := 0; i < len(klists); i++ {
		if klists[i].LlVolume > 0 {
			var kp Klist
			kp.Date = klists[i].Date
			if klists[i].NLastPx > 0 && klists[i].NLastPx/10000 > sumOfFieldByNum(klists, i, i+10, "nLastPx") {
				kp.Value = 1
			} else {
				kp.Value = 0
			}
			kp.NPreCPx = klists[i].NPreCPx / 10000
			kp.NOpenPx = klists[i].NOpenPx / 10000
			kp.NHighPx = klists[i].NHighPx / 10000
			kp.NLastPx = klists[i].NLastPx / 10000
			kp.LlVolume = klists[i].LlVolume
			kp.LlValue = klists[i].LlValue
			newKlist = append(newKlist, kp)
		}
	}
	logging.Info("newKlist%+v", newKlist)
	var public_L, public_M string
	var n int64
	if klists[0].NLastPx > 0 && klists[0].NLastPx/10000 > sumOfFieldByNum(klists, 0, 10, "nLastPx") {
		public_L = "多头"
	} else {
		public_L = "空头"
	}

	var MA5, MA10, MA3, MA3_3 float64
	MA5 = sumOfFieldByNum(klists, 0, 5, "nLastPx") / 5
	MA10 = sumOfFieldByNum(klists, 0, 10, "nLastPx") / 10
	MA3 = sumOfFieldByNum(klists, 0, 3, "nLastPx") / 3
	MA3_3 = sumOfFieldByNum(klists, 3, 6, "nLastPx") / 3
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
	tr.Tendency_chart = newKlist

	var boll_11_agv float64 = 0 //计算boll从1到11的之和

	for i := 1; i <= 11; i++ {
		boll_11_agv = boll_11_agv + boll(klists, i)
	}
	boll_11_agv = boll_11_agv / 11

	var boll_11_nn float64 = 0 //计算平方之和

	for i := 1; i <= 11; i++ {
		boll_11_nn = boll_11_nn + (boll(klists, i)-boll_11_agv)*(boll(klists, i)-boll_11_agv)
	}
	var standard, pressure_level, support_level, degrees float64
	var public_R string
	standard = float64(InvSqrt(float32(boll_11_agv) / 11)) //标准差
	logging.Info("boll_11_agv====%v", boll_11_agv/11)
	logging.Info("standard====%v", standard)
	pressure_level = boll(klists, 0) + 2.1*standard //压力位
	support_level = boll(klists, 0) - 2.1*standard  //支撑位
	logging.Info("pressure_level====%v", pressure_level)
	logging.Info("support_level====%v", support_level)
	if klists[0].NLastPx == 0 {
		degrees = 30 + 180*((0-support_level)/(pressure_level-support_level)) //指针度数
	} else {
		degrees = 30 + 180*(((klists[0].NLastPx/10000)-support_level)/(pressure_level-support_level)) //指针度数
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

	return tr, nil

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

func boll(Klist []Klist, n int) float64 {
	var MA_CLOSE_n_3, MA_CLOSE_n_6, MA_CLOSE_n_12, MA_CLOSE_n_24 float64
	MA_CLOSE_n_3 = sumOfFieldByNumBoll(Klist, n, n+3) / 3
	MA_CLOSE_n_6 = sumOfFieldByNumBoll(Klist, n, n+6) / 6
	MA_CLOSE_n_12 = sumOfFieldByNumBoll(Klist, n, n+12) / 12
	MA_CLOSE_n_24 = sumOfFieldByNumBoll(Klist, n, n+24) / 24
	return (MA_CLOSE_n_3 + MA_CLOSE_n_6 + MA_CLOSE_n_12 + MA_CLOSE_n_24) / 4
}
func sumOfFieldByNumBoll(Klist []Klist, start int, n int) float64 {
	if start >= len(Klist) {
		return 0
	}
	var sum float64 = 0
	for i := start; i < len(Klist); i++ {
		if i < n {
			sum = sum + (Klist[i].NLastPx / 10000)
		} else {
			break
		}
	}
	return sum
}
