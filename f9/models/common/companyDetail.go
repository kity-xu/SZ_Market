package companyDetailModel

import (
	//"niuniu/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type CompanyDetail struct {
	models.Model `db:"-"`
	SYMBOL       string `db:"SYMBOL"`
	SESNAME      string `db:"SESNAME"`
	TOTALSHARE   string `db:"TOTALSHARE"`
	SWLEVEL1CODE string `db:"SWLEVEL1CODE"`
	SWLEVEL1NAME string `db:"SWLEVEL1NAME"`
	COMPCODE     string `db:"COMPCODE"`
	EXCHANGE     string `db:"EXCHANGE"`
	LISTDATE     string `db:"LISTDATE"`
	LISTSTATUS   int    `db:"LISTSTATUS"`
}

func NewCompanyDetail() *CompanyDetail {
	return &CompanyDetail{
		Model: models.Model{
			TableName: "TQ_SK_BASICINFO",
			Db:        models.MyCat,
		},
	}
}
func (this *CompanyDetail) GetCompanyDetail(symbol string, symboType string) (CompanyDetail, error) {
	var data CompanyDetail
	exps := map[string]interface{}{
		"a.SYMBOL=?":   symbol,
		"a.EXCHANGE=?": symboType,
	}

	builder := this.Db.Select("a.SYMBOL,a.SESNAME,a.TOTALSHARE,a.SWLEVEL1CODE, a.SWLEVEL1NAME, a.COMPCODE,a.EXCHANGE,a.LISTDATE,a.LISTSTATUS").From(this.TableName + " as a").
		Where("a.ISVALID = 1  and a.SETYPE = 101 and a.LISTSTATUS=1").Limit(1)
	err := this.SelectWhere(builder, exps).LoadStruct(&data)
	if err != nil {
		logging.Debug("%v", err)
		return data, err
	}
	return data, err
}

//搞行业下的所有公司
func (this CompanyDetail) GetAllCompany(swlevelcode string) ([]*CompanyDetail, error) {
	var data []*CompanyDetail
	builder := this.Db.SelectBySql(`select  COMPCODE,EXCHANGE,SYMBOL  from TQ_SK_BASICINFO
              WHERE SWLEVEL1CODE=? and ISVALID = 1 and LISTSTATUS = 1 and SETYPE = 101 and  EXCHANGE in ( '001002' , '001003')`, swlevelcode)
	_, err := this.SelectWhere(builder, nil).LoadStructs(&data)

	if err != nil {
		logging.Info(err.Error())
		return data, err
	}
	return data, err

}
