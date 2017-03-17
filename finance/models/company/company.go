package company

type Company struct {
	Account string `json:"Account"` //会计师事务所
	Addr    string `json:"Addr"`    //公司办公地址
	Code    string `json:"Code"`    //A股证券代码
	Comp    string `json:"Comp"`    //公司中文名称
	Desc    string `json:"Desc"`    //公司简介
	EDate   int    `json:"EDate"`   //公司创建日期
	Email   string `json:"Email"`   //联系人邮箱
	Indus   string `json:"Indus"`   //公司所属证监会行业
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
}

func NewCompany() *Company {
	return &Company{}
}
