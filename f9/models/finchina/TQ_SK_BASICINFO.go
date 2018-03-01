package finchina

import (
	. "haina.com/market/f9/models"
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
			TableName: TABLE_TQ_SK_BASICINFO,
			Db:        models.MyCat,
		},
	}
}

//获取证券基本信息
func (this *CompanyDetail) GetCompanyDetail(secode string) (CompanyDetail, error) {
	var data CompanyDetail
	exps := map[string]interface{}{
		"SECODE=?": secode,
	}

	builder := this.Db.Select("SYMBOL,SESNAME,TOTALSHARE,SWLEVEL1CODE,SWLEVEL1NAME,COMPCODE,EXCHANGE,LISTDATE,LISTSTATUS").
		From(this.TableName).
		Where("ISVALID = 1  and SETYPE = 101 and LISTSTATUS=1").Limit(1)

	err := this.SelectWhere(builder, exps).LoadStruct(&data)
	if err != nil {
		logging.Debug("%v", err)
	}
	return data, err
}

//获取该行业下的所有公司
func (this *CompanyDetail) GetAllCompany(swlevelcode string) ([]*CompanyDetail, error) {
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
