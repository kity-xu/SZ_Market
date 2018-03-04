package diaService

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"

	"sort"
	"strconv"

	"haina.com/market/f9/models/valueModel"
	"haina.com/market/f9/services"
	"haina.com/market/f9/services/fundService"
	"haina.com/share/logging"
)

type ResponseDiagnose struct {
	Statement1  string  `json:"statement1"`
	Statement2  string  `json:"statement2"`
	Statement3  string  `json:"statement3"`
	Statement4  string  `json:"statement4"`
	Statement5  string  `json:"statement5"`
	Prompt      string  `json:"prompt"`
	Score       int64   `json:"score"`
	Symbol      string  `json:"symbol"`
	Stname      string  `json:"stname"`
	NLastPx     float64 `json:"nLastPx"`
	NPxChg      float64 `json:"nPxChg"`
	NPxChgRatio float64 `json:"nPxChgRatio"`
}

func GetDiaData(sid string) (interface{}, error) {
	var di ResponseDiagnose
	basic, _ := commonService.GetCommonData(sid)

	dataApi, err := GetApi(basic.SymbolParam)
	if err != nil {
		logging.Info("GetApi -- error |%v", err)
	}

	statement3, value, nLastPx, nPxChg, nPxChgRatio := Statement3(sid)
	statement5_s1, statement5_s2 := Statement5(basic.Compcode)
	W, X, Y, Z, Score := diagnose_part(statement5_s1, statement5_s2, value)
	logging.Info("statement5_s1=====", statement5_s1, statement5_s2)

	di.Symbol = basic.Symbol
	di.Stname = basic.Compname
	di.NLastPx = nLastPx / 10000
	di.NPxChg = nPxChg / 10000
	di.NPxChgRatio = nPxChgRatio / 10000
	di.Score = Score
	di.Prompt = Y + "，" + Z

	di.Statement1 = dataApi.Result.Data.Data.Statement1
	di.Statement2 = dataApi.Result.Data.Data.Statement2
	di.Statement3 = statement3
	di.Statement4 = fundService.Statement4(sid)
	di.Statement5 = "该股属于" + basic.Swlevelname + "行业，最近1年该股" + W + "，" + X + "。"

	return &di, nil
}

