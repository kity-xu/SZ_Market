package valueService

import (
	"fmt"
	"strconv"
	"time"

	"haina.com/market/f9/models/common"
	"haina.com/market/f9/services"
	"haina.com/share/logging"
)

var (
	symbolParam   string //证券代码（带字母）
	symbol        string //证券代码（不带字母)
	symboType     string //证券类型
	swlevelcode   string //行业代码
	swlevelname   string //行业名称
	compcode      string //公司代码
	compname      string //公司名称
	companyDetail string //公司详情
	//sumcompay     string //该行业的所有公司
)

type RetVal struct {
	Code int    `json:"code,omitempty"`
	Desc string `json:"desc,omitempty"`
}
type value struct {
	RetVal

	ProfitChartText string         `json:"profit_ability,omitempty"`
	ProfitChartData []*ProfitChart `json:"per_share_earn_chart,omitempty"`

	GrowChartText string       `json:"grow_ability,omitempty"`
	GrowChartData []*growChart `json:"main_net_income_chart,omitempty"`

	PayChartText string      `json:"insolvency_ability,omitempty"`
	PayChartData []*payChart `json:"liability_assets_ratio_chart,omitempty"`

	NumChartText string      `json:"shareholders_ability,omitempty"`
	NumChartData []*numChart `json:"shareholders_total_chart,omitempty"`
}

func GetValueData(scode string) (value, error) {
	symbol = scode[2:]
	symbolParam = scode
	if scode[0:2] == "sz" {
		symboType = "001003"
	} else if scode[0:2] == "sh" {
		symboType = "001002"
	}

	var v value

	detail, err := companyDetailModel.NewCompanyDetail().GetCompanyDetail(symbol, symboType)

	//logging.Info("detail==".)

	if err != nil {
		logging.Info("该股票不存在")
		v.Code = 20000
		v.Desc = "该股票不存在"
		return v, err
	}

	if detail.LISTSTATUS != 1 {
		logging.Info("非股票代码")
		v.Code = 20000
		v.Desc = "非股票代码"
		return v, err
	}

	//redis.Set("kk", []byte("1234567"))

	//logging.Info(redis.Get("kk"))

	now := time.Now()

	m, _ := time.ParseDuration("2017-12-05")
	m1 := now.Add(m)
	fmt.Println(m1)
	fmt.Println(now)

	swlevelcode = detail.SWLEVEL1CODE
	swlevelname = detail.SWLEVEL1NAME
	compcode = detail.COMPCODE
	compname = detail.SESNAME

	allCompany, _ := commonService.IndustryOfAllCompany(swlevelcode)

	leng := strconv.Itoa(len(allCompany))

	logging.Info("allCompany=====%+v", len(allCompany))

	//	logging.Info("------scode=%v-----", scode)
	//	logging.Info("------symbol=%v-----", symbol)
	//	logging.Info("------symbolParam=%v------", symbolParam)
	//	logging.Info("------swlevelcode=%v------", swlevelcode)
	//	logging.Info("------compcode=%v------", compcode)
	//	logging.Info("----公司名称compname=%v------", compname)

	ProfitChartText := make(chan string)
	ProfitChartData := make(chan []*ProfitChart)
	payChartDext := make(chan string)
	payChartData := make(chan []*payChart)
	numChartText := make(chan string)
	numChartData := make(chan []*numChart)
	go GetProfitTextData(ProfitChartText, leng)         //盈利能力本文显示
	go GetProfitChartData(ProfitChartData, compcode)    //盈利能力的8条数据
	go GetPayChartData(payChartData, compcode)          //偿债能力的8条数据
	go GetPayChartText(payChartDext, swlevelcode, leng) //偿债能力文本
	go GetNumChartData(numChartData, compcode)          //股东人数的8条数据
	go GetNumTextData(numChartText, swlevelcode, leng)  //股东人数文本

	v.ProfitChartText = <-ProfitChartText
	v.ProfitChartData = <-ProfitChartData
	logging.Info("等待中...")

	GrowData, err := GetGrowChartData(leng) //成长能力的8条数据以及文本
	v.GrowChartText = GrowData.ChartText
	v.GrowChartData = GrowData.ChartData

	v.PayChartText = <-payChartDext
	v.PayChartData = <-payChartData

	v.NumChartText = <-numChartText
	v.NumChartData = <-numChartData
	logging.Info("完成了 ...")

	return v, err
}
