package company

import (
	"haina.com/market/finance/models/finchina"
)

type AdministratorJson struct {
}
type InfoJson struct {
	//Secode string `json:"Secode"`
	Sname     string `json:"Sname"`     // 股票名称
	PScode    string `json:"PScode"`    // 高管代码
	PSname    string `json:"PSname"`    // 高管名称
	Remark    string `json:"Remark"`    // 简历
	Gender    string `json:"Gender"`    // 性别
	Year      string `json:"Year"`      // 年龄
	Educa     string `json:"Educa"`     // 教育程度
	Post      string `json:"Post"`      // 现任职务
	Prov      string `json:"Prov"`      // 所属地区
	Orgtp     string `json:"Orgtp"`     // 所有制形式
	CSRC      string `json:"CSRC"`      // 证监会行业名称
	SWname    string `json:"SWname"`    // 申万行业名称
	GICS      string `json:"GICS"`      // GICS行业名称
	Isvalid   int64  `json:"Valid"`     // 是否有效
	Tmstamp   int64  `json:"Tmstamp"`   // 时间戳
	EntryDate string `json:"EntryDate"` // 录入日期
	EntryTime string `json:"EntryTime"` // 录入时间
}

type EquityChangeJson struct {
	Date  string  `json:"Date"`  //日期
	Name  string  `json:"Name"`  //变动人
	VolCh float64 `json:"VolCh"` //变动数量（股）
	Vol   float64 `json:"Vol"`   //变动后持股数（结存股票（股））
	Price float64 `json:"price"` //交易均价
	Duty  string  `json:"Duty"`  //职务
}

func (this *AdministratorJson) GetAdminListJson(secode string) (*[]*InfoJson, error) {
	list := make([]*InfoJson, 0)
	pss, err := new(finchina.Administrator).GetPSList(secode)
	if err != nil {
		return &list, err
	}
	list = getAdminList(pss)
	return &list, err
}

func (this *AdministratorJson) GetAdminEquityChangeJson(secode string) (*[]*EquityChangeJson, error) {
	list := make([]*EquityChangeJson, 0)
	eqs, err := new(finchina.Administrator).GetPSEquityChange(secode)
	if err != nil {
		return &list, err
	}
	list = getAdminEquityChange(eqs)
	return &list, err
}
func getAdminList(pss []finchina.AdminInfo) []*InfoJson {
	list := make([]*InfoJson, 0)
	for _, v := range pss {
		var ps InfoJson
		ps.CSRC = v.CSRCNAME.String
		ps.Educa = v.EDUCATIONLEVEL.String
		ps.EntryDate = v.ENTRYDATE.String
		ps.EntryTime = v.ENTRYTIME.String
		ps.Gender = v.GENDER.String
		ps.GICS = v.GICSNAME.String
		ps.Isvalid = v.ISVALID.Int64
		ps.Orgtp = v.ORGTYPE.String
		ps.Post = v.POST.String
		ps.Prov = v.PROVINCENAME.String
		ps.PScode = v.PSCODE.String
		ps.PSname = v.PSNAME.String
		ps.Remark = v.REMARK.String
		ps.Sname = v.SNAME.String
		ps.SWname = v.SWNAME.String
		ps.Tmstamp = v.TMSTAMP.Int64
		ps.Year = v.YEAR.String

		list = append(list, &ps)
	}
	return list
}

func getAdminEquityChange(eqs []finchina.AdminEquityChange) []*EquityChangeJson {
	list := make([]*EquityChangeJson, 0)
	for _, v := range eqs {
		var eq EquityChangeJson
		eq.Vol = v.AFSHAREAMT.Float64
		eq.Price = v.TOTAVGPRICE.Float64
		eq.Date = v.ENDDATE.String
		eq.Duty = v.DUTY.String
		eq.Name = v.SHHOLDERNAME.String
		eq.VolCh = v.AFSHAREAMT.Float64 - v.BFSHAREAMT.Float64

		list = append(list, &eq)
	}
	return list
}
