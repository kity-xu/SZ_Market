package valueService

import (
	"fmt"
	"sort"
	"strconv"

	"haina.com/market/f9/models/valueModel"
)

type Grow struct {
	Text      string       `json:"text"`
	ChartData []*growChart `json:"chartData"`
	ChartText string       `json:"chartText"`
}

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

type growChart struct {
	Date      string  `json:"date"`
	Data_line float64 `json:"data_line"`
	Data_bar  float64 `json:"data_bar"`
}

func GetGrowChartData(leng string) (Grow, error) {
	g := []*growChart{}
	data, err := valueModel.NewGrow().GetGrowChartData(compcode)
	BIZTOTINCO_tmp := make([]float64, 4)
	PARENETP_tmp := make([]float64, 4)
	BIZTOTINCO := make([]float64, 4)
	PARENETP := make([]float64, 4)
	for key, val := range data {
		if key < 8 {
			var ch growChart
			var t1 string = strconv.FormatInt(data[key].ENDDATE.Int64, 10)
			var t2 string = strconv.FormatInt(data[key+1].ENDDATE.Int64, 10)
			ch.Date = t1[0:4] + "-" + t1[4:6] + "-" + t1[6:8]
			if t1[0:4] == t2[0:4] {
				if data[key].REPORTDATETYPE.Int64 != 1 {
					ch.Data_line = data[key].PARENETP.Float64 - data[key+1].PARENETP.Float64
					ch.Data_bar = data[key].BIZTOTINCO.Float64 - data[key+1].BIZTOTINCO.Float64
				} else {
					ch.Data_line = data[key].PARENETP.Float64
					ch.Data_bar = data[key].BIZTOTINCO.Float64

				}
			} else {
				ch.Data_line = data[key].PARENETP.Float64
				ch.Data_bar = data[key].BIZTOTINCO.Float64
			}
			g = append(g, &ch)
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
			}
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

	var W, X string
	if S1 > 0 && S1 < 0.3 {
		W = "主营收入平稳增长"
		if S2 < 0 {
			X = "但成本控制较差，业绩一般"
		}
		if S2 > 0 && S2 <= 0.3 {
			X = "公司总体资产亦较为平稳，但增速较低"
		}
		if S2 > 0.3 && S2 <= 0.7 {
			X = "且成本控制较好，净资产增长较快"
		}
		if S2 > 0.7 {
			X = "总资产增长迅速，看好长期前景"
		}
	}
	if S1 > 0.3 && S1 <= 0.7 {
		W = "主营收入增长良好"
		if S2 < 0 {
			X = "但成本控制较差，业绩一般"
		}
		if S2 > 0 && S2 <= 0.3 {
			X = "净资产稳定增长，预期未来增长将延续"
		}
		if S2 > 0.3 && S2 <= 0.7 {
			X = "且净资产增长显著，看好长期前景"
		}
		if S2 > 0.7 {
			X = "净利润迅猛增长，总资产有迅速扩张之势，建议关注长期价值投资"
		}
	}
	if S1 > 0.7 && S1 <= 1 {
		W = "营收超高速增长"
		if S2 < 0 {
			X = "但成本控制较差，成长性一般"
		}
		if S2 > 0 && S2 <= 0.3 {
			X = "但净利润增长相对缓慢，需要重点关注成本控制模式"
		}
		if S2 > 0.3 && S2 <= 0.7 {
			X = "且净利润大幅上升，资产有扩张之势，建议关注长期价值投资"
		}
		if S2 > 0.7 {
			X = "且净利润迅猛增长，总资产扩张迅速，建议强烈关注公司发展"
		}
	}
	if S1 > -0.1 && S1 <= 0 {
		W = "营收增长有限"
		if S2 < 0 {
			X = "但成本控制较差，成长性一般"
		}
		if S2 > 0 && S2 <= 0.3 {
			X = "总资产增长平稳"
		}
		if S2 > 0.3 && S2 <= 0.7 {
			X = "但成本控制较好，净资产增长迅速，但或后续增长乏力"
		}
		if S2 > 0.7 {
			X = "净资产快速增长，建议重点关注公司负债变化"
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
		}
		if S2 > 0 && S2 <= 0.3 {
			X = "净资产有所增长，但增长缓慢"
		}
		if S2 > 0.3 && S2 <= 0.7 {
			X = "虽然总资产平稳增长，成本控制能力较强，但成长方面存在一定隐患"
		}
		if S2 > 0.7 {
			X = "虽然净资产快速增长，但成长性较差"
		}
	}
	fmt.Println(S1, "=================", S2)
	fmt.Println(W, "=================", X)
	textData, err := valueModel.NewGrow().GetGrowTextData(swlevelcode)
	var gd []GrowData
	for i := 0; i < len(textData)-1; i++ {
		if textData[i].COMPCODE.Int64 == textData[i+1].COMPCODE.Int64 {
			if textData[i].REPORTDATETYPE.Int64 != 1 {
				textData[i].BIZTOTINCO.Float64 = textData[i].BIZTOTINCO.Float64 - textData[i+1].BIZTOTINCO.Float64
				textData[i].PARENETP.Float64 = textData[i].PARENETP.Float64 - textData[i+1].PARENETP.Float64
			}
			var t bool = true
			if len(gd) > 0 {
				for j := 0; j < len(gd); j++ {
					if textData[i].COMPCODE.Int64 == gd[j].COMPCODE {
						t = false
						break
					}
				}
			}
			if t {
				var g GrowData
				g.COMPCODE = textData[i].COMPCODE.Int64
				g.BIZTOTINCO = textData[i].BIZTOTINCO.Float64
				g.PARENETP = textData[i].PARENETP.Float64
				gd = append(gd, g)
			}
			i++
		}
	}
	text1 := []GrowData{}
	text2 := []GrowData{}
	for _, val := range gd {
		var t1 GrowData
		t1.COMPCODE = val.COMPCODE
		t1.BIZTOTINCO = val.BIZTOTINCO
		t1.PARENETP = val.PARENETP
		text1 = append(text1, t1)
		var t2 GrowData
		t2.COMPCODE = val.COMPCODE
		t2.BIZTOTINCO = val.BIZTOTINCO
		t2.PARENETP = val.PARENETP
		text2 = append(text2, t2)
	}
	sort.Sort(GrowData1Sort(text1))
	sort.Sort(GrowData2Sort(text2))
	var BIZTOTINCO_key, PARENETP_key int
	for key, val := range text1 {
		com, _ := strconv.ParseInt(compcode, 10, 64) //转为变量类型为整型 然后进行比较
		if com == val.COMPCODE {
			BIZTOTINCO_key = key + 1
			break
		}
	}
	for key, val := range text2 {
		com, _ := strconv.ParseInt(compcode, 10, 64) //转为变量类型为整型 然后进行比较
		if com == val.COMPCODE {
			PARENETP_key = key + 1
			break
		}
	}
	fmt.Println(BIZTOTINCO_key, "=================", PARENETP_key)
	fmt.Println(gd, "=================", X)
	fmt.Println("%+v", gd)
	//logging.Info(BIZTOTINCO[0])u9
	//logging.Info(PARENETP)
	var textString string = "最近一年，" + compname + W + "，" + X + "，指标分别在" + swlevelname + "行业排名 " + strconv.Itoa(BIZTOTINCO_key) + "/" + leng + "，" + strconv.Itoa(PARENETP_key) + "/" + leng

	var y Grow
	y.ChartData = g
	y.ChartText = textString
	return y, err
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
