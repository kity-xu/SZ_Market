package company

type Dividend struct {
}

//Dividend
type Div struct {
	Bene     string `json:"Bene"`     //分红对象
	Bonus    string `json:"Bonus"`    //送股（股）
	Date     string `json:"Data"`     //年度
	Dividend string `json:"Dividend"` //分红（元，税前）
	DivDate  string `json:"DivDate"`  //红利发放日
	DivRate  string `json:"DivRate"`  //股利支付率（%）
	DNDate   string `json:"DNDate"`   //决案公布日
	Evolve   string `json:"Evolve"`   //事情进展
	ExDate   string `json:"ExDate"`   //除权除息日
	INDate   string `json:"INDate"`   //实施公告日
	LisDate  string `json:"LisDate"`  //转股上市日
	PNDate   string `json:"PNDate"`   //预案公布日
	RegDate  string `json:"RegDate"`  //股权登记日
	Tran     string `json:"Tran"`     //转股（股）
}

//Seasoned Equity Offerings
type SEO struct {
	AGMD    string `json:"AGMD"`    //股东大会决议公告日
	IECD    string `json:"IECD"`    //发审委公告日
	LisDate string `json:"LisDate"` //新股上市日
	PNDate  string `json:"PNDate"`  //预案公布日
	PPrice  string `json:"PPrice"`  //预案发行价格
	Price   string `json:"Price"`   //实际发行价格
	PVal    string `json:"PVal"`    //预案募资金额
	PVol    string `json:"PVol"`    //预案发行数量
	Range   string `json:"Range"`   //发行对象类型
	SEOD    string `json:"SEOD"`    //发行新股日
	SRCD    string `json:"SRCD"`    //证监会核准公告日
	Step    string `json:"Step"`    //事情进展
	Type    string `json:"Type"`    //发行方式
	Val     string `json:"Val"`     //实际募资金额
	Vol     string `json:"Vol"`     //实际发行数量
}

//Rights Offering
type RO struct {
	AGMD    string `json:"AGMD"`    //股东大会决议公告日
	Code    string `json:"Code"`    //配股代码
	DNDate  string `json:"DNDate"`  //决案公布日
	ERDate  string `json:"ERDate"`  //配股除权日
	LisDate string `json:"LisDate"` //配股上市日
	PNDate  string `json:"PNDate"`  //预案公布日
	PProp   string `json:"PProp"`   //计划配股比例
	Price   string `json:"Price"`   //实际配股价格
	Prop    string `json:"Prop"`    //实际配股比例
	PVol    string `json:"PVol"`    //计划配股数量
	RegDate string `json:"RegDate"` //股权登记日
	ROPD    string `json:"ROPD"`    //配股缴款起止日
	Short   string `json:"Short"`   //配股简称
	Vol     string `json:"Vol"`     //实际配股数量
}

// Repo(repurchase agreement)
type Repo struct {
}

func (this *Dividend) GetDiv() *Div {
	return &Div{}
}
func (this *Dividend) GetSEO() *SEO {
	return &SEO{}
}
func (this *Dividend) GetRO() *RO {
	return &RO{}
}
func (this *Dividend) GetRepo() *Repo {
	return &Repo{}
}
