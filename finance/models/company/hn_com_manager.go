package company

import (
	"haina.com/market/finance/models/finchina"
)

type HnManager struct {
	//Pscode     string  `json:"Pscode"` //高管代码
	Name        string  `json:"Name"`    //高管姓名
	Duty        string  `json:"Duty"`    //职务
	HoldAMT     float64 `json:"Holdamt"` //持股数
	UpdateDate  string  `json:"UpDate"`  //公司高管资料更新日期
	PublistDate string  `json:"PubDate"` //持股变动公布日期
}

func (this *HnManager) GetManagerList(scode string) (*[]*HnManager, error) {
	list := make([]*HnManager, 0)
	primal, err := new(finchina.TQ_COMP_MANAGER).GetManagersFromFC(scode) //公司高管表
	if err != nil {
		return &list, err
	}

	managers := managersToOnly(primal)

	amts, err := this.getManagersHoldAMT(scode)
	if err != nil {
		return &list, err
	}

	for _, v := range managers {
		var js HnManager
		//js.Pscode = v.CNAME.String
		js.PublistDate = amts[v.CNAME.String].PUBLISHDATE.String
		js.Duty = v.ACTDUTYNAME.String
		js.Name = v.CNAME.String
		js.UpdateDate = v.UPDATEDATE.String
		js.HoldAMT = amts[v.CNAME.String].HOLDAFAMT.Float64 //- amts[v.CNAME.String].HOLDBEAMT.Float64  //注：TQ_COMP_MANAGER中的PERSONALCODE（人物代码）与TQ_COMP_SKHOLDERCHG中的PSCODE（高管代码）是一回事
		list = append(list, &js)

	}

	return &list, err
}

//获取高管持股数
func (this *HnManager) getManagersHoldAMT(scode string) (map[string]finchina.HolderChange, error) {
	amts := make(map[string]finchina.HolderChange)
	holders, err := new(finchina.TQ_COMP_SKHOLDERCHG).GetHoldAMTlistFromFC(scode) //高管及关联人持股变动表

	for _, v := range hoderChangeToOnly(holders) { //以高管代码（PSCODE）为key与结构体建立一一对应关系
		amts[v.HOLDNAME.String] = v
	}
	return amts, err //返回map
}

//高管信息表数据去重，取UPDATEDATE最新
func managersToOnly(primal []finchina.TQ_COMP_MANAGER) []finchina.TQ_COMP_MANAGER {
	swap := make(map[string]finchina.TQ_COMP_MANAGER)

	var managers []finchina.TQ_COMP_MANAGER
	for _, v := range primal {
		if _, ok := swap[v.CNAME.String]; !ok {
			swap[v.CNAME.String] = v //不存在
		} else { //存在
			update := swap[v.CNAME.String]
			update.ACTDUTYNAME.String = v.ACTDUTYNAME.String + ", " + update.ACTDUTYNAME.String
			swap[v.CNAME.String] = update
		}
	}
	var count int = 0
	for _, v := range swap {
		managers = append(managers, v)
		count++
	}

	return managers
}

//高管持股变动表去重，取PUBLISHDATE最新
func hoderChangeToOnly(primal []finchina.HolderChange) []finchina.HolderChange {
	swap := make(map[string]finchina.HolderChange)

	var holders []finchina.HolderChange
	for _, v := range primal {
		if _, ok := swap[v.HOLDNAME.String]; !ok {
			swap[v.HOLDNAME.String] = v //不存在
		}
	}
	var count int = 0
	for _, v := range swap {
		holders = append(holders, v)
		count++
	}

	return holders
}
