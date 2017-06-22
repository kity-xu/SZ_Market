package financemysql

import (
	"github.com/gocraft/dbr"
)

// 数据对象名称：TQ_FIN_PROBALSHEETNEW    中文名称：一般企业资产负债表(新准则产品表)

type TQ_FIN_PROBALSHEETNEW struct {
	TOTASSET      dbr.NullFloat64 `db:"TOTASSET"`      // 资产总计
	TOTALCURRLIAB dbr.NullFloat64 `db:"TOTALCURRLIAB"` // 流动负债合计
	TOTLIAB       dbr.NullFloat64 `db:"TOTLIAB"`       // 负债合计
	CAPISURP      dbr.NullFloat64 `db:"CAPISURP"`      // 资本公积
	TOTCURRASSET  dbr.NullFloat64 `db:"TOTCURRASSET"`  // 流动资产合计
}

//
func (this *TQ_FIN_PROBALSHEETNEW) GetSingleInfo(sess *dbr.Session, comc string) (TQ_FIN_PROBALSHEETNEW, error) {
	var tsp TQ_FIN_PROBALSHEETNEW

	err := sess.Select("TOTASSET,TOTALCURRLIAB,TOTLIAB,CAPISURP,TOTCURRASSET").From("TQ_FIN_PROBALSHEETNEW").
		Where("COMPCODE=" + comc).
		Where("REPORTTYPE=1").
		Where("ISVALID=1").
		OrderBy("ENDDATE DESC ").
		Limit(1).
		LoadStruct(&tsp)
	return tsp, err
}
