package f10

import (
"haina.com/market/finance/models/finchina"
"haina.com/share/logging"
)

var Province map[string]string

type CompInfo struct {
	Account string `json:"Account"` //会计师事务所
	Addr    string `json:"Addr"`    //公司办公地址
	Code    string `json:"Code"`    //A股证券代码
	Comp    string `json:"Comp"`    //公司中文名称
	Desc    string `json:"Desc"`    //公司简介
	EDate   string `json:"EDate"`   //公司成立日期
	Email   string `json:"Email"`   //联系人电子邮箱
	Indus   string `json:"Indus"`   //公司所属证监会行业（聚源）
	Legal   string `json:"Legal"`   //法人代表
	License string `json:"License"` //企业法人营业执照注册号
	Main    string `json:"Main"`    //经营范围-主营
	Manager string `json:"Manager"` //总经理
	Other   string `json:"Other"`   //经营范围-兼营
	Postc   string `json:"Postc"`   //公司办公地址邮编
	Prov    string `json:"Prov"`    //省份
	Short   string `json:"Short"`   //A股证券简称
	Site    string `json:"Site"`    //首次注册登记地点
	Tele    string `json:"Tele"`    //联系人电话

	RegCap   float64 `json:"RegCap"`   //注册资本
	OrgCode  string  `json:"OrgCode"`  //机构组织 代码
	ListDate string  `json:"ListDate"` //上市日期
}

func (this *CompInfo) GetCompInfo(scode string, market string) (*CompInfo, error) {
	var js CompInfo
	v, err := new(finchina.TQ_COMP_INFO).GetCompInfoFromFC(scode, market)
	if err != nil {
		return &js, err
	}

	indus, e := finchina.NewTQ_COMP_INDUSTRY().GetCompTrade(scode, market)
	if e != nil {
		logging.Error("Error accessing company industry information...")
	}

	js.Account = v.ACCFIRM.String
	js.Addr = v.OFFICEADDR.String
	js.Code = scode
	js.Comp = v.COMPNAME.String
	js.Desc = v.COMPINTRO.String
	js.EDate = v.FOUNDDATE.String
	js.Email = v.COMPEMAIL.String
	js.Indus = indus.LEVEL2NAME.String
	js.Legal = v.LEGREP.String
	js.License = v.BIZLICENSENO.String
	js.Main = v.MAJORBIZ.String
	js.Manager = v.MANAGER.String
	js.Other = v.BIZSCOPE.String
	js.Postc = v.OFFICEZIPCODE.String
	js.Prov = getProvince(v.REGION.String)
	js.Short = v.COMPSNAME.String
	js.Site = v.REGADDR.String
	js.Tele = v.COMPTEL.String
	js.RegCap = v.REGCAPITAL.Float64
	js.OrgCode = v.ORGCODE.String
	js.ListDate = this.getListDate(scode, market)

	return &js, err
}

//LISTDATE
func (this *CompInfo) getListDate(scode string, market string) string {
	info, err := new(finchina.SecurityInfo).GetSecurityBasicInfo(scode, market)
	if err != nil {
		return ""
	}
	return info.LISTDATE.String
}

func getProvince(pro string) string {

	if Province[pro] == "" {
		initProvince()
		logging.Debug("----------initProvince")
	}
	return Province[pro]
}

func initProvince() {
	Province = make(map[string]string)
	Province["CN"] = "全国"
	Province["CN110000"] = "北京"
	Province["CN120000"] = "天津"
	Province["CN130000"] = "河北"
	Province["CN140000"] = "山西"
	Province["CN150000"] = "内蒙古"
	Province["CN210000"] = "辽宁"
	Province["CN220000"] = "吉林"
	Province["CN230000"] = "黑龙江"
	Province["CN310000"] = "上海"
	Province["CN320000"] = "江苏"
	Province["CN320100"] = "南京"
	Province["CN320500"] = "苏州"
	Province["CN330000"] = "浙江"
	Province["CN330100"] = "杭州"
	Province["CN330200"] = "宁波"
	Province["CN340000"] = "安徽"
	Province["CN350000"] = "福建"
	Province["CN350200"] = "厦门"
	Province["CN360000"] = "江西"
	Province["CN370000"] = "山东"
	Province["CN410000"] = "河南"
	Province["CN410100"] = "郑州"
	Province["CN420000"] = "湖北"
	Province["CN420100"] = "武汉"
	Province["CN430000"] = "湖南"
	Province["CN440000"] = "广东"
	Province["CN440100"] = "广州"
	Province["CN440300"] = "深证"
	Province["CN450000"] = "广西"
	Province["CN460000"] = "海南"
	Province["CN510000"] = "四川"
	Province["CN520000"] = "贵州"
	Province["CN530000"] = "云南"
	Province["CN540000"] = "西藏"
	Province["CN500000"] = "重庆"
	Province["CN610000"] = "陕西"
	Province["CN620000"] = "甘肃"
	Province["CN630000"] = "青海"
	Province["CN640000"] = "宁夏"
	Province["CN650000"] = "新疆"
	Province["CN810000"] = "香港"
	Province["CN820000"] = "澳门"
	Province["CN710000"] = "台湾"
	Province["99999998"] = "境外"
	Province["99999999"] = "其他"
}
