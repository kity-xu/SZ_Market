package valueService

import (
	"strconv"

	"haina.com/market/f9/models/finchina"
	"haina.com/market/f9/services"
	"haina.com/share/logging"
)

type value struct {
	ProfitChartText string         `json:"profit_ability,omitempty"`
	ProfitChartData []*ProfitChart `json:"per_share_earn_chart,omitempty"`

	GrowChartText string       `json:"grow_ability,omitempty"`
	GrowChartData []*growChart `json:"main_net_income_chart,omitempty"`

	PayChartText string      `json:"insolvency_ability,omitempty"`
	PayChartData []*payChart `json:"liability_assets_ratio_chart,omitempty"`

	NumChartText string      `json:"shareholders_ability,omitempty"`
	NumChartData []*numChart `json:"shareholders_total_chart,omitempty"`
}

func GetValueData(sid string) (*value, error) {
	var v value
	sc := finchina.NewTQ_OA_STCODE()
	if err := sc.GetCompcode(sid); err != nil {
		return nil, err
	}
	detail, err := finchina.NewCompanyDetail().GetCompanyDetail(sc.SECODE.String)
	if err != nil || detail.LISTSTATUS != 1 {
		logging.Error("该股票不存在或股票 已退市")
		return nil, err
	}

	allCompany, _ := commonService.IndustryOfAllCompany(detail.SWLEVEL1CODE)

	leng := strconv.Itoa(len(allCompany))

	logging.Info("allCompany=====%+v", len(allCompany))

	ProfitChartText := make(chan string)
	ProfitChartData := make(chan []*ProfitChart)
	payChartDext := make(chan string)
	payChartData := make(chan []*payChart)
	numChartText := make(chan string)
	numChartData := make(chan []*numChart)
	go GetProfitTextData(ProfitChartText, detail, leng)     //盈利能力本文显示
	go GetProfitChartData(ProfitChartData, detail.COMPCODE) //盈利能力的8条数据
	go GetPayChartData(payChartData, detail.COMPCODE)       //偿债能力的8条数据
	go GetPayChartText(payChartDext, detail, leng)          //偿债能力文本
	go GetNumChartData(numChartData, detail.COMPCODE)       //股东人数的8条数据
	go GetNumTextData(numChartText, detail, leng)           //股东人数文本

	v.ProfitChartText = <-ProfitChartText
	v.ProfitChartData = <-ProfitChartData
	logging.Info("等待中...")

	GrowData, err := GetGrowChartData(leng, sc.COMPCODE.String, detail) //成长能力的8条数据以及文本
	v.GrowChartText = GrowData.ChartText
	v.GrowChartData = GrowData.ChartData

	v.PayChartText = <-payChartDext
	v.PayChartData = <-payChartData

	v.NumChartText = <-numChartText
	v.NumChartData = <-numChartData
	logging.Info("完成了 ...")

	return &v, err
}