type result struct {
	Result resultData `jsn"result"`
}
type resultData struct {
	Status struct {
		Code int64 `json:"code"`
	} `json:"status"`
	Data struct {
		Code int64  `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Statement1 string `json:"statement1"`
			Statement2 string `json:"statement2"`
			Statement3 string `json:"statement3"`
			Statement4 string `json:"statement4"`
			Statement5 string `json:"statement5"`
			Prompt     string `json:"prompt"`
		} `json:"data"`
	} `json:"data"`
}

func GetApi(symbol string) (result, error) {
	var url string = "http://touzi.sina.com.cn/api/openapi.php/PerspectiveSzService.diagnose?symbol=" + symbol
	resp, err := http.Get(url)
	if err != nil {
		logging.Error(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.Error(err.Error())
	}
	logging.Info(string(body))

	var _param result
	err = json.Unmarshal([]byte(body), &_param)
	return _param, err
}

func Statement3(KeyNsid string) (string, int64, float64, float64, float64) {
	var nLastPx, nPxChg, nPxChgRatio float64
	klistNow, _ := commonService.GetKlineNow(KeyNsid)      //当日K线
	klistPastAsc, _ := commonService.GetKlinePast(KeyNsid) //历史K线

	logging.Info("klistPastAsc===%+v", klistPastAsc)
	nLastPx = klistNow.NLastPx
	nPxChg = klistNow.NPxChg
	nPxChgRatio = klistNow.NPxChgRatio
	var klistPastDesc []commonService.Klist
	for i := len(klistPastAsc) - 1; i >= 0; i-- {
		var kp commonService.Klist
		kp.Date = klistPastAsc[i].NTime
		kp.NPreCPx = klistPastAsc[i].NPreCPx
		kp.NOpenPx = klistPastAsc[i].NOpenPx
		kp.NHighPx = klistPastAsc[i].NHighPx
		kp.NLastPx = klistPastAsc[i].NLastPx
		kp.LlVolume = klistPastAsc[i].LlVolume
		klistPastDesc = append(klistPastDesc, kp)
	}

	var klists []commonService.Klist //合并后的总K线

	klists = append(klists, klistNow)
	klists = append(klists, klistPastDesc...)

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
	var m int64
	if n == 1 && MA5 >= MA10 && MA3 >= MA3_3 {
		public_M = "且上涨趋势有所加快，建议强烈关注。"
		m = 1
	}
	if n == 1 && MA5 >= MA10 && MA3 < MA3_3 {
		public_M = "但上涨趋势有所减缓，可考虑波段操作。"
		m = 2
	}
	if n == 1 && MA5 < MA10 && MA3 < MA3_3 {
		public_M = "但小趋势回调，且回调有所加快，建议暂时观望。"
		m = 3
	}
	if n == 1 && MA5 < MA10 && MA3 >= MA3_3 {
		public_M = "小趋势回调，但回调有所减缓，可考虑波段操作。"
		m = 2
	}
	if n == 0 && MA5 < MA10 && MA3 < MA3_3 {
		public_M = "且下行趋势有所加快，建议观望或减仓。"
		m = 4
	}
	if n == 0 && MA5 < MA10 && MA3 >= MA3_3 {
		public_M = "但下行趋势有所减缓，建议持币观望。"
		m = 3
	}
	if n == 0 && MA5 >= MA10 && MA3 < MA3_3 {
		public_M = "小趋势反弹，但反弹强度有所降低，建议谨慎操作。"
		m = 3
	}
	if n == 0 && MA5 >= MA10 && MA3 >= MA3_3 {
		public_M = "小趋势反弹，且反弹程度有所加强，建议抄底关注。"
		m = 2
	}

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

	var result string

	result = "该股阻力位" + strconv.FormatFloat(pressure_level, 'f', -1, 64) + "，支撑位" + strconv.FormatFloat(support_level, 'f', -1, 64) + "，" + public_R + "该股" + public_L + "，" + public_M + "。"

	return result, m, nLastPx, nPxChg, nPxChgRatio
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

func boll(Klist []commonService.Klist, n int) float64 {
	var MA_CLOSE_n_3, MA_CLOSE_n_6, MA_CLOSE_n_12, MA_CLOSE_n_24 float64
	MA_CLOSE_n_3 = sumOfFieldByNumBoll(Klist, n, n+3) / 3
	MA_CLOSE_n_6 = sumOfFieldByNumBoll(Klist, n, n+6) / 6
	MA_CLOSE_n_12 = sumOfFieldByNumBoll(Klist, n, n+12) / 12
	MA_CLOSE_n_24 = sumOfFieldByNumBoll(Klist, n, n+24) / 24
	return (MA_CLOSE_n_3 + MA_CLOSE_n_6 + MA_CLOSE_n_12 + MA_CLOSE_n_24) / 4
}
func sumOfFieldByNumBoll(Klist []commonService.Klist, start int, n int) float64 {
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

//计算成交量乘以收盘价的之和
//从一个数组中计算指定开始日期向后推n天的数据总和，如果不够n天，有几天算几天 指定开始日期从0开始,
func sumOfday(Klist []commonService.Klist, start int, n int) float64 {
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
func sumOfFieldByNum(Klist []commonService.Klist, start int, n int, feild string) float64 {
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

//=================================

type GrowData struct {
	COMPCODE int64 `json:"COMPCODE"` //公司内码
	//ENDDATE        int64   `json:"ENDDATE"`        //截止日期
	//PUBLISHDATE    int64   `json:"PUBLISHDATE"`    //发版日期
	BIZTOTINCO float64 `json:"BIZTOTINCO"` //营业收入
	PARENETP   float64 `json:"PARENETP"`   //归属于母公司所有者的净利润
	//REPORTDATETYPE int64   `json:"REPORTDATETYPE"` //报告期类型
}

type GrowData1Sort []GrowData
type GrowData2Sort []GrowData

func Statement5(compcode string) (float64, float64) {
	data, _ := valueModel.NewGrow().GetGrowChartData(compcode)
	BIZTOTINCO_tmp := make([]float64, 4)
	PARENETP_tmp := make([]float64, 4)
	BIZTOTINCO := make([]float64, 4)
	PARENETP := make([]float64, 4)
	for key, val := range data {

		if key < 4 {
			var BIZTOTINCO_t, PARENETP_t float64
			if val.REPORTDATETYPE.Int64 == 1 { //判断季度并相减
				BIZTOTINCO_t = val.BIZTOTINCO.Float64
				PARENETP_t = val.PARENETP.Float64
			} else {
				BIZTOTINCO_t = data[key].BIZTOTINCO.Float64 - data[key+1].BIZTOTINCO.Float64
				PARENETP_t = data[key].PARENETP.Float64 - data[key+1].PARENETP.Float64
			}
			BIZTOTINCO_tmp[key] = BIZTOTINCO_t
			PARENETP_tmp[key] = PARENETP_t
			BIZTOTINCO[3-key] = BIZTOTINCO_t
			PARENETP[3-key] = PARENETP_t
		} else {
			break
		}
	}

	sort.Float64s(BIZTOTINCO_tmp)
	sort.Float64s(PARENETP_tmp)
	x := []float64{1, 2, 3, 4}
	BIZTOTINCO_y := make([]float64, 4)
	PARENETP_y := make([]float64, 4)
	for i := 0; i < 4; i++ {
		BIZTOTINCO_y[i] = (BIZTOTINCO[i] - BIZTOTINCO_tmp[0]) / (BIZTOTINCO_tmp[3] - BIZTOTINCO_tmp[0])
		PARENETP_y[i] = (PARENETP[i] - PARENETP_tmp[0]) / (PARENETP_tmp[3] - PARENETP_tmp[0])
	}

	BIZTOTINCO_s1 := 1*BIZTOTINCO_y[0] + 2*BIZTOTINCO_y[1] + 3*BIZTOTINCO_y[2] + 4*BIZTOTINCO_y[3]
	BIZTOTINCO_s2 := (1 / 4) * FloatsSum(x) * FloatsSum(BIZTOTINCO_y)
	BIZTOTINCO_s3 := x[0]*x[0] + x[1]*x[1] + x[2]*x[2] + x[3]*x[3]
	BIZTOTINCO_s4 := (1 / 4) * FloatsSum(x) * FloatsSum(x)

	PARENETP_s1 := 1*PARENETP_y[0] + 2*PARENETP_y[1] + 3*PARENETP_y[2] + 4*PARENETP_y[3]
	PARENETP_s2 := (1 / 4) * FloatsSum(x) * FloatsSum(PARENETP_y)
	PARENETP_s3 := x[0]*x[0] + x[1]*x[1] + x[2]*x[2] + x[3]*x[3]
	PARENETP_s4 := (1 / 4) * FloatsSum(x) * FloatsSum(x)

	S1 := (BIZTOTINCO_s1 - BIZTOTINCO_s2) / (BIZTOTINCO_s3 - BIZTOTINCO_s4)
	S2 := (PARENETP_s1 - PARENETP_s2) / (PARENETP_s3 - PARENETP_s4)
	return S1, S2
}

//浮点型切片求和
func FloatsSum(flo []float64) float64 {
	var sum float64 = 0
	for _, v := range flo {
		sum += v
	}
	return sum
}
func (I GrowData1Sort) Len() int {
	return len(I)
}
func (I GrowData1Sort) Less(i, j int) bool {
	return I[i].BIZTOTINCO > I[j].BIZTOTINCO
}
func (I GrowData1Sort) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

func (I GrowData2Sort) Len() int {
	return len(I)
}
func (I GrowData2Sort) Less(i, j int) bool {
	return I[i].PARENETP > I[j].PARENETP
}
func (I GrowData2Sort) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

//===========================
//概述部分内容
func diagnose_part(S1 float64, S2 float64, value int64) (string, string, string, string, int64) {
	var W, X, Y, Z string
	var score int64
	if S1 > 0 && S1 < 0.3 {
		W = "主营收入平稳增长"
		if S2 < 0 {
			X = "但成本控制较差,业绩一般"
			Y = "业绩一般"
			switch value {
			case 1:
				Z = "但短期走势加强，可考虑低吸"
				score = 45
				break
			case 2:
				Z = "走势一般，建议趋势明朗后进行交易"
				score = 44
				break
			case 3:
				Z = "趋势疲软，建议继续观望"
				score = 43
				break
			case 4:
				Z = "走势趋弱，短期慎入"
				score = 42
			}
		}
		if S2 > 0 && S2 <= 0.3 {
			X = "公司总体资产亦较为平稳，但增速较低"
			Y = "业绩平稳"
			switch value {
			case 1:
				Z = "且短期走势加强，可继续持有或买入"
				score = 75 //50
				break
			case 2:
				Z = "走势一般，建议考虑波段操作"
				score = 49
				break
			case 3:
				Z = "走势较弱，短期需持币观望"
				score = 48
				break
			case 4:
				Z = "但不是短线交易时机"
				score = 47
			}
		}
		if S2 > 0.3 && S2 <= 0.7 {
			X = "且成本控制较好，净资产增长较快"
			Y = "业绩不错"
			switch value {
			case 1:
				Z = "走势较强，短期回档后应考虑低吸"
				score = 60 //55
				break
			case 2:
				Z = "处于盘整震荡期，可继续持有"
				score = 65 //54
				break
			case 3:
				Z = "趋势疲软，建议继续观望"
				score = 53
				break
			case 4:
				Z = "但不是短线交易时间"
				score = 52
			}
		}
		if S2 > 0.7 {
			X = "总资产增长迅速，看好长期前景"
			Y = "业绩不错"
			switch value {
			case 1:
				Z = "且处于上升形态，建议继续持有或买入"
				score = 75 //60
				break
			case 2:
				Z = "可长期关注"
				score = 80 //59
				break
			case 3:
				Z = "走势较弱，可继续关注"
				score = 58
				break
			case 4:
				Z = "长期可关注，但不是短期交易时机"
				score = 57
			}
		}
	}
	if S1 > 0.3 && S1 <= 0.7 {
		W = "主营收入增长良好"
		if S2 < 0 {
			X = "但成本控制较差,业绩一般"
			Y = "业绩平淡"
			switch value {
			case 1:
				Z = "但短期走势加强，可考虑低吸"
				score = 65
				break
			case 2:
				Z = "走势平稳，可继续持有"
				score = 67
				break
			case 3:
				Z = "趋势渐弱，中短期难有上佳表现"
				score = 63
				break
			case 4:
				Z = "且短期处于空头趋势，建议冲高出局换股"
				score = 62
			}
		}
		if S2 > 0 && S2 <= 0.3 {
			X = "净资产稳定增长，预期未来增长将延续"
			Y = "具备中长期投资价值"
			switch value {
			case 1:
				Z = "且处于强势多头阶段，短线回档时可考虑低吸布局"
				score = 70
				break
			case 2:
				Z = "走势平稳，可长期关注"
				score = 85
				break
			case 3:
				Z = "中期趋势一般，建议继续观望"
				score = 68
				break
			case 4:
				Z = "但不是短线交易时机"
				score = 67
			}
		}
		if S2 > 0.3 && S2 <= 0.7 {
			X = "且净资产增长显著，看好长期前景"
			Y = "看好长期前景"
			switch value {
			case 1:
				Z = "且走势强劲，短线回档时可考虑低吸布局"
				score = 75
				break
			case 2:
				Z = "走势平稳，可继续持有"
				score = 74
				break
			case 3:
				Z = "但短期量能不足，建议趋势明朗后进行交易"
				score = 63
				break
			case 4:
				Z = "但短期不具备交易机会"
				score = 60
			}
		}
		if S2 > 0.7 {
			X = "净利润迅猛增长，总资产有迅速扩张之势，建议关注长期价值投资"
			Y = "基本面非常优秀"
			switch value {
			case 1:
				Z = "且处于强势多头阶段，短线回档时可考虑低吸布局"
				score = 80
				break
			case 2:
				Z = "走势平稳，可继续持有"
				score = 79
				break
			case 3:
				Z = "但短期量能不足，建议趋势明朗后进行交易"
				score = 62
				break
			case 4:
				Z = "但短期不具备交易机会"
				score = 63
			}
		}
	}
	if S1 > 0.7 && S1 <= 1 {
		W = "营收超高速增长"
		if S2 < 0 {
			X = "但成本控制较差,成长性一般"
			Y = "业绩一般"
			switch value {
			case 1:
				Z = "但趋势增强，投机者可适当关注"
				score = 85
				break
			case 2:
				Z = "走势平稳，可暂时持有，注意风险"
				score = 70 //84
				break
			case 3:
				Z = "处于盘整调整期，建议观望"
				score = 47 //83
				break
			case 4:
				Z = "处于下降阶段，短期慎入"
				score = 43 //82
			}
		}
		if S2 > 0 && S2 <= 0.3 {
			X = "但净利润增长相对缓慢，需要重点关注成本控制模式"
			Y = "可长期关注"
			switch value {
			case 1:
				Z = "但短期走势加强，可考虑低吸"
				score = 90
				break
			case 2:
				Z = "走势平稳，可继续持有"
				score = 89
				break
			case 3:
				Z = "后市不明朗，可长期跟踪其业绩变化"
				score = 60 //88
				break
			case 4:
				Z = "且短期处于空头趋势，建议冲高出局换股"
				score = 58 //87
			}
		}
		if S2 > 0.3 && S2 <= 0.7 {
			X = "且净利润大幅上升，资产有扩张之势，建议关注长期价值投资"
			Y = "适合价值投资"
			switch value {
			case 1:
				Z = "短线回档时可考虑低吸布局"
				score = 95
				break
			case 2:
				Z = "走势平稳，可继续持有"
				score = 94
				break
			case 3:
				Z = "但短期量能不足，建议趋势明朗后进行交易"
				score = 64 //93
				break
			case 4:
				Z = "但短期不具备交易机会"
				score = 58 //92
			}
		}
		if S2 > 0.7 {
			X = "且净利润迅猛增长，总资产扩张迅速，建议强烈关注公司发展"
			Y = "业绩喜人"
			switch value {
			case 1:
				Z = "可中长期持有"
				score = 100
				break
			case 2:
				Z = "可中长期持有"
				score = 99
				break
			case 3:
				Z = "可中长期持有"
				score = 98
				break
			case 4:
				Z = "短期需谨慎，长期可关注"
				score = 75 //97
			}
		}
	}
	if S1 > -0.1 && S1 <= 0 {
		W = "营收增长有限"
		if S2 < 0 {
			X = "但成本控制较差,成长性一般"
			Y = "业绩疲软"
			switch value {
			case 1:
				Z = "走势较强，可考虑波段操作"
				score = 67 //25
				break
			case 2:
				Z = "走势平淡"
				score = 24
				break
			case 3:
				Z = "短期需观望"
				score = 23
				break
			case 4:
				Z = "空头趋势，建议调仓换股"
				score = 22
			}
		}
		if S2 > 0 && S2 <= 0.3 {
			X = "总资产增长平稳"
			Y = "基本面一般"
			switch value {
			case 1:
				Z = "走势较强，可考虑波段操作"
				score = 68 //30
				break
			case 2:
				Z = "走势平淡"
				score = 29
				break
			case 3:
				Z = "短期需观望"
				score = 28
				break
			case 4:
				Z = "空头趋势，建议调仓换股"
				score = 27
			}
		}
		if S2 > 0.3 && S2 <= 0.7 {
			X = "但成本控制较好，净资产增长迅速，但或后续增长乏力"
			Y = "成长性一般"
			switch value {
			case 1:
				Z = "走势较强，可考虑波段操作"
				score = 55 //35
				break
			case 2:
				Z = "走势平淡"
				score = 34
				break
			case 3:
				Z = "短期需观望"
				score = 33
				break
			case 4:
				Z = "空头趋势，建议调仓换股"
				score = 32
			}
		}
		if S2 > 0.7 {
			X = "净资产快速增长，建议重点关注公司负债变化"
			Y = "业绩一般"
			switch value {
			case 1:
				Z = "走势较强，可考虑波段操作"
				score = 40
				break
			case 2:
				Z = "走势平淡"
				score = 39
				break
			case 3:
				Z = "短期需观望"
				score = 38
				break
			case 4:
				Z = "空头趋势，建议调仓换股"
				score = 37
			}
		}
	}
	if S1 > -0.3 && S1 <= -0.1 {
		W = "业绩下降"
	}
	if S1 > -0.7 && S1 <= -0.3 {
		W = "营业收入连续负增长"
	}
	if S1 <= -0.7 {
		W = "业绩剧烈下降"
	}
	if (S1 > -0.3 && S1 <= -0.1) || (S1 > -0.7 && S1 <= -0.3) || (S1 <= -0.7) {
		if S2 < 0 {
			X = "且盈利能力较低，不具备长期投资价值"
			Y = "基本面差"
			switch value {
			case 1:
				Z = "走势较强，可考虑波段操作"
				score = 36 //5
				break
			case 2:
				Z = "走势平淡"
				score = 4
				break
			case 3:
				Z = "短期需观望"
				score = 3
				break
			case 4:
				Z = "空头趋势，建议调仓换股"
				score = 2
			}
		}
		if S2 > 0 && S2 <= 0.3 {
			X = "净资产有所增长，但增长缓慢"
			Y = "不具备长期投资价值"
			switch value {
			case 1:
				Z = "走势较强，可考虑波段操作"
				score = 42 //10
				break
			case 2:
				Z = "走势平淡"
				score = 9
				break
			case 3:
				Z = "短期需观望"
				score = 8
				break
			case 4:
				Z = "空头趋势，建议调仓换股"
				score = 7
			}
		}
		if S2 > 0.3 && S2 <= 0.7 {
			X = "虽然总资产平稳增长，成本控制能力较强，但成长方面存在一定隐患"
			Y = "缺乏成长性"
			switch value {
			case 1:
				Z = "走势较强，可考虑波段操作，注意风险"
				score = 39 //15
				break
			case 2:
				Z = "走势平淡"
				score = 14
				break
			case 3:
				Z = "短期需观望"
				score = 13
				break
			case 4:
				Z = "空头趋势，建议调仓换股"
				score = 12
			}
		}
		if S2 > 0.7 {
			X = "虽然净资产快速增长，但成长性较差"
			Y = "不具备长期投资价值"
			switch value {
			case 1:
				Z = "走势较强，可考虑波段操作，注意风险"
				score = 37 //20
				break
			case 2:
				Z = "走势平淡"
				score = 19
				break
			case 3:
				Z = "短期需观望"
				score = 18
				break
			case 4:
				Z = "空头趋势，建议调仓换股"
				score = 17
			}
		}
	}
	return W, X, Y, Z, score
}
