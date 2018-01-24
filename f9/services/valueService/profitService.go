package valueService

import (
	"sort"
	"strconv"

	"haina.com/market/f9/models/valueModel"
)

//盈利能力开始=================================================================
type Profit struct {
	Text      string         `json:"text"`
	ChartData []*ProfitChart `json:"chartData"`
	ChartText string         `json:"chartText"`
}

type ProfitChart struct {
	Date      string  `json:"date"`
	Data_line float64 `json:"data_line"`
	Data_bar  float64 `json:"data_bar"`
}

type ProfitText struct {
	COMPCODE       int64
	ENDDATE        int64
	REPORTDATETYPE int64
	EPSBASIC       float64
	ROEWEIGHTED    float64
}
type ProfitText1Sort []ProfitText
type ProfitText2Sort []ProfitText

//盈利能力本文显示
func GetProfitTextData(ProfitChartText chan string, leng string) (chan string, error) {
	data, err := valueModel.NewProfit().GetProfitTextData(swlevelcode)

	text := []ProfitText{}
	for i := 0; i < len(data)-1; i++ {
		if data[i].COMPCODE == data[i+1].COMPCODE {
			if data[i].REPORTDATETYPE.Int64 != 1 {
				data[i].EPSBASIC.Float64 = data[i].EPSBASIC.Float64 - data[i+1].EPSBASIC.Float64
				data[i].ROEWEIGHTED.Float64 = data[i].ROEWEIGHTED.Float64 - data[i+1].ROEWEIGHTED.Float64
			}
			var t bool = true
			if len(text) != 0 {
				for j := 0; j < len(text); j++ {
					if text[j].COMPCODE == data[i].COMPCODE.Int64 {
						t = false
						break
					}
				}
			}
			if t {
				var te ProfitText
				te.COMPCODE = data[i].COMPCODE.Int64
				te.ENDDATE = data[i].ENDDATE.Int64
				te.EPSBASIC = data[i].EPSBASIC.Float64
				te.ROEWEIGHTED = data[i].ROEWEIGHTED.Float64
				text = append(text, te)
			}
		}
	}

	text1 := []ProfitText{} //按照EPSBASIC倒序排序
	text2 := []ProfitText{}
	for _, val := range text {
		var t1 ProfitText
		t1.COMPCODE = val.COMPCODE
		t1.ENDDATE = val.ENDDATE
		t1.REPORTDATETYPE = val.REPORTDATETYPE
		t1.EPSBASIC = val.EPSBASIC
		t1.ROEWEIGHTED = val.ROEWEIGHTED
		text1 = append(text1, t1)

		var t2 ProfitText
		t2.COMPCODE = val.COMPCODE
		t2.ENDDATE = val.ENDDATE
		t2.REPORTDATETYPE = val.REPORTDATETYPE
		t2.EPSBASIC = val.EPSBASIC
		t2.ROEWEIGHTED = val.ROEWEIGHTED
		text2 = append(text2, t2)
	}
	//logging.Info("排序前==============================")
	sort.Sort(ProfitText1Sort(text1))
	sort.Sort(ProfitText2Sort(text2))

	var sort1, sort2 int
	var FinancialRatios40, FinancialRatios59 float64
	var textString string
	for key, val := range text1 {
		com, _ := strconv.ParseInt(compcode, 10, 64) //转为变量类型为整型 然后进行比较
		if com == val.COMPCODE {
			//logging.Info("=====比较获得-key:%v==", key)
			sort1 = key + 1
			FinancialRatios40 = val.EPSBASIC
			FinancialRatios59 = val.ROEWEIGHTED
			break
		}
	}
	for key, val := range text2 {
		com, _ := strconv.ParseInt(compcode, 10, 64) //转为变量类型为整型 然后进行比较
		if com == val.COMPCODE {
			//logging.Info("=====比较获得-key:%v==", key)
			sort2 = key + 1
			break
		}

		//logging.Info("-----compcode=%v------", reflect.TypeOf(text1[0].COMPCODE))   //判断变量类型
	}

	if sort1 > 0 && len(text1)/3 > sort1 {
		textString = "前列，建议强烈关注"
	} else if sort1 > (len(text1)/3)*2 && sort1 < (len(text1)/3)*3 {
		textString = "靠后，注意风险"
	} else {
		textString = "中间，适当关注"
	}
	//	logging.Info("=====textString======", textString)
	//	logging.Info("=====FinancialRatios40======", strconv.FormatFloat(FinancialRatios40, 'f', -1, 64))
	//	logging.Info("=====FinancialRatios59======", FinancialRatios59)
	//	logging.Info("=====sort1======", sort1)
	//	logging.Info("=====sort2======", sort2)
	ProfitChartText <- "最近报告期，公司每股收益为 " + strconv.FormatFloat(FinancialRatios40, 'f', -1, 64) + "，净资产收益率为 " + strconv.FormatFloat(FinancialRatios59, 'f', -1, 64) +
		"，位居行业" + textString + ",指标分别在" + swlevelname + "行业中排名" + strconv.Itoa(sort1) + "/" + leng + "," + strconv.Itoa(sort2) + "/" + leng
	return ProfitChartText, err
}

//盈利能力的8条数据
func GetProfitChartData(ProfitChartData chan []*ProfitChart, scode string) (chan []*ProfitChart, error) {
	chart := []*ProfitChart{}
	data, err := valueModel.NewProfit().GetProfitChartData(scode)
	for key, _ := range data {
		//logging.Info("====%+v", val)
		var ch ProfitChart
		if key < 8 {
			var t1 string = strconv.FormatInt(data[key].ENDDATE.Int64, 10)
			var t2 string = strconv.FormatInt(data[key+1].ENDDATE.Int64, 10)
			//logging.Info("key==%v", key)
			ch.Date = t1[0:4] + "-" + t1[4:6] + "-" + t1[6:8]
			//ch.Date = t1[0:4]
			if t1[0:4] == t2[0:4] {
				if data[key].REPORTDATETYPE.Int64 != 1 {
					ch.Data_line = (data[key].ROEWEIGHTED.Float64*10000 - data[key+1].ROEWEIGHTED.Float64*10000) / 10000
					ch.Data_bar = (data[key].EPSBASIC.Float64*10000 - data[key+1].EPSBASIC.Float64*10000) / 10000
				} else {
					ch.Data_line = data[key].ROEWEIGHTED.Float64
					ch.Data_bar = data[key].EPSBASIC.Float64

				}
			} else {
				ch.Data_line = data[key].ROEWEIGHTED.Float64
				ch.Data_bar = data[key].EPSBASIC.Float64
			}
			//logging.Info("===%+v", ch)
			chart = append(chart, &ch)
		}
	}
	ProfitChartData <- chart
	return ProfitChartData, err
}

func (I ProfitText1Sort) Len() int {
	return len(I)
}
func (I ProfitText1Sort) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}
func (I ProfitText1Sort) Less(i, j int) bool {
	return I[i].EPSBASIC > I[j].EPSBASIC
}

func (I ProfitText2Sort) Len() int {
	return len(I)
}
func (I ProfitText2Sort) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}
func (I ProfitText2Sort) Less(i, j int) bool {
	return I[i].ROEWEIGHTED > I[j].ROEWEIGHTED
}

//盈利能力开始=================================================================
