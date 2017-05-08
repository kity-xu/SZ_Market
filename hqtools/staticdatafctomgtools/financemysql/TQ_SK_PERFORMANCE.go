package financemysql

import (
	"github.com/gocraft/dbr"
)

// 数据对象名称：TQ_SK_PERFORMANCE    中文名称：上市公司业绩快报表

type TQ_SK_PERFORMANCE struct {
	TOTASSET dbr.NullFloat64 `db:"TOTASSET"` // 总资产

}

// 查询公司业绩报表
func (this *TQ_SK_PERFORMANCE) GetSingleInfo(sess *dbr.Session, comc string) (TQ_SK_PERFORMANCE, error) {
	var tsp TQ_SK_PERFORMANCE

	err := sess.Select("*").From("TQ_FIN_PROBALSHEETNEW").
		Where("COMPCODE=" + comc + " and  ISVALID=1").OrderBy("PUBLISHDATE  DESC ").Limit(1).LoadStruct(&tsp)
	return tsp, err
}
