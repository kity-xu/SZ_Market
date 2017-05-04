package company

import (
	"haina.com/market/finance/models/finchina"
)

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

func (this *CompInfo) GetCompInfo(scode string) (*CompInfo, error) {
	var js CompInfo
	v, err := new(finchina.TQ_COMP_INFO).GetCompInfoFromFC(scode)
	if err != nil {
		return &js, err
	}

	js.Account = v.ACCFIRM.String
	js.Addr = v.OFFICEADDR.String
	js.Code = scode
	js.Comp = v.COMPNAME.String
	js.Desc = v.COMPINTRO.String
	js.EDate = v.FOUNDDATE.String
	js.Email = v.COMPEMAIL.String
	//js.Indus = ?
	js.Legal = v.LEGREP.String
	js.License = v.BIZLICENSENO.String
	js.Main = v.MAJORBIZ.String
	js.Manager = v.MANAGER.String
	js.Other = v.BIZSCOPE.String
	js.Postc = v.OFFICEZIPCODE.String
	js.Prov = v.REGION.String
	js.Short = v.COMPSNAME.String
	js.Site = v.REGADDR.String
	js.Tele = v.COMPTEL.String
	js.RegCap = v.REGCAPITAL.Float64
	js.OrgCode = v.ORGCODE.String
	js.ListDate = this.getListDate(scode)

	return &js, err
}

//LISTDATE
func (this *CompInfo) getListDate(scode string) string {
	info, err := new(finchina.SecurityInfo).GetSecurityBasicInfo(scode)
	if err != nil {
		return ""
	}
	return info.LISTDATE.String
}
