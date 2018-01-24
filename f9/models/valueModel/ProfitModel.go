package valueModel

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type Profit struct {
	models.Model   `db:"-"`
	COMPCODE       dbr.NullInt64   `db:"COMPCODE"` //公司内码
	ENDDATE        dbr.NullInt64   `db:"ENDDATE"`  //截止日期
	REPORTDATETYPE dbr.NullInt64   `db:"REPORTDATETYPE"`
	EPSBASIC       dbr.NullFloat64 `db:"EPSBASIC"`    //基本每股收益
	ROEWEIGHTED    dbr.NullFloat64 `db:"ROEWEIGHTED"` //净资产收益率——加权
}

func NewProfit() *Profit {
	return &Profit{
		Model: models.Model{
			TableName: "TQ_FIN_PROFINMAININDEX",
			Db:        models.MyCat,
		},
	}
}

func (this *Profit) GetProfitTextData(swlevelcode string) ([]*Profit, error) {
	//logging.Info(scode)
	var data []*Profit
	exps := map[string]interface{}{
		"b.ENDDATE>?":        20161001,
		"a.SWLEVEL1CODE = ?": swlevelcode,
		"a.ISVALID  =?":      1,
		"a.LISTSTATUS =?":    1,
		"a.SETYPE =?":        101,
		"b.ISACTPUB =?":      1,
		"b.REPORTTYPE =?":    1,
		"b.ACCSTACODE =?":    11002,
	}
	builder := this.Db.Select("b.COMPCODE,b.ENDDATE,b.REPORTDATETYPE,b.EPSBASIC,b.ROEWEIGHTED").From("TQ_SK_BASICINFO as a").
		LeftJoin("TQ_FIN_PROFINMAININDEX as b", "a.COMPCODE=b.COMPCODE").Where("a.EXCHANGE in ( '001002' , '001003')").OrderBy("b.COMPCODE,b.ENDDATE desc")
	_, err := this.SelectWhere(builder, exps).LoadStructs(&data)
	if err != nil {
		logging.Debug("%v", err)
		return data, err
	}
	return data, err
}

func (this *Profit) GetProfitChartData(scode string) ([]*Profit, error) {
	logging.Info(scode)
	var data []*Profit
	exps := map[string]interface{}{
		"COMPCODE=?":    scode,
		"ISVALID = ?":   1,
		"ISACTPUB =?":   1,
		"REPORTTYPE =?": 1,
		"ACCSTACODE =?": 11002,
	}
	builder := this.Db.Select("COMPCODE,ENDDATE,REPORTDATETYPE,EPSBASIC,ROEWEIGHTED").From(this.TableName).OrderBy("ENDDATE desc").Limit(10)
	_, err := this.SelectWhere(builder, exps).LoadStructs(&data)
	if err != nil {
		logging.Debug("%v", err)
		return data, err
	}
	return data, err
}
