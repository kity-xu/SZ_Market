package fcmysql

import (
	"github.com/gocraft/dbr"
)

// 数据对象名称：TQ_SK_PERFORMANCE    中文名称：上市公司业绩快报表

type TQ_SK_PERFORMANCE struct {
	BASICEPS    dbr.NullFloat64 `db:"TOTASSET"`    // 基本每股收益
	WEIGHTEDROE dbr.NullFloat64 `db:"WEIGHTEDROE"` // 净资产收益率(加权)
	TOTSHAREQUI dbr.NullFloat64 `db:"TOTSHAREQUI"` // 股东权益合计
	NAPS        dbr.NullFloat64 `db:"NAPS"`        // 每股净资产
}

// 查询公司业绩报表
func (this *TQ_SK_PERFORMANCE) GetSingleInfo(sess *dbr.Session, comc string) (TQ_SK_PERFORMANCE, error) {
	var tsp TQ_SK_PERFORMANCE

	err := sess.Select("BASICEPS,WEIGHTEDROE,TOTSHAREQUI,NAPS").From("TQ_SK_PERFORMANCE").
		Where("COMPCODE=" + comc + " and  ISVALID=1").
		OrderBy("PUBLISHDATE  DESC ").
		Limit(1).
		LoadStruct(&tsp)
	return tsp, err
}
