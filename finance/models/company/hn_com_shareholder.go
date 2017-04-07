package company

import (
	"haina.com/market/finance/models/finchina"
)

// 机构持股
type OrganizationCom struct {
	EntAnn        float64 // 企业年金持股数
	EntAnns       float64 // 企业年金持股所占比例
	Fund          float64 // 基金持股数
	Funds         float64 // 基金持股所占比例
	FinCom        float64 // 财务公司持股数
	FinComs       float64 // 财务公司持股所占比例
	InsCom        float64 // 保险公司持股数
	InsComs       float64 // 保险公司持股所占比例
	LoTruCom      float64 // 信托公司持股数
	LoTruComs     float64 // 信托公司持股所占比例
	OthAge        float64 // 其他机构持股数
	OthAges       float64 // 其他机构持股所占比例
	QFII          float64 // QFII持股数
	QFIIs         float64 // QFII持股所占比例
	SecTra        float64 // 券商持股数
	SecTras       float64 // 券商持股所占比例
	SecTraFinPro  float64 // 券商理财产品持股数
	SecTraFinPros float64 // 券商理财产品持股所占比例
	SocInsFun     float64 // 社保基金持股数
	SocInsFuns    float64 // 社保基金持股所占比例
}

// 返回机构持股json
type OrganZAJson struct {
	EntAnn        float64 // 企业年金持股数
	EntAnns       float64 // 企业年金持股所占比例
	Fund          float64 // 基金持股数
	Funds         float64 // 基金持股所占比例
	FinCom        float64 // 财务公司持股数
	FinComs       float64 // 财务公司持股所占比例
	InsCom        float64 // 保险公司持股数
	InsComs       float64 // 保险公司持股所占比例
	LoTruCom      float64 // 信托公司持股数
	LoTruComs     float64 // 信托公司持股所占比例
	OthAge        float64 // 其他机构持股数
	OthAges       float64 // 其他机构持股所占比例
	QFII          float64 // QFII持股数
	QFIIs         float64 // QFII持股所占比例
	SecTra        float64 // 券商持股数
	SecTras       float64 // 券商持股所占比例
	SecTraFinPro  float64 // 券商理财产品持股数
	SecTraFinPros float64 // 券商理财产品持股所占比例
	SocInsFun     float64 // 社保基金持股数
	SocInsFuns    float64 // 社保基金持股所占比例
}

type OrganZList interface{}
type OrganAList interface{}
type RetOrganInfoJson struct {
	SCode      string      `json:"scode"`
	OrganZList interface{} `json:"OrgZ"`
	OrganAList interface{} `json:"OrgA"`
}

/**获取机构持股信息
 */
func GetCompGroup(scode string) (RetOrganInfoJson, error) {
	// 根据证卷代码查询公司内码跟公告截止日期
	org, compcode, err := finchina.NewTQ_SK_SHAREHOLDER().GetSingleByScode(scode)

	exps := map[string]interface{}{
		"ENDDATE=?":  org.ENDDATE,
		"COMPCODE=?": compcode,
	}
	// 根据公司内码和截止日期获取获取所有的股东信息
	data, err := finchina.NewTQ_SK_SHAREHOLDER().GetListByExps(exps)
	var roij RetOrganInfoJson
	var ozj OrganZAJson // 总股本
	var oaj OrganZAJson // 无限售流通A股
	for _, item := range data {
		switch item.SHHOLDERTYPE {
		case 1:
			// QFII
			ozj.QFII += item.HOLDERAMT
			ozj.QFIIs += item.HOLDERRTO
			if item.SHARESTYPE == "流通A股" {
				oaj.QFII += item.UNLIMHOLDERAMT
				//oaj.QFIIs+=
				// 无限售流通A股所占比例暂时没找到
			}
		case 2:
			// 保险公司
			ozj.InsCom += item.HOLDERAMT
			ozj.InsComs += item.HOLDERRTO
			if item.SHARESTYPE == "流通A股" {
				oaj.InsCom += item.UNLIMHOLDERAMT
			}
		case 3:
			// 财务公司
			ozj.FinCom += item.HOLDERAMT
			ozj.FinComs += item.HOLDERRTO
			if item.SHARESTYPE == "流通A股" {
				oaj.FinCom += item.UNLIMHOLDERAMT
			}
		case 6:
			// 基金
			ozj.Fund += item.HOLDERAMT
			ozj.Funds += item.HOLDERRTO
			if item.SHARESTYPE == "流通A股" {
				oaj.Fund += item.UNLIMHOLDERAMT
			}
		case 8:
			// 全国社保基金
			ozj.SocInsFun += item.HOLDERAMT
			ozj.SocInsFuns += item.HOLDERRTO
			if item.SHARESTYPE == "流通A股" {
				oaj.SocInsFun += item.UNLIMHOLDERAMT
			}
		case 11:
			// 信托公司
			ozj.LoTruCom += item.HOLDERAMT
			ozj.LoTruComs += item.HOLDERRTO
			if item.SHARESTYPE == "流通A股" {
				oaj.LoTruCom += item.UNLIMHOLDERAMT
			}
		case 12:
			// 券商
			ozj.SecTra += item.HOLDERAMT
			ozj.SecTras += item.HOLDERRTO
			if item.SHARESTYPE == "流通A股" {
				oaj.SecTra += item.UNLIMHOLDERAMT
			}
		case 13:
			// 券商理财产品
			ozj.SecTraFinPro += item.HOLDERAMT
			ozj.SecTraFinPros += item.HOLDERRTO
			if item.SHARESTYPE == "流通A股" {
				oaj.SecTraFinPro += item.UNLIMHOLDERAMT
			}
		case 17:
			// 企业年金
			ozj.EntAnn += item.HOLDERAMT
			ozj.EntAnns += item.HOLDERRTO
			if item.SHARESTYPE == "流通A股" {
				oaj.EntAnn += item.UNLIMHOLDERAMT
			}
		case 99:
			// 其他机构
			ozj.OthAge += item.HOLDERAMT
			ozj.OthAges += item.HOLDERRTO
			if item.SHARESTYPE == "流通A股" {
				oaj.OthAge += item.UNLIMHOLDERAMT
			}
		}
	}
	roij.OrganZList = ozj
	roij.OrganAList = oaj
	roij.SCode = scode
	return roij, err

}
