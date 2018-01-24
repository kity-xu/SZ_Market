package valueModel

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type Grow struct {
	models.Model   `db:"-"`
	COMPCODE       dbr.NullInt64   `db:"COMPCODE"`       //公司内码
	ENDDATE        dbr.NullInt64   `db:"ENDDATE"`        //截止日期
	PUBLISHDATE    dbr.NullInt64   `db:"PUBLISHDATE"`    //发版日期
	BIZTOTINCO     dbr.NullFloat64 `db:"BIZTOTINCO"`     //营业收入
	PARENETP       dbr.NullFloat64 `db:"PARENETP"`       //归属于母公司所有者的净利润
	REPORTDATETYPE dbr.NullInt64   `db:"REPORTDATETYPE"` //报告期类型
}

func NewGrow() *Grow {
	return &Grow{
		Model: models.Model{
			TableName: "TQ_FIN_PROINCSTATEMENTNEW",
			Db:        models.MyCat,
		},
	}
}

func (this *Grow) GetGrowChartData(compcode string) ([]*Grow, error) {
	var data []*Grow
	exps := map[string]interface{}{
		"COMPCODE=?":   compcode,
		"ISVALID=?":    1,
		"ISACTPUB=?":   1,
		"REPORTTYPE=?": 1,
	}
	builder := this.Db.Select("COMPCODE,PUBLISHDATE,ENDDATE, BIZTOTINCO,PARENETP,REPORTDATETYPE").From(this.TableName).OrderBy("ENDDATE desc").Limit(10)
	_, err := this.SelectWhere(builder, exps).LoadStructs(&data)

	if err != nil {
		logging.Info("---", err.Error())
		return data, err
	}
	return data, err
}
func (this *Grow) GetGrowTextData(swlevelcode string) ([]*Grow, error) {
	var data []*Grow
	exps := map[string]interface{}{
		"a.SWLEVEL1CODE=?": swlevelcode,
		"b.ENDDATE>?":      "20161001",
		"a.ISVALID=?":      1,
		"a.LISTSTATUS=?":   1,
		"a.SETYPE=?":       101,
		"b.REPORTTYPE=?":   1,
		"b.ISACTPUB=?":     1,
	}
	builder := this.Db.Select("b.COMPCODE,b.PUBLISHDATE,b.ENDDATE,b.BIZTOTINCO,b.PARENETP,b.REPORTDATETYPE").
		From("TQ_SK_BASICINFO as a").LeftJoin("TQ_FIN_PROINCSTATEMENTNEW as b", "a.COMPCODE=b.COMPCODE").Where("a.EXCHANGE in ( '001002' , '001003')").OrderBy("b.COMPCODE,b.ENDDATE desc")
	_, err := this.SelectWhere(builder, exps).LoadStructs(&data)
	if err != nil {
		logging.Info(err.Error())
		return data, err
	}
	return data, err
}
