package valueService

import (
	"sort"
	"strconv"

	"haina.com/market/f9/models/finchina"
	"haina.com/market/f9/models/valueModel"
)

type pay struct {
	COMPCODE  int64
	ENDDATE   int64
	ASSLIABRT float64
	QUICKRT   float64
}
type paySort1 []pay
type paySort2 []pay
type payChart struct {
	Date      string  `json:"date"`
	Data_line float64 `'json:"data_line"`
	Data_bar  float64 `json:"data_bar"`
}

var ASSLIABRT, QUICKRT chan float64

func init() {
	ASSLIABRT = make(chan float64)
	QUICKRT = make(chan float64)
}

//获取8条数据
func GetPayChartData(payChartData chan []*payChart, compcode string) (chan []*payChart, error) {
	data := []*payChart{}
	chartData, err := valueModel.NewPay().GetPayChartData(compcode)
	for _, val := range chartData {
		var d payChart
		d.Date = strconv.FormatInt(val.ENDDATE.Int64, 10)
		d.Data_line = val.QUICKRT.Float64
		d.Data_bar = val.ASSLIABRT.Float64
		data = append(data, &d)
	}
	ASSLIABRT <- chartData[0].ASSLIABRT.Float64
	QUICKRT <- chartData[0].QUICKRT.Float64
	//fmt.Println(ASSLIABRT, "====", QUICKRT)
	payChartData <- data
	return payChartData, err
}

//获取文本
func GetPayChartText(payChartDext chan string, detail finchina.CompanyDetail, leng string) (chan string, error) {
	var text1, text2 []pay
	textData, err := valueModel.NewPay().GetPayChartText(detail.SWLEVEL1CODE)
	//fmt.Println(err)
	for _, val := range textData {
		var t pay
		t.COMPCODE = val.COMPCODE.Int64
		t.ASSLIABRT = val.ASSLIABRT.Float64
		t.QUICKRT = val.QUICKRT.Float64
		text1 = append(text1, t)
		text2 = append(text2, t)
	}
	sort.Sort(paySort1(text1))
	sort.Sort(paySort2(text2))
	var ASSLIABRT_key, QUICKRT_key int
	for key, val := range text1 {
		com, _ := strconv.ParseInt(detail.COMPCODE, 10, 64) //string转int64
		if com == val.COMPCODE {
			ASSLIABRT_key = key + 1
		}
	}
	for key, val := range text2 {
		com, _ := strconv.ParseInt(detail.COMPCODE, 10, 64) //string转int64     strconv.Itoa int转string
		if com == val.COMPCODE {
			QUICKRT_key = key + 1
		}
	}

	payChartDext <- "最近报告期，公司资产负债率为" + strconv.FormatFloat(<-ASSLIABRT, 'f', -1, 64) + "，速动比率为" + strconv.FormatFloat(<-QUICKRT, 'f', -1, 64) + "，两个指标分别在" + detail.SWLEVEL1NAME + "行业排名 " + strconv.Itoa(ASSLIABRT_key) + "/" + leng + "，" + strconv.Itoa(QUICKRT_key) + "/" + leng
	return payChartDext, err
}

func (I paySort1) Len() int {
	return len(I)
}
func (I paySort1) Less(i, j int) bool {
	return I[i].ASSLIABRT > I[j].ASSLIABRT
}
func (I paySort1) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

func (I paySort2) Len() int {
	return len(I)
}
func (I paySort2) Less(i, j int) bool {
	return I[i].QUICKRT > I[j].QUICKRT
}
func (I paySort2) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}
