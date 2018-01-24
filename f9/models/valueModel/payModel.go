package valueModel

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type Pay struct {
	models.Model `db:"-"`
	COMPCODE     dbr.NullInt64   `db:"COMPCODE"`
	ENDDATE      dbr.NullInt64   `db:"ENDDATE"`
	ASSLIABRT    dbr.NullFloat64 `db:"ASSLIABRT"`
	QUICKRT      dbr.NullFloat64 `db:"QUICKRT"`
}

func NewPay() *Pay {
	return &Pay{
		Model: models.Model{
			TableName: "TQ_FIN_PROINDICDATA",
			Db:        models.MyCat,
		},
	}
}

func (this *Pay) GetPayChartData(compcode string) ([]*Pay, error) {
	var data []*Pay //data := []*payChart{}
	exps := map[string]interface{}{
		"COMPCODE=?":   compcode,
		"ISVALID=?":    1,
		"REPORTTYPE=?": 3,
	}

	builder := this.Db.Select("COMPCODE,ENDDATE,ASSLIABRT,QUICKRT").From("TQ_FIN_PROINDICDATA").OrderBy("ENDDATE DESC").Limit(8)
	_, err := this.SelectWhere(builder, exps).LoadStructs(&data)
	if err != nil {
		logging.Info(err.Error())
		return data, err
	}
	return data, err
}
func (this *Pay) GetPayChartText(swlevelcode string) ([]*Pay, error) {
	data := []*Pay{}
	//	ssql := "(SELECT b.COMPCODE,b.ENDDATE,b.ASSLIABRT,b.QUICKRT FROM TQ_SK_BASICINFO AS a LEFT JOIN TQ_FIN_PROINDICDATA AS b ON a.COMPCODE = b.COMPCODE" +
	//		"WHERE a.SWLEVEL1CODE = " + swlevelcode + " and a.ISVALID = 1 AND a.LISTSTATUS = 1 AND a.SETYPE = 101 AND a.EXCHANGE IN ('001002', '001003') AND b.ISVALID = 1" +
	//		"AND b.REPORTTYPE = 3 order BY b.ENDDATE DESC) as c"
	//	builder := this.Db.Select("c.COMPCODE,c.ENDDATE,c.ASSLIABRT,c.QUICKRT").From(ssql)
	//	_, err := this.SelectWhere(builder, nil).LoadStructs(&data)
	//	if err != nil {
	//		logging.Info(err.Error())
	//		return data, err
	//	}
	//	return data, err

	builder := this.Db.SelectBySql(`(SELECT
                c.COMPCODE,
                c.ENDDATE,
                c.ASSLIABRT,
                c.QUICKRT
                FROM (SELECT
                        b.COMPCODE,
                        b.ENDDATE,
                        b.ASSLIABRT,
                        b.QUICKRT
                        FROM
                            TQ_SK_BASICINFO AS a
                        LEFT JOIN TQ_FIN_PROINDICDATA AS b ON a.COMPCODE = b.COMPCODE
                        WHERE
                             a.SWLEVEL1CODE = ?
                        and a.ISVALID = 1
                        AND a.LISTSTATUS = 1
                        AND a.SETYPE = 101
                        AND a.EXCHANGE IN ('001002', '001003')
                        AND b.ISVALID = 1
                        AND b.REPORTTYPE = 3 order BY b.ENDDATE DESC) as c
                 GROUP BY c.COMPCODE)`, swlevelcode)
	_, err := this.SelectWhere(builder, map[string]interface{}{}).LoadStructs(&data)

	if err != nil {
		logging.Info(err.Error())
		return data, err
	}
	return data, err

}
