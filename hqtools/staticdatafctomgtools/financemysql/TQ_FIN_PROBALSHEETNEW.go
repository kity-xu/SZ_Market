package financemysql

import (
	"github.com/gocraft/dbr"
)

// 数据对象名称：TQ_FIN_PROBALSHEETNEW    中文名称：一般企业资产负债表(新准则产品表)

type TQ_FIN_PROBALSHEETNEW struct {
	TOTCURRASSET dbr.NullFloat64 `db:"TOTCURRASSET"` // 流动资产合计
}

// 查询公司业绩报表
func (this *TQ_FIN_PROBALSHEETNEW) GetSingleInfo(sess *dbr.Session, comc string) (TQ_FIN_PROBALSHEETNEW, error) {
	var tsp TQ_FIN_PROBALSHEETNEW

	err := sess.Select("*").From("TQ_FIN_PROBALSHEETNEW").
		Where("COMPCODE=" + comc + " and  ISVALID=1").OrderBy("PUBLISHDATE DESC ").Limit(1).LoadStruct(&tsp)
	return tsp, err
}
