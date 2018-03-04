package valueService

import (
	"sort"
	"strconv"

	"haina.com/market/f9/models/finchina"
	"haina.com/market/f9/models/valueModel"
)

type Num struct {
	COMPCODE   int64
	ENDDATE    int64
	TOTALSHAMT int64
}

type sortText []*Num

type numChart struct {
	Date     string `json:"date"`
	Data_bar int64  `json:"data_bar"`
}

var numText chan string

func init() {
	numText = make(chan string)
}

//获取8天数据
func GetNumChartData(numChartData chan []*numChart, compcode string) (chan []*numChart, error) {
	num := []*numChart{}
	data, err := valueModel.NewNum().GetNumChartData(compcode)
	for key, val := range data {
		if key < 8 {
			var n numChart
			n.Date = strconv.FormatInt(val.ENDDATE.Int64, 10)
			n.Data_bar = val.TOTALSHAMT.Int64
			num = append(num, &n)
		} else {
			break
		}
	}
	t1 := data[0].TOTALSHAMT.Int64
	t2 := data[1].TOTALSHAMT.Int64
	t3 := data[2].TOTALSHAMT.Int64
	if t1-t2 > 0 && t2-t3 > 0 {
		numText <- "连续增加"
	} else if t1-t2 < 0 && t2-t3 < 0 {
		numText <- "连续减少"
	} else {
		numText <- "股东人数无明显变化"
	}
	numChartData <- num
	return numChartData, err
}

//获取文本信息
func GetNumTextData(numChartText chan string, detail finchina.CompanyDetail, leng string) (chan string, error) {
	num := []*Num{}
	data, err := valueModel.NewNum().GetNumTextData(detail.SWLEVEL1CODE)
	for key, val := range data {
		if key < 8 {
			var n Num
			n.COMPCODE = val.COMPCODE.Int64
			n.ENDDATE = val.ENDDATE.Int64
			n.TOTALSHAMT = val.TOTALSHAMT.Int64
			num = append(num, &n)
		} else {
			break
		}
	}
	sort.Sort(sortText(num))
	var TOTALSHAMT_key int
	for key, val := range num {
		com, _ := strconv.ParseInt(detail.COMPCODE, 10, 64)
		if com == val.COMPCODE {
			TOTALSHAMT_key = key + 1
			break
		}
	}
	numChartText <- "近期，该公司股东人数" + <-numText + "，在" + detail.SWLEVEL1NAME + "行业排名 " + strconv.Itoa(TOTALSHAMT_key) + "/" + leng + "。"
	return numChartText, err
}

func (I sortText) Len() int {
	return len(I)
}
func (I sortText) Less(i, j int) bool {
	return I[i].TOTALSHAMT > I[j].TOTALSHAMT
}
func (I sortText) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}
